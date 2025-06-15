package main

type potMgr struct {
	lastConV    uint16
	lastCyanV   uint16
	lastMagentV uint16
	lastYellowV uint16
}

func (mgr *potMgr) process(t int64) {
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
