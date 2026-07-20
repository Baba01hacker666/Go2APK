package main

import "testing"

func TestCalculate(t *testing.T) {
	tests := []struct {
		left, op, right, want string
	}{
		{"2", "+", "3", "5"},
		{"7", "−", "9", "-2"},
		{"6", "×", "7", "42"},
		{"8", "÷", "2", "4"},
		{"10", "%", "3", "1"},
		{"5", "÷", "0", "Cannot divide by 0"},
	}
	for _, tt := range tests {
		got, err := Calculate(tt.left, tt.op, tt.right)
		if err != nil {
			t.Fatalf("Calculate(%q, %q, %q) unexpected error: %v", tt.left, tt.op, tt.right, err)
		}
		if got != tt.want {
			t.Fatalf("Calculate(%q, %q, %q) = %q, want %q", tt.left, tt.op, tt.right, got, tt.want)
		}
	}
}
