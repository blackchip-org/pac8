package audio

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	SampleRate = 22050
	Channels   = 2
	QueueRate  = 60
	Format     = sdl.AUDIO_U16LSB
	Buffer     = 5
)

type Audio interface {
	Queue() error
}

type NullAudio struct{}

func (n NullAudio) Queue() error {
	return nil
}

type Voice struct {
	Freq     int
	Vol      float64
	Waveform []float64

	cycleFreq  int
	cycleVol   float64
	cycleWave  []float64
	sampleRate int
	phase      int
	period     int
}

func NewVoice(sampleRate int32) *Voice {
	return &Voice{sampleRate: int(sampleRate)}
}

func (v *Voice) Fill(out []float64, n int) {
	for i := 0; i < n; i++ {
		if v.phase == 0 {
			v.phase = 0
			v.cycleFreq = v.Freq
			v.cycleVol = v.Vol
			v.cycleWave = v.Waveform
			if v.cycleFreq > 0 {
				v.period = (v.sampleRate / v.cycleFreq) - 1
			}
		}
		if v.cycleFreq == 0 || v.cycleWave == nil {
			out[i] = 0
		} else {
			pct := float64(v.phase) / float64(v.period)
			pos := math.Round(pct * float64(len(v.cycleWave)-1))
			out[i] = float64(v.cycleWave[int(pos)]) * v.cycleVol
		}
		v.phase++
		if v.phase > v.period {
			v.phase = 0
		}
	}
}

type Synth struct {
	Spec sdl.AudioSpec
	V    []*Voice

	samples [][]float64
	mixed   []float64
	data    []byte
}

func NewSynth(spec sdl.AudioSpec, voiceN int) (*Synth, error) {
	s := &Synth{}
	s.Spec = spec
	/*
		if s.Spec.Format != sdl.AUDIO_U16LSB {
			return nil, fmt.Errorf("unsupported format: %v", s.Spec.Format)
		}
	*/
	if s.Spec.Channels != Channels {
		return nil, fmt.Errorf("unsupported channel count: %v", s.Spec.Channels)
	}
	s.V = make([]*Voice, voiceN)
	samplesLen := s.Spec.Samples * Buffer
	s.samples = make([][]float64, voiceN, voiceN)
	for v := 0; v < voiceN; v++ {
		s.V[v] = NewVoice(s.Spec.Freq)
		s.samples[v] = make([]float64, samplesLen, samplesLen)
	}
	s.mixed = make([]float64, samplesLen)
	dataLen := 4 * int(samplesLen)
	s.data = make([]byte, dataLen, dataLen)
	return s, nil
}

func (s *Synth) Queue() error {
	q := sdl.GetQueuedAudioSize(1) / 4
	n := int(s.Spec.Samples*Buffer) - int(q)
	if n <= 0 {
		return nil
	}
	for i := 0; i < len(s.V); i++ {
		s.V[i].Fill(s.samples[i], n)
	}
	for i := 0; i < n; i++ {
		sample := float64(0)
		for j := 0; j < len(s.V); j++ {
			sample += s.samples[j][i]
		}
		s.mixed[i] = sample / float64(len(s.V))
		if s.mixed[i] > 1 {
			fmt.Println("over")
		}
		if s.mixed[i] < -1 {
			fmt.Println("under")
		}
	}
	for i, d := 0, 0; i < n; i, d = i+1, d+4 {
		sample := convert(s.mixed[i])
		fmt.Printf("%04x\n", sample)
		s.data[d+0] = byte(sample & 0xff)
		s.data[d+1] = byte((sample >> 8) & 0xff)
		s.data[d+2] = byte(sample & 0xff)
		s.data[d+3] = byte((sample >> 8) & 0xff)
	}
	return sdl.QueueAudio(1, s.data[0:n*4])
}

func convert(f float64) uint16 {
	return uint16(((f + 1) / 2) * (1<<16 - 1))
}
