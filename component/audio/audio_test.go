package audio

import (
	"fmt"
	"testing"

	. "github.com/blackchip-org/pac8/expect"
)

func TestFill(t *testing.T) {
	v := NewVoice(8)
	v.Freq = 2
	v.Vol = 1.0
	v.Waveform = []float64{-1, 0, 1, 0}
	buf := make([]float64, 8, 8)
	v.Fill(buf, len(buf))

	WithFormat(t, "%.2f").
		Expect(buf).ToBe([]float64{-1, 0, 1, 0, -1, 0, 1, 0})
}

func TestFillHalfVol(t *testing.T) {
	v := NewVoice(8)
	v.Freq = 2
	v.Vol = 0.5
	v.Waveform = []float64{-1, 0, 1, 0}
	buf := make([]float64, 8, 8)
	v.Fill(buf, len(buf))

	WithFormat(t, "%.2f").
		Expect(buf).ToBe([]float64{-0.5, 0, 0.5, 0, -0.5, 0, 0.5, 0})
}

func TestFillStretch(t *testing.T) {
	v := NewVoice(8)
	v.Freq = 1
	v.Vol = 1.0
	v.Waveform = []float64{-1, 0, 1, 0}
	buf := make([]float64, 8, 8)
	v.Fill(buf, len(buf))

	WithFormat(t, "%.2f").
		Expect(buf).ToBe([]float64{-1, -1, 0, 0, 1, 1, 0, 0})
}

func TestConvertSample(t *testing.T) {
	tests := []struct {
		from float64
		to   uint16
	}{
		{-1, 0},
		{1, 0xffff},
		{0, 0x7fff},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v to %v", test.from, test.to), func(t *testing.T) {
			WithFormat(t, "%04x").Expect(convert(test.from)).ToBe(test.to)
		})
	}
}
