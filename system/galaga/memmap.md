# galaga

## Files

| Name | Length | Description |
|-|-|-|
| `04m_g01.bin` | `$1000` | Code, 1 of 4, CPU 1
| `04k_g02.bin` | `$1000` | Code, 2 of 4, CPU 1
| `04j_g03.bin` | `$1000` | Code, 3 of 4, CPU 1
| `04h_g04.bin` | `$1000` | Code, 4 of 4, CPU 1
| `04e_g05.bin` | `$1000` | Code, CPU 2
| `04d_g06.bin` | `$1000` | Code, CPU 3
| `07m_g08.bin` | `$1000` | Tile set
| `07e_g10.bin` | `$1000` | Sprites, 1 of 2
| `07h_g09.bin` | `$1000` | Sprites, 2 of 2

## ROMs

| Address | CPU1 | CPU2 | CPU3 |
|-|-|-|-|
| `0000 - 0fff` | `04m_g01.bin` | `04e_g05.bin` | `04d_g06.bin` |
| `1000 - 1fff` | `04k_g02.bin` | | |
| `2000 - 2fff` | `04j_g03.bin` | | |
| `3000 - 3fff` | `04h_g04.bin` | | |

## References

- MAME source code, version 0.37b5
- MAME source code, latest