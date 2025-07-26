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
	isTest    bool
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
		adj = int32(num.Mul(num.Num(adj), halfStop))
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
		return bound((int32(b) / 100) * int32(v))
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
	es.exposures[exp].colVals = [3]int16{0, 0, 0}

	return true
}

func (es *exposureSet) adjustExposureTime(exp uint8, col uint8, delta int16) bool {
	// TODO: cap these values

	switch es.exposures[exp].expUnit {
	case expUnitOff, expUnitFreeHand:
		return false
	case expUnitAbsolute, expUnitPercent:
		es.exposures[exp].colVals[col] += delta
	default:
		if delta > 0 {
			es.exposures[exp].colVals[col] += 1
		} else {
			es.exposures[exp].colVals[col] -= 1
		}
	}

	return true
}
