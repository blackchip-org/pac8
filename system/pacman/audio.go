package pacman

// typedef unsigned char Uint8;
// void Callback(void *userdata, Uint8 *stream, int len);
import "C"

import (
	"math"
	"reflect"
	"unsafe"

	"github.com/blackchip-org/pac8/console"
	"github.com/veandco/go-sdl2/sdl"
)

const bufSize = 800

type Audio struct {
	SamplesPerSecond int32
	ToneHz           int32
	ToneVolume       int
}

type voice struct {
	Acc      [5]uint8
	Waveform uint8
	Freq     [5]uint8
	Vol      uint8
}

var voices [3]voice

func NewAudio(reg *Registers) (*Audio, error) {
	a := &Audio{
		SamplesPerSecond: 48000,
		ToneHz:           440,
		ToneVolume:       300,
	}
	spec := sdl.AudioSpec{
		Freq:     a.SamplesPerSecond,
		Format:   sdl.AUDIO_U8,
		Channels: 2,
		Samples:  765,
		Callback: sdl.AudioCallback(C.Callback),
		UserData: unsafe.Pointer(&voices),
	}
	var got sdl.AudioSpec
	if err := sdl.OpenAudio(&spec, &got); err != nil {
		return nil, err
	}
	console.Printf("%+v\n", got)
	sdl.PauseAudio(false)
	return a, nil
}

func (a *Audio) Queue() error {
	/*
		//sdl.ClearQueuedAudio(1)
		v := a.voices[1]
		freq := uint32(v.Freq[0])&0x0f +
			(uint32(v.Freq[1])&0x0f)<<4 +
			(uint32(v.Freq[2])&0x0f)<<8 +
			(uint32(v.Freq[3])&0x0f)<<12 +
			(uint32(v.Freq[4])&0x0f)<<16
		vol := v.Vol
		vol = 15
		freq = 440
		console.Printf("freq: %v\n", freq)
		console.Printf("v0 %+v\n", a.voices[0])
		console.Printf("v1 %+v\n", a.voices[1])
		console.Printf("v2 %+v\n\n", a.voices[2])
		if freq == 0 {
			return nil
		}
		svol := vol * 8
		squareWavePeriod := uint32(a.SamplesPerSecond) / freq
		if squareWavePeriod == 0 {
			console.Printf("nil\n")
			return nil
		}
		halfSquareWavePeriod := squareWavePeriod / 2
		for i := uint32(0); i < bufSize; i++ {
			level := (i / halfSquareWavePeriod) % 2
			vol := uint8(128 + svol)
			if level == 1 {
				vol = uint8(128 - svol)
			}
			a.buf[i] = vol
		}
		return sdl.QueueAudio(1, a.buf)
	*/
	return nil
}

const (
	toneHz   = 440
	sampleHz = 48000
)

var phase float64

//export Callback
func Callback(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	v := voices[0]
	freq := uint32(v.Freq[0])&0x0f +
		(uint32(v.Freq[1])&0x0f)<<4 +
		(uint32(v.Freq[2])&0x0f)<<8 +
		(uint32(v.Freq[3])&0x0f)<<12 +
		(uint32(v.Freq[4])&0x0f)<<16
	dPhase := 2 * math.Pi * float64(freq) / float64(sampleHz)
	sample := C.Uint8(0)
	for i := 0; i < n; i += 2 {
		phase += dPhase
		sample = C.Uint8((math.Sin(phase) + 0.999999) * 128)
		if v.Vol == 0 {
			sample = 128
		}
		console.Printf("phase: %v, sample: %v\n", phase, sample)
		buf[i] = sample
		buf[i+1] = sample
	}
}
