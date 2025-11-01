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

package main

import "intrepidfstopper/num"

type testStripMethod uint8

const (
	testStripMethodCover testStripMethod = iota // each step covers ppaper
	testStripMethodAbs                          // each strip is the same
)

var testMethodStrs = [3][]byte{
	[]byte(`cov`),
	[]byte(`abs`),
}

type testStrip struct {
	method   testStripMethod
	steps    uint8
	exposure exposure
}

type exposureSet struct {
	baseTime  uint16 // Only one base time is ever configured
	isTest    bool
	testStrip testStrip
	exposures [maxExposures]exposure
}

type exposure struct {
	// These are set by the user in printMode
	expUnit expUnit // What's the setting for this exposure
	colVal  int16
	rgb     [3]uint8
	enabled bool

	time uint16
}

func (es *exposureSet) adjustBaseTime(long, neg bool) bool {
	delta := int16(10)
	if long {
		delta = 100
	}
	if neg {
		delta *= -1
	}

	switch {
	case delta > 0:
		es.baseTime += uint16(delta)
		if es.baseTime >= 25500 {
			es.baseTime = 25500
		}
		return true
	case delta < 0:
		if es.baseTime < uint16(-1*delta) {
			es.baseTime = 0
		} else {
			es.baseTime -= uint16(delta * -1)
		}
		return true
	default:
		return false
	}
}

func bound(v int32) uint16 {
	switch {
	case v >= 600_00:
		return 600_00
	case v <= 0:
		return 0
	default:
		return uint16(v)
	}
}

func halfStops(b uint16, v int16) uint16 {
	if v == 0 {
		return b
	}

	neg := v < 0
	if neg {
		v = v * -1
	}

	adj := int32(b)

	i := v
	for {
		if i <= 1 {
			break
		}
		if !neg {
			adj = int32(uint16(adj) << 1)
			if adj > 600_00 {
				return 600_00
			}
		} else {
			adj = int32(uint16(adj) >> 1)
			if adj < 0 {
				return 0
			}
		}
		i -= 2
	}

	if i == 0 {
		return bound(adj)
	}

	if !neg {
		adj = int32(num.Mul(num.Num(adj), halfStop))
		if adj > 600_00 {
			return 600_00
		}
	} else {
		adj = int32(num.Mul(num.Num(adj), negHalfStop))
		if adj < 0 {
			return 0
		}
	}

	return bound(adj)
}

func thirdStops(b uint16, v int16) uint16 {
	if v == 0 {
		return b
	}

	neg := v < 0
	if neg {
		v = v * -1
	}

	adj := int32(b)

	i := v
	for {
		if i <= 2 {
			break
		}
		if !neg {
			adj = int32(uint16(adj) << 1)
			if adj > 600_00 {
				return 600_00
			}
		} else {
			adj = int32(uint16(adj) >> 1)
			if adj < 0 {
				return 0
			}
		}
		i -= 3
	}

	if i == 0 {
		return bound(adj)
	}

	for i = i; i > 0; i -= 1 {
		if !neg {
			adj = int32(num.Mul(num.Num(adj), thirdStop))
			if adj > 600_00 {
				return 600_00
			}
		} else {
			adj = int32(num.Mul(num.Num(adj), negThirdStop))
			if adj < 0 {
				return 0
			}
		}
	}

	return bound(adj)
}

func tenthStops(b uint16, v int16) uint16 {
	if v == 0 {
		return b
	}

	neg := v < 0
	if neg {
		v = v * -1
	}

	adj := int32(b)

	i := v
	for {
		if i <= 9 {
			break
		}
		if !neg {
			adj = int32(uint16(adj) << 1)
			if adj > 600_00 {
				return 600_00
			}
		} else {
			adj = int32(uint16(adj) >> 1)
			if adj < 0 {
				return 0
			}
		}
		i -= 10
	}

	if i == 0 {
		return bound(adj)
	}

	for i = i; i > 0; i -= 1 {
		if !neg {
			adj = int32(num.Mul(num.Num(adj), tenthStop))
			if adj > 600_00 {
				return 600_00
			}
		} else {
			adj = int32(num.Mul(num.Num(adj), negTenthStop))
			if adj < 0 {
				return 0
			}
		}
	}

	return bound(adj)
}

