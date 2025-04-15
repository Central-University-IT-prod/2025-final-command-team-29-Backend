package models_test

import (
	"backend/internal/models"
	"testing"
)

func TestCondition_Scan(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		value   interface{}
		wantErr bool
	}{
		{
			name:    "valid condition",
			value:   []byte(`{"type": "buy", "buy": 2}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var a models.Condition
			gotErr := a.Scan(tt.value)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Scan() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Scan() succeeded unexpectedly")
			}
		})
	}
}

func TestCondition_Value(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		wantErr bool
	}{
		{
			name: "valid value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var a models.Condition
			_, gotErr := a.Value()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Value() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Value() succeeded unexpectedly")
			}
		})
	}
}
