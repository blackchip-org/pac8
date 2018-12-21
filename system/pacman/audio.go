package pacman

import (
	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/component/audio"
	"github.com/blackchip-org/pac8/component/memory"
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

func NewAudio(spec sdl.AudioSpec, roms AudioROM) (*Audio, error) {
	a := &Audio{}
	synth, err := audio.NewSynth(spec, 3)
	if err != nil {
		return nil, err
	}
	a.Synth = synth
	for i := 0; i < 8; i++ {
		addr := uint16(i * 32)
		a.waveforms[i+0] = rescale(roms.R1, addr)
		a.waveforms[i+8] = rescale(roms.R2, addr)
	}
	return a, nil
}

func (a *Audio) toFreq(v uint32) uint32 {
	f := (375.0 / 4096.0) * float32(v)
	return uint32(f)
}

func (a *Audio) Queue() error {
	v0 := a.Voices[0]
	freq0 := uint32(v0.Freq[0])&0x0f +
		(uint32(v0.Freq[1])&0x0f)<<4 +
		(uint32(v0.Freq[2])&0x0f)<<8 +
		(uint32(v0.Freq[3])&0x0f)<<12 +
		(uint32(v0.Freq[4])&0x0f)<<16
	freq0 = a.toFreq(freq0)
	wf0 := bits.Slice(v0.Waveform, 0, 2)
	a.Synth.V[0].Freq = int(freq0)
	a.Synth.V[0].Vol = float64(v0.Vol&0xf) / 15
	a.Synth.V[0].Waveform = a.waveforms[wf0]

	a.Synth.V[0].Freq = 440
	a.Synth.V[0].Vol = 1
	a.Synth.V[0].Waveform = a.waveforms[0]

	/*
		v1 := a.Voices[1]
		freq1 := (uint32(v1.Freq[0])&0x0f)<<4 +
			(uint32(v1.Freq[1])&0x0f)<<8 +
			(uint32(v1.Freq[2])&0x0f)<<12 +
			(uint32(v1.Freq[3])&0x0f)<<16
		freq1 = a.toFreq(freq1)
		wf1 := bits.Slice(v1.Waveform, 0, 2)
		a.Synth.V[1].Freq = int(freq1)
		a.Synth.V[1].Vol = float64(v1.Vol&0xf) / 15
		a.Synth.V[1].Waveform = a.waveforms[wf1]

		v2 := a.Voices[2]
		freq2 := (uint32(v2.Freq[0])&0x0f)<<4 +
			(uint32(v2.Freq[1])&0x0f)<<8 +
			(uint32(v2.Freq[2])&0x0f)<<12 +
			(uint32(v2.Freq[3])&0x0f)<<16
		freq2 = a.toFreq(freq2)
		wf2 := bits.Slice(v2.Waveform, 0, 2)
		a.Synth.V[2].Freq = int(freq2)
		a.Synth.V[2].Vol = float64(v2.Vol%0xf) / 15
		a.Synth.V[2].Waveform = a.waveforms[wf2]
	*/

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
