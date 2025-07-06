package num

import "testing"

func TestOut(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    Num
	}{
		{"", " 123", 123_00},
		{"", " 655", Max}, // TODO(tcm): could round up here
		{"", "99.0", 99_00},
		{"", "50.0", 50_00},
		{"", " 1.6", 1_60},
		{"", " 0.0", 0_00},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var b NumBuf
			Out(&b, tt.given)
			actual := string(b[:])
			if actual != tt.expected {
				t.Errorf("(%d): expected %q, actual %q", tt.given, tt.expected, actual)
			}
		})
	}
}

func TestOutLeftJust(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    Num
	}{
		{"", "123 ", 123_00},
		{"", "655 ", Max}, // TODO(tcm): could round up here
		{"", "99.0", 99_00},
		{"", "50.0", 50_00},
		{"", "1.6 ", 1_60},
		{"", "0.0 ", 0_00},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var b NumBuf
			OutLeft(&b, tt.given)
			actual := string(b[:])
			if actual != tt.expected {
				t.Errorf("(%d): expected %q, actual %q", tt.given, tt.expected, actual)
			}
		})
	}
}

func TestLen(t *testing.T) {
	var tests = []struct {
		name     string
		expected int
		given    Num
	}{
		{"", 3, 123_00},
		{"", 3, Max}, // TODO(tcm): could round up here
		{"", 4, 99_00},
		{"", 4, 50_00},
		{"", 3, 1_60},
		{"", 3, 0_00},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := Len(tt.given)
			if actual != tt.expected {
				t.Errorf("(%d): expected %q, actual %q", tt.given, tt.expected, actual)
			}
		})
	}
}
