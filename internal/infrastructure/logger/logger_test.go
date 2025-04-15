package logger_test

import (
	"backend/internal/infrastructure/logger"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want *logger.Logger
	}{
		{
			name: "new logger",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := logger.NewLogger()
			if got == nil {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
