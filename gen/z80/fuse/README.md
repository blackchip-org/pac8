# fuse

Z80 tests designed for the Free Unix Spectrum Emulator (FUSE).

Test files can be found in the resource pack at:

- fuse
    - tests.expected
    - <nolink>tests.in</nolink>

The original location of FUSE is here:

- http://fuse-emulator.sourceforge.net/

The files found in the resource pack were downloaded from:

- https://github.com/descarte1/fuse-emulator-fuse/tree/fuse-1-3-6/z80/tests

Generate `in.go` and `expected.go` with:

```bash
go generate
```