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
		name     string
		expected Num
		a        Num
		b        int32
	}{
		{"", 1_230, 1_230, 100},
		{"", 123, 1_230, 10},
		{"", 12_300, 1_230, 1000},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := Mul(tt.a, tt.b)
			if actual != tt.expected {
				t.Errorf("(%d * %d): expected %d, actual %d", tt.a, tt.b, tt.expected, actual)
			}
		})
	}
}
