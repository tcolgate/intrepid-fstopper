package main

import "intrepidfstopper/num"

type testStripMethod uint8

const (
	testStripMethodAbs     testStripMethod = iota // each strip is the same
	testStripMethodCover                          // each step covers previous
	testStripMethodUncover                        // each step uncovers paper
)

type testStrip struct {
	method testStripMethod
	steps  uint8
	// we can take exposure settings for exposureSet exposures[0]
}

type exposureSet struct {
	baseTime  uint16 // Only one base time is ever configured
	isTest    uint8
	testStrip testStrip
	exposures [maxExposures]exposure
}

type exposure struct {
	// These are set by the user in printMode
	expUnit expUnit // What's the setting for this exposure
	colVals [3]int16

	// These are read by exposureMode
	colTime [3]uint16
}

func (es *exposureSet) adjustBaseTime(delta int16) bool {
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

func expUnitToS(b uint16, u expUnit, v int16) uint16 {
	switch u {
	case expUnitAbsolute:
		if v >= 0 {
			return b + uint16(v)
		}
		return b - uint16(v*-1)
	case expUnitPercent:
		off := ((int32(b) / 100) * int32(v))
		if v >= 0 {
			return b + uint16(off)
		}
		return b - uint16(off*-1)

	case expUnitHalfStop:
		wholeStops := v / 2
		needsHalfStop := (v % 2) != 0
		if wholeStops < 0 {
			wholeStops = wholeStops * -1
		}
		var adj num.Num

		if v >= 0 {
			adj = adj << wholeStops
		} else {
			if needsHalfStop {
				// easier for the math if we go down an extra
				// stop and back up one
				adj = adj >> (wholeStops + 1)
			} else {
				adj = adj >> (wholeStops)
			}
		}

		if needsHalfStop {
			adj = num.Mul(num.Num(adj), halfStop)
		}

		if v > 0 {
			return b + uint16(adj)
		} else {
			return b - uint16(adj)
		}
	case expUnitThirdStop:
		return b
	case expUnitTenthStop:
		return b
	default:
		return 0
	}
}

func expSToUnit(b uint16, u expUnit, s uint16) int16 {
	switch u {
	case expUnitAbsolute:
		return int16(int32(s) - int32(b))
	case expUnitPercent:
	default:
		return 0
	}
	return 0
}

// convExpUnit converts between different exposure units
// to give nicer UX when changing expUnit used
func convExpUnit(t, f expUnit, b uint16, v int16) int16 {
	s := expUnitToS(b, f, v)
	return expSToUnit(b, t, s)
}

func (es *exposureSet) cycleExpUnit(exp uint8, up bool) bool {
	og := es.exposures[exp].expUnit

	curr := int(og)
	if up {
		curr++
	} else {
		curr--
	}

	switch {
	case exp == 0 && curr < 0:
		curr = int(expUnitOff - 1)
	case exp == 0 && curr >= int(expUnitOff):
		curr = 0
	case curr < 0:
		curr = int(expUnitLast - 1)
	case curr >= int(expUnitLast):
		curr = 0
	}

	es.exposures[exp].expUnit = expUnit(curr)
	es.exposures[exp].colVals[0] = convExpUnit(expUnit(curr), og, es.baseTime, es.exposures[exp].colVals[0])
	es.exposures[exp].colVals[1] = convExpUnit(expUnit(curr), og, es.baseTime, es.exposures[exp].colVals[1])
	es.exposures[exp].colVals[2] = convExpUnit(expUnit(curr), og, es.baseTime, es.exposures[exp].colVals[2])

	return true
}

func (es *exposureSet) adjustExposureTime(exp uint8, col uint8, delta int16) bool {
	switch es.exposures[exp].expUnit {
	case expUnitOff, expUnitFreeHand:
		return false
	case expUnitAbsolute:
		es.exposures[exp].colVals[col] += delta
	default:
		es.exposures[exp].colVals[col] += 1
	}

	return true
}
