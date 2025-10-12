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
