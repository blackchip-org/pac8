package audio

import (
	"math"
)

type Synth interface {
	Queue() error
}

type NullSynth struct{}

func (n NullSynth) Queue() error {
	return nil
}

type ToneGenerator struct {
	Freq     int
	Vol      float32
	Waveform []float32

	cycleFreq  int
	cycleVol   float32
	cycleWave  []float32
	sampleRate int
	phase      int
	period     int
}

func NewToneGenerator(sampleRate int) *ToneGenerator {
	return &ToneGenerator{sampleRate: sampleRate}
}

func (t *ToneGenerator) Fill(out []float32, n int) {
	for i := 0; i < n; i++ {
		if t.phase == t.period {
			t.phase = 0
			t.cycleFreq = t.Freq
			t.cycleVol = t.Vol
			t.cycleWave = t.Waveform
			if t.cycleFreq > 0 {
				t.period = t.sampleRate / t.cycleFreq
			}
		}
		if t.cycleFreq == 0 || t.cycleWave == nil {
			out[i] = 0
		} else {
			pos := float64(t.phase) / float64(t.period)
			wff := math.Round(pos * float64(len(t.cycleWave)-1))
			wfi := int(wff)
			out[i] = t.cycleWave[wfi] * t.cycleVol
		}
		t.phase++
		if t.phase > t.period {
			t.phase = 0
		}
	}
}
