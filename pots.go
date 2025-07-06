package main

const (
	conPotUpdated = 1 << iota
	cyanPotUpdated
	magentaPotUpdated
	yellowPotUpdated
)

type potUpdate struct {
	vals    [4]uint16
	updated uint8
}

func potChanged(o, n uint16) bool {
	// This gives us 128 valid pot positions
	minDiff := uint16(128)

	if o > n {
		return (minDiff < (o - n))
	}

	if o < n {
		return (minDiff < (n - o))
	}

	return false
}


type potMgr struct {
	lastConV    uint16
	lastCyanV   uint16
	lastMagentV uint16
	lastYellowV uint16
}

func (mgr *potMgr) Process(t int64) {
	var updated uint8

	if newConV := contrast.Get(); potChanged(newConV, mgr.lastConV) {
		updated |= conPotUpdated
		mgr.lastConV = newConV
	}

	if newCyanV := cyan.Get(); potChanged(newCyanV, mgr.lastCyanV) {
		updated |= cyanPotUpdated
		mgr.lastCyanV = newCyanV
	}

	if newMagentaV := magenta.Get(); potChanged(newMagentaV, mgr.lastMagentV) {
		updated |= magentaPotUpdated
		mgr.lastMagentV = newMagentaV
	}

	if newYellowV := yellow.Get(); potChanged(newYellowV, mgr.lastYellowV) {
		updated |= yellowPotUpdated
		mgr.lastYellowV = newYellowV
	}

	if updated == 0 {
		return
	}

	update := potUpdate{
		updated: updated,
	}

	update.vals[0] = mgr.lastConV
	update.vals[1] = mgr.lastCyanV
	update.vals[2] = mgr.lastMagentV
	update.vals[3] = mgr.lastYellowV

	select {
	case potUpdateChan <- update:
	default:
	}
}
