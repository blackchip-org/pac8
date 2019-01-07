# zex

The Z80 Instruction Exerciser written by Frank D. Cringle.

Source code and compiled binaries can be found in the resource pack at:

- zex
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
go test -v -tags fn -timeout 60m
```

Running the full zexdoc can take more than 10 minutes. This test instead breaks up each test into an individual run. The HL register is loaded with the address of the test and the program counter is set to the beginning of the normal test loop. Execution is stopped when the program counter returns to the top of the loop. Output is then checked for "ERROR" to determine if the test passes or fails.

The following tests are currently failing:

- TestZexdoc/cpd1: `cpd<r>`
- TestZexdoc/rot8080: `<rlca,rrca,rla,rra>`

Run the benchmarks with:

```bash
go test -run=X -tags=fn -bench=.
```
