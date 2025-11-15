# User Manual

## Installation

In order to flash the device you will need a 
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
* Download the Fstopper firmware from the releases page
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
## Usage

### Basic Controls

- *Contrast dial*: The Grey "contrast" control is the main means of navigating
  the UI. As you turn it a cursor will move between the various settings. If
  you are in the Print Exposure screem, or Test Strip screen there are two
  pages of settings. The second page control exposure brightness or colour.
  These settings will be discussed later. The three C/M/Y control wheel
  are currently unused. You can turn them if that kind of thing gives you a
  kick.
- *+/-*: These controls will variously increas/descrease values, or cycle
  between settings. In some contexts a long press will increase/decrease values
  faster.
- *Cancel*: The cancel button is context dependent.
  - In most context it will reset the value under the cursor
  - If the cursor is on the Exposure Unit setting or the Expsure Number
    setting, Cancel will turn the currently selected exposure on/off
  - During an exposure or focus Cancel will stop everything and return you to
    the relevant Print or Test mode.
- *Mode*: The Mode button has two functions. A quick press alternates between
  Print Exposure mode and Test Strip mode. A long press changes between the
  White light and RGB LED modes.
- *Run*:
  - In Print or Test Strip modes the Run button will begin an exposure
  - When running an exposure the Run button will pause. This will turn off the
    light. You can press again to restart the exposure. If you press it by accident
    there is no need to panic, you can just press again and finish the exposure. The
    LED warm up and cool down are fast enough that this will no adversely impact your
    exposure
- The Focus button can be pressed while on the Print or Test Strip screens. By
  this will turn on the Red focus light. Pressing again returns you to the
  previous mode. If you want a White light you can long hold the Focus button.
  Long hold works from the Print, Test Strip and Focus screens.


### Print Exposure settings

When you first power on the unit you will be dropped into the Print Exposure
mode. This is the main mode used for exposing prints.

- *Base Time*
- *Exposure Value*
- *Exposure Unit*

| Unit | Interpretation of Exposure Value   |
| ---- | ---------------------------------- |
| s    | Fixed additional amount of seconds |
| /2º  | Number of half stops               |
| /3º  | Number of third stops              |
| /10º | Number of tenth stops              |
| %    | Percentage of base time            |
| Free | A free-hand exposure               |

### Test Strip Screen

A quick press of the Mode button switch from Print exposure mode to the Test
Strip exposure mode.

- *Base Time*
- *Exposure Value*
- *Exposure Unit*

### Light colour control


