# pac8

[![Build Status](https://travis-ci.com/blackchip-org/pac8.svg?branch=master)](https://travis-ci.com/blackchip-org/pac8)

After finding the following document I decided to give it a try:

https://www.lomont.org/Software/Games/PacMan/PacmanEmulation.pdf

All projects need a name so the working title for this one is pac8, the Portable Arcade Cabinet, 8-bit edition.

## Requirements

Go and SDL2 are needed to build the application.

### macOS

Install go and SDL with:

```bash
brew install go sdl2{,_image,_mixer,_ttf,_gfx} pkg-config
```

### ROMs

Download the ROMs from [somewhere](http://example.org). Unpack the
tarball into your home directory.

## Installation

```bash
go get github.com/blackchip-org/pac8
```

## Run

```bash
~/go/bin/pac8
```

Use the `-m` flag to enable the [monitor](monitor.md).

## Inputs

- `5`: Coin slot 1
- `6`: Coin slot 2
- `1`: One Player Start
- `2`: Two Player Start
- Arrow keys: Joystick

## Status

Game is playable with some bugs. Hexadecimal score being one of those bugs.

## License

MIT



