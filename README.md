# midi-hid

This software allows mapping and translating of MIDI commands to HID inputs on Linux.

## Installation and Usage

Install with

```bash
go install git.datalore.sh/datalore/midi-hid@latest
```

and run it with `midi-hid`. If no config file is specified, it reads `~/.config/midi-hid/config.yaml`.
Every configured midi controller will be represented by a virtual gamepad, and the inputs will be translated until SIGINT is received.

See the provided example config on how to configure your controller, it should be pretty self-explanatory.

## Known issues

The midi library used seems to recognise NoteOff messages as NoteOn messages. However, they can still be recognised by checking the velocity, which is always 0 in NoteOff messages. A workaround has been implemented.

## Third-party libraries

 - <https://github.com/bendahl/uinput>
 - <https://gitlab.com/gomidi/midi>
