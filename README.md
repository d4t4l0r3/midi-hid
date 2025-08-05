# midi-hid

This software allows mapping and translating of MIDI commands to HID inputs on Linux.

## Known issues

The midi library used seems to recognise NoteOff messages as NoteOn messages. However, they can still be recognised by checking the velocity, which is always 0 in NoteOff messages. A workaround has been implemented.
