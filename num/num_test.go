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
				t.Errorf("(%d): expected %s, actual %s", tt.given, tt.expected, actual)
			}
		})
	}
}
