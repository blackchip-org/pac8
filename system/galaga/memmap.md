# galaga

## Files

| Name | Length | Description |
|-|-|-|
| `04m_g01.bin` | `$1000` | Code, CPU 1, #1
| `04k_g02.bin` | `$1000` | Code, CPU 1, #2
| `04j_g03.bin` | `$1000` | Code, CPU 1, #3
| `04h_g04.bin` | `$1000` | Code, CPU 1, #4
| `04e_g05.bin` | `$1000` | Code, CPU 2
| `04d_g06.bin` | `$1000` | Code, CPU 3
| `07m_g08.bin` | `$1000` | Tile set
| `07e_g10.bin` | `$1000` | Sprites, #1
| `07h_g09.bin` | `$1000` | Sprites, #2

## ROMs

| Address | CPU1 | CPU2 | CPU3 |
|-|-|-|-|
| `0000 - 0fff` | `04m_g01.bin` | `04e_g05.bin` | `04d_g06.bin` |
| `1000 - 1fff` | `04k_g02.bin` | | |
| `2000 - 2fff` | `04j_g03.bin` | | |
| `3000 - 3fff` | `04h_g04.bin` | | |

### ROM, CPU 1


## RAM

| Address | Description |
|-|-|
| `8000 - 83ff` | Video RAM |
| `8400 - 87ff` | Color RAM |
| `9100` | CPU 2 notify CPU 1 (ROM check, more?)
| `9101` | CPU 3 notify CPU 1 (ROM check, more?)

## Strings (CPU 1)

Type: "R", reversed string. "LP", length prefixed.

| Address | Type | Contents |
|-|-|-|
| `00eb - 00fb` | R  | '1UP  HIGHSCORE'
| `09ca - 09cf` | R  | 'CREDIT'
| `09d0 - 09da` | R  | 'FREE PLAY'
| `3283 - 3297` |    | 'ENTER YOUR INITIALS !'
| `329b - 32aa` |    | 'SCORE  NAME'
| `32af - 32b3` |    | 'TOP 5'
| `32b7 - 32c4` |    | 'SCORE  NAME'
| `32c7 - 32ee` |    | Default high score table?
| `3349 - 335b` |    | 'THE GALACTIC HEROES'
| `3363 - 3368` |    | 'BEST 5'
| `3a4a - 3a4e` | LP | 'SOUND'
| `3ad3 - 3ad9` | LP | 'UPRIGHT'
| `3add - 3ae1` | LP | 'TABLE'
| `3ae7 - 3aea` | LP | 'RANK'
| `3aee - 3af2` | LP | 'SHIPS'
| `3af7 - 3afa` | LP | 'COIN'
| `3afe - 3b03` | LP | 'CREDIT'
| `3b07 - 3b0f` | LP | 'FREE PLAY'
| `3b21 - 3b2a` | LP | '1ST BONUS'
| `3b2d - 3b35` | LP | '0000 PTS'
| `3b39 - 3b41` | LP | '2ND BONUS'
| `3b44 - 3b4c` | LP | '0000 PTS'
| `3b50 - 3b58` | LP | 'AND EVERY'
| `3b5b - 3b63` | LP | '0000 PTS'
| `3b67 - 3b7d` | LP | 'BONUS NOTHING'
| `3b81 - 3b87` | LP | 'RAM OK'
| `3b8b - 3b91` | LP | 'ROM OK'


## References

- MAME source code, version 0.37b5
- MAME source code, latest