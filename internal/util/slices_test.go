package util

import "testing"

func TestStringIn(t *testing.T) {
	tests := []struct {
		name   string
		ss     []string
		target string
		want   bool
	}{
		{"found", []string{"a", "b", "c"}, "b", true},
		{"not found", []string{"a", "b", "c"}, "d", false},
		{"empty slice", nil, "a", false},
		{"empty target", []string{"a", "b"}, "", false},
		{"single match", []string{"x"}, "x", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringIn(tt.ss, tt.target); got != tt.want {
				t.Errorf("StringIn(%v, %q) = %v, want %v", tt.ss, tt.target, got, tt.want)
			}
		})
	}
}
