# Intrepid Enlarger Alternative Firemware

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

Controls and Display are not great at the moment, but there are lots of ways
to improve things.

## Goals

The intention is to provide:

- F-Stop timing for printing and test strips
  - covering,uncovering or fully exposing on each section of a test strip
  - adjust by 1/2,1/3 or 1/10 stops
  - 9,7,5 or 3 strip  test strips (could be arbitrary, but that's not that useful)
- Multi-exposure
- Tri-colour printing (with f-stop timing)
- Freehand exposure

Possible additions
- Percentage from base time, rather than stops

## Non-Goals

- Colour or multi-grade LEDs. Intrepid put lots of effort in to
  calibrating the multigrade and colour LEDs. The LEDs themselves
  have quite broad spectral ranges which makes accurate colour
  calibration difficult, if not really possible. As such I am
  focused on white light usage only.