func expUnitToS(b uint16, u expUnit, v int16) uint16 {
	switch u {
	case expUnitAbsolute:
		return bound(int32(b) + int32(v))
	case expUnitPercent:
		return bound((int32(b) + (int32(b)/100)*int32(v)))
	case expUnitHalfStop:
		return halfStops(b, v)
	case expUnitThirdStop:
		return thirdStops(b, v)
	case expUnitTenthStop:
		return tenthStops(b, v)
	default:
		return 0
	}
}

func (es *exposureSet) cycleExpUnit(exp uint8, up bool) bool {
	expP := &es.exposures[exp]
	if es.isTest {
		expP = &es.testStrip.exposure
	} else {
		if !expP.enabled {
			expP.enabled = true
			return true
		}
	}

	og := expP.expUnit

	curr := int(og)
	if up {
		curr++
	} else {
		curr--
	}

	switch {
	case exp == 0 && curr < 0:
		curr = int(expUnitLast - 1)
	case exp == 0 && curr >= int(expUnitLast):
		curr = 0
	case curr < 0:
		curr = int(expUnitLast - 1)
	case curr >= int(expUnitLast):
		curr = 0
	}

	if es.isTest && curr == int(expUnitFreeHand) {
		if up {
			curr = 0
		} else {
			curr -= 1
		}
	}

	expP.expUnit = expUnit(curr)
	expP.colVal = 0

	return true
}

func (es *exposureSet) adjustExposureTime(exp uint8, long, neg bool) bool {
	// TODO: cap these values
	expP := &es.exposures[exp]
	if es.isTest {
		expP = &es.testStrip.exposure
	}

	var delta int16
	switch expP.expUnit {
	case expUnitPercent:
		delta = int16(5)
		if long {
			delta = 10
		}
	default:
		delta = int16(10)
		if long {
			delta = 100
		}
	}

	if neg {
		delta *= -1
	}

	switch expP.expUnit {
	case expUnitFreeHand:
		return false
	case expUnitAbsolute:
		expP.colVal += delta
	case expUnitPercent:
		expP.colVal += delta
		if expP.colVal < -95 {
			expP.colVal = -95
		}
	default:
		if delta > 0 {
			expP.colVal += 1
		} else {
			expP.colVal -= 1
		}
	}

	return true
}

func (es *exposureSet) tpAdjustExposureSet(touchPointIndex uint8, exp uint8, long, neg bool) bool {
	switch touchPointIndex {
	case 0: // baseTime
		return es.adjustBaseTime(long, neg)
	case 1: // exposure adjustment
		return es.adjustExposureTime(exp, long, neg)
	case 2: // adjustment unit
		return es.cycleExpUnit(exp, true)
	case 3: // test strip step count
		switch es.isTest {
		case true:
			if !neg {
				if es.testStrip.steps == 2 {
					return false
				}
				es.testStrip.steps++
			} else {
				if es.testStrip.steps == 0 {
					return false
				}
				es.testStrip.steps--
			}
		case false:
			return false
		}
		return true
	case 4: // test strip step method
		switch es.isTest {
		case true:
			if es.testStrip.method == 1 {
				es.testStrip.method = 0
			} else {
				es.testStrip.method = 1
			}
		case false:
			return false
		}
		return true
	default:
		return false
	}
}

func (es *exposureSet) calcTestInto(out *[maxExposures]int64) uint8 {
	allsteps := 1 + (es.testStrip.steps+1)*2

	v := (-1 * es.testStrip.exposure.colVal * (int16(es.testStrip.steps + 1)))

	for i := uint8(0); i < allsteps; i++ {
		out[i] = (int64)(expUnitToS(
			es.baseTime,
			es.testStrip.exposure.expUnit,
			v,
		)) * int64(tick)

		v += es.testStrip.exposure.colVal
	}

	switch es.testStrip.method {
	case testStripMethodCover:
		run := out[0]
		for i := uint8(1); i < allsteps; i++ {
			out[i] -= run
			run += out[i]
		}
	default:
		// testStripMethodAbs
	}

	return allsteps
}

func (es *exposureSet) calcInto(out *[maxExposures]int64) uint8 {
	if es.isTest {
		return es.calcTestInto(out)
	}

	expCnt := uint8(0)
expCnt:
	for i := range es.exposures {
		if !es.exposures[i].enabled {
			break expCnt
		}
		switch es.exposures[i].expUnit {
		case expUnitFreeHand:
		default:
			out[i] = (int64)(expUnitToS(
				es.baseTime,
				es.exposures[i].expUnit,
				es.exposures[i].colVal,
			)) * int64(tick)
		}
		expCnt++
	}
	return expCnt
}
