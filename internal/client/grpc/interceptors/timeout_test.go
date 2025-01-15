package interceptors

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestTimeoutInterceptor(t *testing.T) {
	invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return nil
		}
	}

	timeoutInterceptor := Timeout(50 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := timeoutInterceptor(ctx, "testMethod", nil, nil, nil, invoker)
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded error, but got %v", err)
	}
}
