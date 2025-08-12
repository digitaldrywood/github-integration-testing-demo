package tests

import (
	"testing"
)

func TestBasicMath(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"add positive", 2, 3, 5},
		{"add negative", -2, -3, -5},
		{"add mixed", -2, 3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a + tt.b
			if result != tt.expected {
				t.Errorf("Expected %d + %d = %d, got %d", tt.a, tt.b, tt.expected, result)
			}
		})
	}
}

func TestStringOperations(t *testing.T) {
	t.Run("concatenation", func(t *testing.T) {
		result := "Hello, " + "World!"
		expected := "Hello, World!"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("length", func(t *testing.T) {
		str := "Test"
		if len(str) != 4 {
			t.Errorf("Expected length 4, got %d", len(str))
		}
	})
}

func TestSliceOperations(t *testing.T) {
	t.Run("append", func(t *testing.T) {
		slice := []int{1, 2, 3}
		slice = append(slice, 4)
		if len(slice) != 4 || slice[3] != 4 {
			t.Errorf("Append failed: %v", slice)
		}
	})

	t.Run("capacity", func(t *testing.T) {
		slice := make([]int, 3, 5)
		if cap(slice) != 5 {
			t.Errorf("Expected capacity 5, got %d", cap(slice))
		}
	})
}