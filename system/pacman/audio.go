package pacman

import (
	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/component/audio"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	bufSize  = 367
	channels = 2
	nvoices  = 3
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
	SampleRate int
	Voices     [3]Voice
	ToneGen    [3]*audio.ToneGenerator
	buf        [3][]float32
	mix        []byte
	waveforms  [16][]float32
}

func NewAudio(spec sdl.AudioSpec, roms AudioROM) (*Audio, error) {
	a := &Audio{SampleRate: int(spec.Freq)}
	for i := 0; i < nvoices; i++ {
		a.ToneGen[i] = audio.NewToneGenerator(a.SampleRate)
		a.buf[i] = make([]float32, bufSize, bufSize)
	}
	mixLen := int(bufSize) * int(spec.Channels)
	a.mix = make([]byte, mixLen, mixLen)

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
		//console.Printf("freq: %4d, vol: %3d\n", freq, v.Vol)
	freq0 = a.toFreq(freq0)

	a.ToneGen[0].Freq = int(freq0) // 440
	a.ToneGen[0].Vol = float32(v0.Vol&0xf) / 15
	wf0 := bits.Slice(v0.Waveform, 0, 2)
	a.ToneGen[0].Waveform = a.waveforms[wf0]

	v1 := a.Voices[1]
	freq1 := (uint32(v1.Freq[0])&0x0f)<<4 +
		(uint32(v1.Freq[1])&0x0f)<<8 +
		(uint32(v1.Freq[2])&0x0f)<<12 +
		(uint32(v1.Freq[3])&0x0f)<<16
		//console.Printf("freq: %4d, vol: %3d\n", freq, v.Vol)
	freq1 = a.toFreq(freq1)

	a.ToneGen[1].Freq = int(freq1) // 440
	a.ToneGen[1].Vol = float32(v1.Vol&0xf) / 15
	wf1 := bits.Slice(v1.Waveform, 0, 2)
	a.ToneGen[1].Waveform = a.waveforms[wf1]

	v2 := a.Voices[2]
	freq2 := (uint32(v2.Freq[0])&0x0f)<<4 +
		(uint32(v2.Freq[1])&0x0f)<<8 +
		(uint32(v2.Freq[2])&0x0f)<<12 +
		(uint32(v2.Freq[3])&0x0f)<<16
		//console.Printf("freq: %4d, vol: %3d\n", freq, v.Vol)
	freq2 = a.toFreq(freq2)

	a.ToneGen[2].Freq = int(freq2) // 440
	a.ToneGen[2].Vol = float32(v2.Vol%0xf) / 15
	wf2 := bits.Slice(v2.Waveform, 0, 2)
	a.ToneGen[2].Waveform = a.waveforms[wf2]

	/*
		console.Printf("v[0]: %5d, %v, wf %2v\n", a.ToneGen[0].Freq, a.ToneGen[0].Vol, wf0)
		console.Printf("v[1]: %5d, %v, wf %2v\n", a.ToneGen[1].Freq, a.ToneGen[1].Vol, wf1)
		console.Printf("v[2]: %5d, %v, wf %2v\n", a.ToneGen[2].Freq, a.ToneGen[2].Vol, wf2)
		console.Println()
	*/

	/*

		a.ToneGen[0].Freq = 440
		a.ToneGen[0].Vol = 0.8
	*/

	q := sdl.GetQueuedAudioSize(1)
	// FIXME: For now, skip if something is in the queue
	if q > 0 {
		return nil
	}
	n := bufSize - int(q)
	//n := bufSize
	//console.Printf("q: %v, n: %v\n", q, n)
	a.ToneGen[0].Fill(a.buf[0], n)
	a.ToneGen[1].Fill(a.buf[1], n)
	a.ToneGen[2].Fill(a.buf[2], n)

	for i, d := 0, 0; i < n; i, d = i+1, d+2 {
		mix := (a.buf[0][i] + a.buf[1][i] + a.buf[2][i]) / 3.0
		//mix := a.buf[1][i]
		umix := uint16(((mix + 1) / 2) * float32(255))
		a.mix[d+0] = byte(umix)
		a.mix[d+1] = byte(umix)
	}
	return sdl.QueueAudio(1, a.mix)
}

func rescale(mem memory.Memory, addr uint16) []float32 {
	out := make([]float32, 32, 32)
	for i := uint16(0); i < 32; i++ {
		v := mem.Load(addr + i)
		out[i] = (float32(v) - 7.5) / 8
	}
	return out
}
