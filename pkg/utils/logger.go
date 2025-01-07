// Package utils предоставляет различные вспомогательные функции и инструменты,
// используемые в проекте GophKeeper. В данном пакете реализован функционал
// для создания и конфигурирования логгера на базе библиотеки go.uber.org/zap.
//
// Основная цель пакета — облегчить процесс инициализации журнала (логгирования)
// и обеспечить корректную обработку и вывод диагностической информации в
// продакшн-окружении.
package utils

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"syscall"
)

// NewLogger создает новый экземпляр логгера *zap. Logger в конфигурации,
// предназначенной для использования в продакшн-окружении GophKeeper.
// Время в сообщениях логгера форматируется в стиле ISO8601 (e.g. 2024-01-02T15:04:05Z).
//
// При ошибках в процессе инициализации логгера (например, неверная конфигурация),
// функция выполняет zap.L().Fatal, что приводит к завершению работы приложения.
//
// Возвращаемое значение:
//
//	*zap. Logger – сконфигурированный экземпляр логгера, готовый к записи логов.
//
// Особенности реализации:
//  1. Создается zap.NewProductionConfig(), где задается используемый формат логов
//     и дополнительные настройки для продакшн-среды.
//  2. Кодирует метку времени в формате ISO8601 для удобства чтения.
//  3. При завершении функции зарегистрирован отложенный вызов (defer) logger.Sync(),
//     обеспечивающий сброс буферов логгирования. Ошибки, связанные с некорректными
//     файловыми дескрипторами (EBADF, ENOTTY, EINVAL), игнорируются; остальные –
//     обрабатываются как фатальные.
func NewLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		zap.L().Fatal("Server failed to create logger instance", zap.Error(err))
	}

	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil && (!errors.Is(err, syscall.EBADF) &&
			!errors.Is(err, syscall.ENOTTY) &&
			!errors.Is(err, syscall.EINVAL)) {
			logger.Fatal("Server failed on Sync()-method call on zap.Logger.", zap.Error(err))
		}
	}(logger)

	return logger
}
