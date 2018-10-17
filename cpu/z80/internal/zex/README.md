# zex

The Z80 Instruction Exerciser written by Frank D. Cringle.

Source code and compiled binaries can be found in the resource pack at:

- data/zex
    - <nolink>zexall.com</nolink>
    - zexall.z80
    - <nolink>zexdoc.com</nolink>
    - zexdoc.z80

The original location of the exerciser seems to be here:

- http://mdfs.net/Software/Z80/Exerciser/

The sources found in the resource pack were downloaded from:

- https://github.com/anotherlin/z80emu/blob/master/testfiles

Helpful references:

- https://floooh.github.io/2016/07/12/z80-rust-ms1.html
- http://jeffavery.ca/computers/macintosh_z80exerciser.html

Run the functional test with:

```bash
go test -v -tags fn
```

The test can take up to ten minutes to complete. The following tests
are currently failing:

- `cpd<r>`
- `<rlca,rrca,rla,rra>`

Run the benchmarks with:

```bash
go test -run=X -tags=fn -bench=.
```
