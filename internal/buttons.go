package internal

import (
	"math/bits"
	"time"
)

type Button uint8

const (
	ButTimePlus Button = 1 << iota
	ButTimeMinus
	ButRun
	ButFocus
	ButCancel
	ButSafelight
	ButMode

	ButDoesLongPress = ButFocus | ButMode
	ButDoesHold      = ButTimePlus | ButTimeMinus

	// not sure why, but these two are inverted
	butInverted = ButMode | ButSafelight

	ButLongPressTime = 1 * time.Second
	ButHoldPressTime = 100 * time.Millisecond
)

type ButtonEventType uint8

const (
	ButEventPress      ButtonEventType = iota // a single press
	ButEventLongPress                         // a long press, interpretted as one events
	ButEventHoldRepeat                        // a held press, interpretted as an ongoing event
)

// ButIntEvent holds information from the button interrupts
type ButIntEvent struct {
	Button
	Status bool
}

type ButEvent struct {
	Button
	ButtonEventType
}

// BufMgr translates button interrupt events into button UI events
type ButMgr struct {
	IntEvents <-chan ButIntEvent
	Events    chan<- ButEvent

	DownTimes     [7]int64
	lastSentTimes [7]int64
}

func buttonToIndex(b Button) int {
	return bits.TrailingZeros8(uint8(b))
}

func (m *ButMgr) Process(now int64) {
	// process any buffered interrupt events
loop:
	for {
		select {
		case u := <-m.IntEvents:
			butIndx := buttonToIndex(u.Button)
			if butIndx > 6 {
				// unknown button, should log this
				continue
			}

			status := u.Status
			if 0 != (u.Button & butInverted) {
				status = !status
			}

			switch status {
			case false: // Down
				m.DownTimes[butIndx] = now
				switch {
				case 0 != (u.Button & ButDoesLongPress):
				case 0 != (u.Button & ButDoesHold):
					fallthrough // hold buttons should produce a press on quick press too
				default:
					m.Events <- ButEvent{
						Button:          u.Button,
						ButtonEventType: ButEventPress,
					}
					m.lastSentTimes[butIndx] = now
				}
			case true: // Up
				dt := m.DownTimes[butIndx]
				if dt == 0 {
					break
				}
				switch {
				case 0 != (u.Button & ButDoesHold):
				case 0 != (u.Button & ButDoesLongPress):
					d := now - dt
					if d > int64(ButLongPressTime) {
						m.Events <- ButEvent{
							Button:          u.Button,
							ButtonEventType: ButEventLongPress,
						}
					} else {
						m.Events <- ButEvent{
							Button:          u.Button,
							ButtonEventType: ButEventPress,
						}
					}
					m.lastSentTimes[butIndx] = now
				default:
				}

				m.DownTimes[butIndx] = 0
			}
		default:
			break loop
		}
	}

	// produce events for long press and hold
	var dt int64
	var lst int64
	for i := range len(m.DownTimes) {
		dt = m.DownTimes[i]
		lst = m.lastSentTimes[i]
		if dt == 0 {
			continue
		}
		d := now - dt
		but := Button(1 << i)
		switch {
		case 0 != (but & ButDoesHold):
			if d < int64(ButLongPressTime-ButHoldPressTime) {
				continue
			}

			if (now - lst) < int64(ButHoldPressTime) {
				continue
			}
			m.Events <- ButEvent{
				Button:          but,
				ButtonEventType: ButEventHoldRepeat,
			}
			m.lastSentTimes[i] = now
		case 0 != (but & ButDoesLongPress):
			if d < int64(ButLongPressTime) {
				continue
			}
			m.Events <- ButEvent{
				Button:          but,
				ButtonEventType: ButEventLongPress,
			}
			m.DownTimes[i] = 0
			m.lastSentTimes[i] = now
		}
	}
}
