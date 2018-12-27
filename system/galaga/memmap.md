# galaga

## Files

| Name | Length | Description |
|-|-|-|
| `04m_g01.bin` | `$1000` | Code, CPU 0, #1
| `04k_g02.bin` | `$1000` | Code, CPU 0, #2
| `04j_g03.bin` | `$1000` | Code, CPU 0, #3
| `04h_g04.bin` | `$1000` | Code, CPU 0, #4
| `04e_g05.bin` | `$1000` | Code, CPU 1
| `04d_g06.bin` | `$1000` | Code, CPU 2
| `07m_g08.bin` | `$1000` | Tile set
| `07e_g10.bin` | `$1000` | Sprites, #1
| `07h_g09.bin` | `$1000` | Sprites, #2

## ROMs

| Address | CPU0 | CPU1 | CPU2 |
|-|-|-|-|
| `0000 - 0fff` | `04m_g01.bin` | `04e_g05.bin` | `04d_g06.bin` |
| `1000 - 1fff` | `04k_g02.bin` | | |
| `2000 - 2fff` | `04j_g03.bin` | | |
| `3000 - 3fff` | `04h_g04.bin` | | |

## References

- MAME source code, version 0.37b5
- MAME source code, latest