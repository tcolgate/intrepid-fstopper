package button

import (
	"math/bits"
	"time"
)

type Button uint8

const (
	Plus Button = 1 << iota
	Minus
	Run
	Focus
	Cancel
	Safelight
	Mode

	DoesLongPress = Focus | Mode | Cancel
	DoesHold      = Plus | Minus

	// not sure why, but these two are inverted
	inverted = Mode | Safelight

	LongPressTime = 1 * time.Second
	HoldPressTime = 100 * time.Millisecond
)

type EventType uint8

const (
	EventPress      EventType = iota // a single press
	EventLongPress                   // a long press, interpreted as one events
	EventHoldRepeat                  // a held press, interpreted as an ongoing event
)

// IntEvent holds information from the button interrupts
type IntEvent struct {
	Button
	Status bool
}

type Event struct {
	Button
	EventType
}

// Mgr translates button interrupt events into button UI events
type Mgr struct {
	IntEvents <-chan IntEvent
	Events    chan<- Event

	DownTimes     [7]int64
	lastSentTimes [7]int64
}

func toIndex(b Button) int {
	return bits.TrailingZeros8(uint8(b))
}

func (m *Mgr) Process(now int64) {
	// process any buffered interrupt events
loop:
	for {
		select {
		case u := <-m.IntEvents:
			butIndx := toIndex(u.Button)
			if butIndx > 6 {
				// unknown button, should log this
				continue
			}

			status := u.Status
			if 0 != (u.Button & inverted) {
				status = !status
			}

			switch status {
			case false: // Down
				m.DownTimes[butIndx] = now
				switch {
				case 0 != (u.Button & DoesLongPress):
				case 0 != (u.Button & DoesHold):
					fallthrough // hold buttons should produce a press on quick press too
				default:
					m.Events <- Event{
						Button:    u.Button,
						EventType: EventPress,
					}
					m.lastSentTimes[butIndx] = now
				}
			case true: // Up
				dt := m.DownTimes[butIndx]
				if dt == 0 {
					break
				}
				switch {
				case 0 != (u.Button & DoesHold):
				case 0 != (u.Button & DoesLongPress):
					d := now - dt
					if d > int64(LongPressTime) {
						m.Events <- Event{
							Button:    u.Button,
							EventType: EventLongPress,
						}
					} else {
						m.Events <- Event{
							Button:    u.Button,
							EventType: EventPress,
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
		case 0 != (but & DoesHold):
			if d < int64(LongPressTime-HoldPressTime) {
				continue
			}

			if (now - lst) < int64(HoldPressTime) {
				continue
			}
			m.Events <- Event{
				Button:    but,
				EventType: EventHoldRepeat,
			}
			m.lastSentTimes[i] = now
		case 0 != (but & DoesLongPress):
			if d < int64(LongPressTime) {
				continue
			}
			m.Events <- Event{
				Button:    but,
				EventType: EventLongPress,
			}
			m.DownTimes[i] = 0
			m.lastSentTimes[i] = now
		}
	}
}
