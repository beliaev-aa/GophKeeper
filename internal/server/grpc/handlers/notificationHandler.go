// Package handlers содержит обработчики для серверных вызовов gRPC, в частности для уведомлений.
package handlers

import (
	"beliaev-aa/GophKeeper/pkg/proto"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"slices"
	"sync"
)

// Представляет подписчика на уведомления.
type subscriber struct {
	stream   proto.Notification_SubscribeServer
	id       uint64
	finished chan<- bool
}

// NotificationHandler управляет подписчиками на уведомления.
type NotificationHandler struct {
	proto.UnimplementedNotificationServer
	logger      *zap.Logger
	subscribers sync.Map
}

// NewNotificationHandler создаёт новый экземпляр NotificationHandler.
func NewNotificationHandler(logger *zap.Logger) *NotificationHandler {
	return &NotificationHandler{logger: logger}
}

// Subscribe обрабатывает подписку клиента на уведомления.
// Возвращает ошибку, если подписка не может быть выполнена.
func (s *NotificationHandler) Subscribe(in *proto.SubscribeRequest, stream proto.Notification_SubscribeServer) error {
	ctx := stream.Context()

	userID, err := extractUserID(ctx)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("received subscribe from client", zap.Int("client_id", int(in.Id)), zap.Int("user_id", int(userID)))

	fin := make(chan bool)
	v, ok := s.subscribers.Load(userID)
	var subscribers []subscriber
	if ok {
		subscribers, ok = v.([]subscriber)
		if !ok {
			return status.Error(codes.Internal, "failed to cast subscribers")
		}
	}

	subscribers = append(subscribers, subscriber{stream: stream, id: in.Id, finished: fin})
	s.subscribers.Store(userID, subscribers)

	for {
		select {
		case <-fin:
			s.logger.Info("closing stream for client", zap.Int("client_id", int(in.Id)))
			return nil
		case <-ctx.Done():
			s.logger.Info("client has disconnected", zap.Int("client_id", int(in.Id)))
			return nil
		}
	}
}

// notifyClients отправляет уведомления всем подписчикам, за исключением инициатора изменений.
// Возвращает ошибку, если не удаётся отправить уведомление или подписчиков нет.
func (s *NotificationHandler) notifyClients(userID uint64, clientID uint64, ID uint64, updated bool) error {
	v, ok := s.subscribers.Load(userID)
	if !ok {
		return errors.New("no subscribers")
	}

	subs, ok := v.([]subscriber)
	if !ok {
		return errors.New("failed to cast to subs")
	}

	var unsubscribes []int
	for i, sub := range subs {
		if sub.id == clientID {
			continue
		}
		resp := &proto.SubscribeResponse{Id: ID, Updated: updated}
		if err := sub.stream.Send(resp); err != nil {
			s.logger.Error("failed to send notification to client", zap.Error(err))
			sub.finished <- true
			unsubscribes = append(unsubscribes, i)
		}
	}

	for _, unsub := range unsubscribes {
		subs = slices.Delete(subs, unsub, unsub)
	}
	if len(subs) > 0 {
		s.subscribers.Store(userID, subs)
	} else {
		s.subscribers.Delete(userID)
	}

	return nil
}
