package utils

import (
	"go.uber.org/zap"
	"testing"
)

func TestNewLogger_returns_logger(t *testing.T) {
	type testCase struct {
		name    string
		wantNil bool
	}

	testCases := []testCase{
		{
			name:    "normal_case",
			wantNil: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := NewLogger()
			logger = AddLoggerFields(logger, "logger test")

			if (logger == nil) != tc.wantNil {
				t.Errorf("NewLogger() got = %v, wantNil = %v", logger, tc.wantNil)
			}

			if _, ok := interface{}(logger).(*zap.Logger); !ok {
				t.Error("NewLogger() did not return a *zap.Logger")
			}
		})
	}
}
