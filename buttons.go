package main

import (
	"machine"
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
	butEventPress    buttonEventType = iota // a single press
	butEventLongPres                        // a long press, interpretted as one events
	butEventHold                            // a held press, interpretted as an ongoing event
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
	tickChan  <-chan struct{}
	intEvents <-chan butIntEvent
	events    chan<- butEvent

	downTimes [7]time.Time
}

func buttonToIndex(b button) int {
	indx := 0
	for {
		if b == 1 {
			return indx
		}
		indx += 1
		b >>= 1
		if b == 0 {
			return 0
		}
	}
}

func (m *butMgr) process(t time.Time) {
	var unset time.Time
	for {
		select {
		case u := <-m.intEvents:
			switch {
			case 0 != (u.button & butDoesHold):
				fallthrough
			case 0 != (u.button & butDoesLongPress):
				fallthrough
			default: // button has simple press behaviour
				butIndx := buttonToIndex(u.button)
				if butIndx < 0 || butIndx > 6 {
					continue
				}
				switch u.status {
				case false: // buttonUp
					d := time.Since(m.downTimes[butIndx])
					print("d ", d)
					m.downTimes[butIndx] = unset
				case true: // buttonDown
					m.downTimes[butIndx] = t
					m.events <- butEvent{
						button:          u.button,
						buttonEventType: butEventPress,
					}
				}
			}
		default:
			return
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
