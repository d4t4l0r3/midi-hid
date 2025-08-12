# midi-hid

This software allows mapping and translating of MIDI commands to HID inputs on Linux.

## Installation and Usage

Install with

```bash
go install git.datalore.sh/datalore/midi-hid@latest
```

and run it with `midi-hid`. If no config file is specified, it reads `~/.config/midi-hid/config.yaml`.
Every configured midi controller will be represented by a virtual gamepad, and the inputs will be translated until SIGINT is received.

## Configuration

See `example-config.yaml` For an annotated configuration file.
The valid names for buttons are:

| Button name  | Description                      |
| ------------ | -------------------------------- |
| `north`      | E.g. Y on XBox or triangle on PS |
| `east`       | E.g. B on XBox or circle on PS   |
| `south`      | E.g. A on XBox or X on PS        |
| `west`       | E.g. X on XBox or square on PS   |
| `l1`         | Left bumper                      |
| `l2`         | Left trigger                     |
| `l3`         | Left stick pressed down          |
| `r1`         | Right bumper                     |
| `r2`         | Right trigger                    |
| `r3`         | Right stick pressed down         |
| `select`     | Select, or Back on XBox          |
| `start`      | Start button                     |
| `dpad-up`    | Directional pad up               |
| `dpad-down`  | Directional pad down             |
| `dpad-left`  | Directional pad left             |
| `dpad-right` | Directional pad right            |

Valid axis are `left-x`, `left-y`, `right-x` and `right-y`.

### Finding the MIDI channel and note / controller

Many vendors provide official documentation on the MIDI commands sent by their devices, but if you are unable to find them, you can use `aseqdump -p <port>` to print all MIDI messages sent by your device. Simply interact with the controls you want to map and you will see the corresponding messages.

## Third-party libraries

 - <https://github.com/bendahl/uinput>
 - <https://gitlab.com/gomidi/midi>
 - <https://github.com/charmbracelet/log>
