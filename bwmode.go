package main

import "intrepidfstopper/num"

type bwMode struct {
	prevMode Mode

	baseTime           uint64
	paused             bool
	running            bool
	remainingTime      int64
	exposureFactor     int32
	exposureFactorUnit expUnit
}

func (e bwMode) SwitchTo(prev Mode) {
	e.prevMode = prev
}

func (e bwMode) SwitchAway() Mode {
	return e.prevMode
}

func (e bwMode) TouchPoints() []touchPoint {
	return touchPoints[0]
}

func (e bwMode) Tick(passed int64) (bool, bool) {
	return false, false
}

func (e bwMode) Run() bool {
	return true
}

func (e bwMode) Focus() bool {
	return true
}

func (e bwMode) LongFocus() bool {
	return true
}

func (e bwMode) Cancel(touchPoint uint8) (bool, bool) {
	// should reset stuff and/or delete the current
	// exposure
	return false, false
}

func (e bwMode) Plus(touchPointIndex uint8) bool {
	switch touchPointIndex {
	case 0:
		if e.baseTime != 25500 {
			e.baseTime += 10
		}
	case 1:
		if e.exposureFactor != 126 {
			e.exposureFactor += 1
		}
	case 2:
		e.exposureFactorUnit++
		if e.exposureFactorUnit > 4 {
			e.exposureFactorUnit = 0
		}
	default:
		return false
	}
	return true
}

func (e bwMode) LongPlus(touchPointIndex uint8) bool {
	switch touchPointIndex {
	case 0:
		if e.baseTime != 25500 {
			e.baseTime += 10
		}
	case 1:
		if e.exposureFactor != 126 {
			e.exposureFactor += 1
		}
	case 2:
		e.exposureFactorUnit++
		if e.exposureFactorUnit > 4 {
			e.exposureFactorUnit = 0
		}
	default:
		return false
	}
	return true
}

func (e bwMode) Minus(touchPointIndex uint8) bool {
	switch touchPointIndex {
	case 0:
		if e.baseTime != 0 {
			e.baseTime -= 10
		}
	case 1:
		if e.exposureFactor != -126 {
			e.exposureFactor -= 1
		}
	case 2:
		e.exposureFactorUnit--
		if e.exposureFactorUnit == 0 {
			e.exposureFactorUnit = 4
		}
	default:
		return false
	}
	return true
}

func (e bwMode) LongMinus(touchPointIndex uint8) bool {
	switch touchPointIndex {
	case 0:
		if e.baseTime != 0 {
			e.baseTime -= 10
		}
	case 1:
		if e.exposureFactor != -126 {
			e.exposureFactor -= 1
		}
	case 2:
		e.exposureFactorUnit--
		if e.exposureFactorUnit < 0 {
			e.exposureFactorUnit = 4
		}
	default:
		return false
	}
	return true
}

func (e bwMode) UpdateDisplay(nextDisplay *[2][]byte) *touchPoint {
	nb := &num.NumBuf{}
	copy(nextDisplay[0], stringTable[0][0])
	copy(nextDisplay[1], stringTable[0][1])

	num.Out(nb, num.Num(e.baseTime))
	copy(nextDisplay[0][2:6], nb[0:4])

	if e.exposureFactor < 0 {
		nextDisplay[1][1] = signMinus
	} else {
		nextDisplay[1][1] = signPlus
	}
	absExpFact := e.exposureFactor
	if absExpFact < 0 {
		absExpFact = absExpFact * -1
	}
	num.OutLeft(nb, num.Num(absExpFact))
	copy(nextDisplay[1][2:6], nb[0:4])

	copy(nextDisplay[0][10:15], expUnitNames[e.exposureFactorUnit][0:4])

	res := num.Num(11_23)
	num.Out(nb, num.Num(res))
	resLen := num.Len(res)
	nextDisplay[1][13-resLen] = []byte("(")[0]
	nextDisplay[1][14-resLen] = []byte("=")[0]
	copy(nextDisplay[1][11:15], nb[0:4])
	nextDisplay[1][15] = []byte(")")[0]

	return nil
}
