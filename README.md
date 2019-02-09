# pac8

Latest development work can be found here:

https://github.com/blackchip-org/retro-cs

[![Build Status](https://travis-ci.com/blackchip-org/pac8.svg?branch=master)](https://travis-ci.com/blackchip-org/pac8)

After finding the following document I decided to give it a try:

https://www.lomont.org/Software/Games/PacMan/PacmanEmulation.pdf

All projects need a name so the working title for this one is pac8, the Portable Arcade Cabinet, 8-bit edition.

## Requirements

Go and SDL2 are needed to build the application.

### Linux

Install SDL with:

```bash
sudo apt-get install libsdl2{,-image,-mixer,-ttf,-gfx}-dev
```

Install go from here:

https://golang.org/dl


### macOS

Install go and SDL with:

```bash
brew install go sdl2{,_image,_mixer,_ttf,_gfx} pkg-config
```

### ROMs

You will need to find [ROMs](ROMS.md) somewhere.

## Installation

```bash
go get github.com/blackchip-org/pac8
```

## Run

```bash
~/go/bin/pac8 -g <game>
```

where game is one of the following:

- `pacman`
- `mspacman`
- `galaga`


Use the `-m` flag to enable the [monitor](monitor.md).

## Inputs

- `c`: Coin slot
- `1`: One Player Start
- `2`: Two Player Start
- Arrow keys: Joystick

## Status

- Pacman & Ms. Pacman
  - Playable
  - Sound works but is a bit glitchy
  - High scores not saved
- Galaga
  - Work in progress
  - Boot to "PUSH START BUTTON" screen
  - IO registers stuffed to advance past startup
- z80
  - Failing two [zexdoc](component/proc/z80/internal/zex/README.md) tests

## License

MIT



