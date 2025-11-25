// Copyright 2025 Tristan Colgate-McFarlane
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func TestIntOut(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    Num
	}{
		{"", " 123", 123},
		{"", "9999", IntMax}, // TODO(tcm): could round up here
		{"", "9900", 9_900},
		{"", "5000", 5_000},
		{"", " 160", 160},
		{"", "   0", 0_00},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var b NumBuf
			IntOut(&b, tt.given)
			actual := string(b[:])
			if actual != tt.expected {
				t.Errorf("(%d): expected %q, actual %q", tt.given, tt.expected, actual)
			}
		})
	}
}

func TestIntOutLeft(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    Num
	}{
		{"", "123 ", 123},
		{"", "9999", IntMax}, // TODO(tcm): could round up here
		{"", "5000", 5000},
		{"", "99  ", 99},
		{"", "160 ", 160},
		{"", "0   ", 0_00},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var b NumBuf
			IntOutLeft(&b, tt.given)
			actual := string(b[:])
			if actual != tt.expected {
				t.Errorf("(%d): expected %q, actual %q", tt.given, tt.expected, actual)
			}
		})
	}
}

func TestIntLen(t *testing.T) {
	var tests = []struct {
		name     string
		expected int
		given    Num
	}{
		{"", 4, 1_230},
		{"", 4, IntMax},
		{"", 4, 9_900},
		{"", 3, 500},
		{"", 2, 16},
		{"", 1, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := IntLen(tt.given)
			if actual != tt.expected {
				t.Errorf("(%d): expected %d, actual %d", tt.given, tt.expected, actual)
			}
		})
	}
}

func TestMul(t *testing.T) {
	var tests = []struct {
		name            string
		expected        int32
		expectedRounded bool
		expectedErr     int32

		a Num
		b int32
	}{
		{"1230 * 1", 1_230, false, 0, 1_230, 1000},
		{"1230 * 0.1", 123, false, 0, 1_230, 100},
		{"1230 * 0.5", 615, false, 0, 1_230, 500},
		{"1231 * 0.5", 616, true, 0, 1_231, 500},
		{"1230 * 10", 12_300, false, 0, 1_230, 10000},

		{"16s + 1/2 stop", 22_63, false, -1, 1_600, HalfStop},
		{"16s - 1/2 stop", 11_31, false, 0, 1_600, NegHalfStop},
		{"16s + 1/3rd stop", 20_17, false, -1, 1_600, ThirdStop},
		{"16s - 1/3rd stop", 12_69, false, 1, 1_600, NegThirdStop},
		{"16s + 1/10th stop", 17_15, false, 0, 1_600, TenthStop},
		{"16s - 1/10th stop", 14_93, true, 0, 1_600, NegTenthStop},

		{"300s + 1/2 stop", 424_26, false, -6, 300_00, HalfStop},
		{"300s - 1/2 stop", 212_13, false, -3, 300_00, NegHalfStop},
		{"300s + 1/3rd stop", 377_98, false, 2, 300_00, ThirdStop},
		{"300s - 1/3rd stop", 238_11, false, 9, 300_00, NegThirdStop},
		{"300s + 1/10th stop", 321_53, false, 7, 300_00, TenthStop},
		{"300s - 1/10th stop", 279_91, false, -1, 300_00, NegTenthStop},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, rounded := Mul1000th(tt.a, tt.b)
			diff := actual - tt.expected
			if diff != tt.expectedErr {
				t.Errorf("(%d * %d): got %d, expected err of %d, actual %d", tt.a, tt.b, actual, tt.expectedErr, diff)
			}
			if rounded != tt.expectedRounded {
				t.Errorf("(%d * %d): expect rounded %v, actual %v", tt.a, tt.b, tt.expectedRounded, rounded)
			}
		})
	}
}
