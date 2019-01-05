package pacman

import (
	"github.com/blackchip-org/pac8/pkg/audio"
	"github.com/blackchip-org/pac8/pkg/memory"
	"github.com/blackchip-org/pac8/pkg/util/bits"
	"github.com/veandco/go-sdl2/sdl"
)

type AudioROM struct {
	R1 memory.Memory
	R2 memory.Memory
}

type Voice struct {
	Acc      [5]uint8
	Waveform uint8
	Freq     [5]uint8
	Vol      uint8
}

type Audio struct {
	Synth     *audio.Synth
	Voices    [3]Voice
	waveforms [16][]float64
}

func NewAudio(spec sdl.AudioSpec, roms memory.Set) (*Audio, error) {
	a := &Audio{}
	synth, err := audio.NewSynth(spec, 3)
	if err != nil {
		return nil, err
	}
	a.Synth = synth
	for i := 0; i < 16; i++ {
		addr := uint16(i * 32)
		a.waveforms[i] = rescale(roms["waveform"], addr)
	}
	return a, nil
}

func (a *Audio) Queue() error {
	for i := 0; i < 3; i++ {
		v := a.Voices[i]
		wf := bits.Slice(v.Waveform, 0, 2)

		// Voice 0 has 5 bytes but Voice 1 and 2 only have 4 bytes with
		// the missing lower byte being zero.
		nFreq := 4
		if i == 0 {
			nFreq = 5
		}
		a.Synth.V[i].Freq = freq(v.Freq, nFreq)
		a.Synth.V[i].Vol = float64(v.Vol&0xf) / 15
		a.Synth.V[i].Waveform = a.waveforms[wf]
	}
	return a.Synth.Queue()
}

func rescale(mem memory.Memory, addr uint16) []float64 {
	out := make([]float64, 32, 32)
	for i := uint16(0); i < 32; i++ {
		v := mem.Load(addr + i)
		out[i] = (float64(v) - 7.5) / 8
	}
	return out
}

func freq(f [5]uint8, n int) int {
	val := uint32(0)
	shift := uint(0)
	if n == 4 {
		shift = 4
	}
	for i := 0; i < n; i++ {
		val += uint32(f[i]&0x0f) << shift
		shift += 4
	}
	freq := (375.0 / 4096.0) * float32(val)
	return int(freq)
}
