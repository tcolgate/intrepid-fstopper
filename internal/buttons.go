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
)

type ButtonEventType uint8

const (
	ButEventPress     ButtonEventType = iota // a single press
	ButEventLongPress                        // a long press, interpretted as one events
	ButEventHold                             // a held press, interpretted as an ongoing event
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

	DownTimes [7]int64
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

			switch u.Status {
			case false: // buttonUp
				dt := m.DownTimes[butIndx]
				if dt == 0 {
					break
				}
				switch {
				case 0 != (u.Button & ButDoesHold):
				case 0 != (u.Button & ButDoesLongPress):
					d := now - dt
					if d > int64(1*time.Second) {
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
				default:
				}

				m.DownTimes[butIndx] = 0
			case true: // buttonDown
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
				}
			}
		default:
			break loop
		}
	}

	// produce events for long press and hold
	var dt int64
	for i := range len(m.DownTimes) {
		dt = m.DownTimes[i]
		if dt == 0 {
			continue
		}
		d := now - dt
		but := Button(1 << i)
		switch {
		case 0 != (but & ButDoesHold):
		case 0 != (but & ButDoesLongPress):
			if d < int64(1*time.Second) {
				continue
			}
			m.Events <- ButEvent{
				Button:          but,
				ButtonEventType: ButEventLongPress,
			}
			m.DownTimes[i] = 0
		}
	}
}
