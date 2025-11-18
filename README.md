# Intrepid Enlarger Alternative Firmware

This repo contains an experiment to create an alternative firmware
for the Intrepid Enlarger timer.

*NOTE* I make not promises that this will not destroy your timer,
burn your house down, or kill a random kitten. That said, it appears
to be safe to flash the timer, and then re-flash with the official
firmware, I am not aware of resultant dead kittens in my area.

The main features of this firmware are:

- Multiple exposures: You can set up to 9 exposures for a print, each with
  separate duration controlled relative to a fixed base time. Each exposure
  can have a different brightness or colour.
- Test strips can be created with the exposure varied by a fixed amount of
  seconds, stops or percentages.
- A safer Focus light than the original, making it harder to accidentally turn
  on the white light, and easier to turn it on when you want it.

## What is F-Stop printing?

The process of f-stop printing can be thought of as allowing us to treat our timer
more like the "shutter speed" dial of our camera, and talk about our paper
exposure more in terms of the Exposure Triangle we are used to with film. It
was pioneered by [Gene Nocon](https://www.youtube.com/watch?v=xoAiBNSpg6Y&pp=ygUPZi1zdG9wIHByaW50aW5n)
and is implemented in many commercially available darkroom enlarger timers.

F-Stop printing can be achieved with regular darkroom timers by manually calculating
the times for individual steps, but this process is time consuming and requires
a fair bit of maths and/or referencing large tables of times. It's particularly
fiddly with the Intrepid timer as we have to `+/-` our timer settings.

F-Stop timing is particularly useful for producing test strips. A traditional
test strip uses linear time spaces between each step. Since we typically think
of print density exponentially, rather than linearly, one end
of the test strip will have much bigger steps in density than the other. With
an f-stop test strip we can get can a spread of densities more like a range
of zones. This can make it easier to quickly establish a more accurate
highlight exposure.

Additional advantages of f-stop printing include:
- easier contrast grading. A 1/3rd stop of exposure time is approximately
  equivalent to a 1/2 filter grade on the Ilford filter kit. When looking
  at an 1/3 step f-stop test strip we can estimate that if our highlights
  are good at 1 step over the base time, and the shadows look good at
  2 steps under, we have 3 * 1/3rd stop difference which we can translate
  to a 1.5 difference in filter grade (3 x 0.5).
- Easier print resizing. If you need to dodge and burn areas of your print
  you can time those in terms of stops or percentages relative to your
  base exposure time. When resizing a print you can calculate your new base
  exposure time, and then use the same percentage, or stops for dodges and
  burns (since they are relative to the base time). This makes it easier to
  translate a test print at a smaller size into any other print size.

In short, f-stop printing lets us more easily establish our print time,
and understand the global contrast. It can also let us think about our local
contrast control steps in a way that is agnostic of the size of the print. For instance
you can state to burn in for "60% of the base time", rather than "10s for an 8x10".

## Usage

The physical Intrepid Enlarger Timer was not designed with f-stop timing in
mind. Adding features results in more fiddly control than the original firmware
and will take a bit of getting used to.

The controls are as follow
- *Contrast dial*: The Grey "contrast" control dial is the main means of
  navigating the UI. As you turn it a cursor will move between the various
  settings. Note that there are two pages of settings, the first controls
  exposure timing, the second controls the settings for the LED.
- *+/-*: These controls will variously increas/descrease values, or cycle
  between settings. Holding the button will change settings faster
- *Cancel*: The cancel button is context dependent.
- *Mode*: The Mode button has two functions. A quick press alternates between
  Print Exposure mode and Test Strip mode. A long press changes between the
  White light and RGB LED modes.
- *Run*: Is used to start or continue an exposure. It can also pause/unpause a
  running exposure.
- The Focus button behaves slightly differently to the default Intrepid timer.

### Focusing light

The Focus button works differently from the original Intrepid firmware. The behaviour
is intended to make it easier, and safer to switch from Red to White light.

When on either the Print or Test Strip screens you can press the Focus button to
switch on the focus light for focusing your print. A short press of the button
will turn on the Red light. A long press will instead turn on the White light.
If you have the red light on you can long hold Focus to switch between white and
red light.

A short press of either Focus or Cancel while the focus light is on will turn
off the light and return you to whichever mode were first in.

### Print Exposure settings

When you first power on the unit you will be dropped into the Print Setup
screen. This is the main screen used for controlling prints.

![The Print Setup screen](/doc/screenshots/print-1.svg)

There are four main settings for controlling a print.

- *Base Time*: This setting is shared between every exposure of a
  multi-exposure print and shared with the Base Time setting of the Test Strip
  mode.
- *Exposure Value*: This controls the quantity of the exposure unit that will
  be applied to the base time to calculate the final time for an exposure.
- *Exposure Unit*: This controls how each individual exposure time will be
  calculated from the Base Time and Exposure Value.
   *s*: Fixed additional number of seconds
   */2º*: Number of half stops
   */3º*: Number of third stop
   */10º*: Number of tenth stops
   *%*: Percentage of base time
   *Free*: A free-hand exposure
- *Multi-exposure selection*: THe timer allows you to setup up to 9 exposures
  for a print. The *+/-* buttons allow you to edit a different exposure. This
  display will be explained further in elsewhere.

The default Exposure Unit of *s* treats the timer as a traditional timer. When you
turn the unit on you can just set the base time, hit *Run* and the light will
come on for the set amount of time. Changing the exposure value will add or subtract
a fixed amount of time from the *Base Time*.

Changing to the  */2º* Exposure Unit is where more interesting things happen. If we set
our Base time to 10s, our Exposure Value to 2, and our Exposure Unit to */2º* we will see
the following display.

![Examples of different exposure unit adjustments](/doc/screenshots/print-2.svg)

We have asked the time to use "base time plus 2 half stops", 2 half stops is a full
stop, which is twice the exposure. The timer shows us that the final calculated time
is 20s , twice the base time. If we had set the Exposure Value to *-2*, we are asking
for "base time minus 2 half stops", which would be 5s. Similarly we could ask
for "base time minus 10% of the base time", which would be 9 seconds.

The *FreeHand* exposure unit ignores the base time and allows you to expose a step for
as long as you wish.

### Running the exposures

When you press *Run* on the Print Screen (or Test Strip Screen), the exposure of your print will
begin.
* If you only have 1 exposure configured the print will begin immediately.
* If you have multiple exposures configured the exposure will start but will be
  immediately paused. You will need to press *Run* a second time to begin the
  first exposure.
* You can press *Run* again at during a running exposure to pause the exposure, pressing
  Run again will continue the exposure. The LED warm-up, and cool-down is essentially
  negligible (in the order of 1us), so we do not need to be concerned about this skewing
  our exposure times.
* If you are using a FreeHand exposure the Run button will still let you pause and restart
  . To stop a FreeHand exposure you must press the Cancel button. This turn off the light
  and advance you to the next exposure.
* Pressing Cancel at any other time will stop this exposure, and return you to
  the Print or Test Strip screens.  Any remaining configured exposures are
  ignored.

### Multiple Exposures

THe timer allows you to setup up to 9 exposures for a print. By pressing *+/-* on
the exposure number (shown directly after the *E:*), you can cycle between
the available exposures.
* Initially only exposure 1 is enabled.
* Exposure 1 cannot be disabled
* The other 8 can be turned on by changing the Exposure Unit from
  *Off* to something else, or toggled on/off with the *Cancel* button.
* When the control is on the active exposure selection the display
  of the number of active exposure switch to a *9* to remind you
  how many exposures you can choose from (this behaviour may change).

Each exposure can have different Exposure Value, Exposure Unit, and LED
settings.

The display shows which exposure you are adjusting the settings for, and the
total number of configured exposures for this print. For example, *1/1* means
there is only one exposure configured. *3/4* means you are editing the 3rd
exposure and have 4 exposures enabled in total. It is possible for the display
to read *5/2*, indicating that you are editing the 5th of the 9 available
exposures, but only have 2 of the 9 exposures enabled.

### Test Strip Screen

A quick press of the Mode button switches from Print Setup mode to the Test
Strip  mode. There are four controls on the Test Strip screen.

![The Test Strip Setup Screen](/doc/screenshots/ts-1.svg)

- *Base Time*: this setting is shared with the Print Exposure settings. If you
  change it here, it updates the Base Time of the Print Exposure too.
- *Exposure Unit*: This is the unit by which we will adjust the test strip steps
- *Exposure Value*: This is the number of Exposure Units that should be added
  and subtracted per-step
- *Step Count*: This is a visual representation of the number of steps. It can
  be +/- 1, 2 or 3 steps (respectively 3, 5 or 7 steps in total). The Central
  step will always be exposed at the base time.
- Method: This chooses how the strip will be exposed. There are two choices
  - *cov*: This indicates that the user will cover the previous steps. For the
         first exposure the strip is completely uncovered. For the first step
         a section of the strip is covered, then first two sections, and so on.
  - abs: This indicates that each section will get a full exposure. This is for
         use with step wedge printers that present a small window to the user that
         a strip is slid under and moved along for each exposure.

The Exposure Units are the same as the Print Screen. On the final printed test strip the
central exposure will always be exposed at the base time. The step to the left will then
be reduce by the exposure unit and value. The steps 2 and 3 to the left will be reduced by 2 and 3
times the value. The steps to the right will be increased by 1, 2 and 3 times the value. The following
shows the seconds of exposure for each given setting:

* `10s ± 2 s`: `4 6 8 10 12 14 16` - each step is 2s apart
* `10s ± 2 /2º`: `1.25 2.5 5 10 20 40` - each step is 1 stop apart (two half stops)
* `10s ± 20 %`: `4 6 8 10 12 14 16` - each stp is 20% of 10s apart (2 seconds)

### Light colour control

On the Print screen the letter to the left of the *E:* for exposure indicates
which of the two colour modes you are in. *W* for white and *C* for colour.

On power on the unit will use white light at full brightness.

To control the light settings can use the grey control dial to scroll onto the
second page of controls. In White Light mode this will show a *Brightness*
setting that can be set between 1 and 255.

By long pressing the Mode button you can toggle between White light and RGB
controls. The second page of setting on the Print and Test Strip screens will
switch to offering settings for individual R, G and B channels, settable from 0
to 255.

The same method of light control will apply to all exposures in a multi-exposure
print.

### White Light

In white mode each exposure can set a different brightness for each exposure in a
multi-exposure print.

For the default white light mode you must use filters for contrast and colour
control. For contrast control you can use regular contrast control filters.
Since the white LED output is not the same as a traditional tungsten bulb,
results will vary slightly from the traditional usage (this is also true of the
standard firmware). You can also use traditional CMY filters for traditional
RA4 printing, or use RGB filters and try Tri-Colour printing for RA4.

The brightness control of the white light can be useful to act as a form of ND
filter to increase exposure times without external NDs or needing to change
aperture. If can also be used as a brightness control if you are using the
light for sensitometry.

### RGB Light

In RGB mode you can set exact R,G & B values for each exposure in a
multi-exposure print.

The LEDs in the enlarger do not have particularly "narrow" R,G and B spectra.
The G and B are quite broad. The blue in particular leans a little greener than
is really desirable for RA4 printing.

In Intrepid's original firmware they present a traditional CMY filtration
interface to the user. This translates CMY values into RGB for the light. That
translation combines with the inherent inaccuracy of the LEDs to make finer
grained colour control trickier than it might be.

Rather than attempt to provide a convenient CMY (or contrast graded) filter
interface, I have opted to instead just provide direct control of the RGB
channels. In time I hope to provide a guide to contrast grading and RA4 printing
directly with the RGB light values.

*NOTE*: it is a quirk of the hardware that *R=255 G=255 B=255* is not a white
light as you might expect.

## Installation

*WARNING*: Obviously I cannot provide any guarantees of safety of the processes
here. I have personally flashed both my Intrepid Compact and 4x5 enlarger
timers with this process literally hundred of times with no impact.

* You will need a USB Mini cable
* If you have your timer plugged in *UNPLUG THE POWER SUPPLY*. During all the
  following steps you should only connect and power the device via the USB Mini
  connection. It is important that you *DO NOT* have both the power supply and
  USB Mini cables connected at the same time.
* Make sure your timer is [updated to the most up to date Intrepid firmware](https://intrepidcamera.co.uk/blogs/guides/upgrading-to-the-new-enlarger-firmware)
* Check you controller has the correct chip set. With the original Intrepid
  firmware loaded:
  - Hold the Safe Light and Focus buttons while the timer boots.
  - *If it shows "V1.2a" you have the A variant of the timer AND WILL NOT
    BE ABLE TO USE THIS FIRMWARE AT THIS TIME*
  - If it shows "V1.2" (and  *not* "V1.2a"), you have the B variant and
    are fine to proceed
* Download and install [avrdude](https://github.com/avrdudes/avrdude). This is
  the tool that the official Intrepid firmware upload tool uses in the
  background.
* Backup your existing firmware using the following command. The `/dev/ttyUSB0`
  used here may need to be changed. There are some tips [here](https://www.ladyada.net/learn/avr/avrdude.html) under the information for the `-P` setting
  On Linux:
```
 avrdude -c stk500v1 -p m328p -b 57600 -P /dev/ttyUSB0 -U flash:r:"download-intrepid-firmware.hex":i
```
  This will write the existing firmware into a file called `download-intrepid-firmware.hex`. You should save this file for safe keeping.
* Download the `.hex` file from the latest intrepid-fstopper [releases page](https://github.com/tcolgate/intrepid-fstopper/releases/latest)
* Upload the new firmware
```
avrdude -c stk500v1 -p m328p -b 57600 -P /dev/ttyUSB0 -U intrepid-fstopper.hex
```
* If you wish to revert to the original firmware you can either re-run the
  official firmware update, or use the following command with the firmware you
  backed up from you device earlier
```
avrdude -c stk500v1 -p m328p -b 57600 -P /dev/ttyUSB0 -U download-intrepid-firmware.hex
```
