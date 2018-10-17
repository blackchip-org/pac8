# dasm

Generates the disassembly table using the "Full Z80 Opcode List Including
Undocumented Opcodes" written by J.G. Harston.

The document is parsed from the following file in the resource pack:

- data/harston
    - z80oplist.txt

Document downloaded from here:

- http://www.z80.info/z80oplist.txt

Expected values from the disassembler are found in `expected.txt` and
are in the following two line format:

- List of bytes in hexadecimal separated by a space
- Expected output of the disassembler

Each test is delimited by a blank line.

Generate `dasm.go` and `harston.go` with:

```bash
go generate
```
