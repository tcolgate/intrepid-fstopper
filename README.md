# Intrepid Enlarger Alternative Firmware

This repo contains an experiment to create an alternative firmware
for the Intrepid Enlarger timer.

*NOTE* I make not promises that this will not destroy your timer,
burn your house down, or kill a random kitten. That said, it appears
to be safe to flash the timer, and then re-flash with the official
firmware, I am not aware of resultant dead kittens in my area.

## Achievments

Not many so far, mostly just trying to implement the needed button behaviour
and work out the internal API

- Only supporing the P variant right now (I have two timers, one old, one new,
  both are P. B may be easy to support, may just need a separate build)
- Basic control of the light, with asingle exposure
  - +/- set the time, you can tap, or long press
  - Run starts exposure
  - During exposure Run pauses exposure, Run restarts
  - Cancel stop the exposure
- Focus mode:
  - Single press Focus for red light
  - Long press Focus for white light
  - Long press Focus during focus switches red/white
  - Single press Focus or Cancel during focus mode exits to previous mode
- F-Stop timing for printing
  - adjust by 1/2,1/3,1/10 stops
  - Percentage from base time, rather than stops
- F-Stop timing for test strips
  - 7,5 or 3 strip patch test strips
  - gradual covering or fully exposing of each test patch

Controls and Display are not great at the moment, but there are lots of ways
to improve things.

## Goals

The intention is to provide:

- Multi-exposure
- Freehand exposure
- Tri-colour printing (with f-stop timing)

Possible additions
- Pulse a burst of red every second during freehand exposure to
  help counting time (these would not contribute to exposure time,
  and is intended to be an alternate to an audible blip (in colour
  mode that would be turning the lamp off.
- During BW print, when exposure is paused, the panel could switch
  to red light

## Non-Goals

- Colour or multi-grade LEDs. Intrepid put lots of effort in to
  calibrating the multigrade and colour LEDs. The LEDs themselves
  have quite broad spectral ranges which makes accurate colour
  calibration difficult, if not really possible. As such I am
  focused on white light usage only.
- I may be lying about the previous point. An RGB exposure mode
  and while it may not be accurate, it will be more "honest" than
  the existing CMY mode. It would be interesting to compare to
  the tri-colour nad filtered options.

## Programming Style

The firmware is implemented using TinyGo, but due to the constraints of the
atmega chips, only a strict subset of Go's functionality can be used.

- Channels are OK, occasionally useful (essentially a pre-existing, generic
  ring buffer.
- Slices are fine
- I have avoided maps (avoided the need for any of the key hashing logic)
- No interfaces, raw storage is fine, but the runtime required to use them
  blows the image size.
- No Go routines (I did use them with success, but they did not provide big win,
  and caused latency, and were hard to debug when they broke stuff, which was
  trivial to do)
- No in-function heap allocation, minimize stack usage, no GC
- Almost no stdlib usage (all attempts resulted in blowing out image size)
- Minimize function call stack depth
- If you want to be able to test code you need to be in a dedicated package, without
  reference to the `machine` package, and should run under a standalone go test on the
  build host.
- THere's a fair bit of global state, I've hidden it in types that should be testable, but
  then these get used as singletons.
