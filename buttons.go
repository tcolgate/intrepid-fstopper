package main

import (
	"machine"
	"math/bits"
	"time"
)

type button uint8

const (
	butTimePlus button = 1 << iota
	butTimeMinus
	butRun
	butFocus
	butCancel
	butSafelight
	butMode

	butDoesLongPress = butFocus | butMode
	butDoesHold      = butTimePlus | butTimeMinus
)

type buttonEventType uint8

const (
	butEventPress     buttonEventType = iota // a single press
	butEventLongPress                        // a long press, interpretted as one events
	butEventHold                             // a held press, interpretted as an ongoing event
)

// bufIntEvent holds information from the button interrupts
type butIntEvent struct {
	button
	status bool
}

type butEvent struct {
	button
	buttonEventType
}

// bufMgr translates button interrupt events into button UI events
type butMgr struct {
	intEvents <-chan butIntEvent
	events    chan<- butEvent

	downTimes [7]time.Time
}

func buttonToIndex(b button) int {
	return bits.TrailingZeros8(uint8(b))
}

func (m *butMgr) process(t time.Time) {
	var unset time.Time

	// process any buffered interrupt events
loop:
	for {
		select {
		case u := <-m.intEvents:
			butIndx := buttonToIndex(u.button)
			if butIndx < 0 || butIndx > 6 {
				// unknown button, should log this
				continue
			}

			switch u.status {
			case false: // buttonUp
				if m.downTimes[butIndx] == unset {
					break
				}
				switch {
				case 0 != (u.button & butDoesHold):
				case 0 != (u.button & butDoesLongPress):
					d := time.Since(m.downTimes[butIndx])
					if d > 1*time.Second {
						m.events <- butEvent{
							button:          u.button,
							buttonEventType: butEventLongPress,
						}
					} else {
						m.events <- butEvent{
							button:          u.button,
							buttonEventType: butEventPress,
						}
					}
				default:
				}

				m.downTimes[butIndx] = unset
			case true: // buttonDown
				m.downTimes[butIndx] = t
				switch {
				case 0 != (u.button & butDoesLongPress):
				case 0 != (u.button & butDoesHold):
					fallthrough // hold buttons should produce a press on quick press too
				default:
					m.events <- butEvent{
						button:          u.button,
						buttonEventType: butEventPress,
					}
				}
			}
		default:
			break loop
		}
	}

	// produce events for long press and hold
	for i, t := range m.downTimes {
		if t == unset {
			continue
		}
		d := time.Since(t)
		but := button(1 << i)
		switch {
		case 0 != (but & butDoesHold):
		case 0 != (but & butDoesLongPress):
			if d < (1 * time.Second) {
				continue
			}
			m.events <- butEvent{
				button:          but,
				buttonEventType: butEventLongPress,
			}
			m.downTimes[i] = unset
		}
	}
}

func (m *butMgr) Int(p machine.Pin) {
	ev := butIntEvent{
		button: pinToButton(p),
		status: p.Get(),
	}

	select {
	case butIntEventChan <- ev:
	default:
	}
}
