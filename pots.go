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

type potMgr struct {
	disabled [4]bool
	quant    [4]uint16
	last     [4]uint16
}

const (
	uint16Max = ^(uint16(0))
)

func (mgr *potMgr) SetDisabled(p uint8, b bool) {
	mgr.disabled[p] = b
}

func (mgr *potMgr) SetPotQuant(p uint8, q uint16) {
	mgr.quant[p] = q
}

func (mgr *potMgr) Process(t int64) {
	var updated uint8

	for i := range pots {
		if mgr.disabled[i] {
			continue
		}

		quant := uint16(128)
		if mgr.quant[i] != 0 {
			quant = mgr.quant[i]
		}
		step := (uint16Max) / quant
		newV := pots[i].Get() / step
		if newV != mgr.last[i] {
			updatePotBits := uint8(1) << i
			updated |= updatePotBits
			mgr.last[i] = newV
		}
	}

	if updated == 0 {
		return
	}

	update := potUpdate{
		updated: updated,
	}

	update.vals = mgr.last
	select {
	case potUpdateChan <- update:
	default:
	}
}
