package machine

type CharDecoder func(uint8) (rune, bool)

var AsciiDecoder = func(code uint8) (rune, bool) {
	printable := code >= 32 && code < 128
	return rune(code), printable
}
