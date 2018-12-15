package audio

type Synth interface {
	Queue() error
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
			pos := float32(t.phase) / float32(t.period)
			wfi := int(pos * float32(len(t.cycleWave)))
			out[i] = t.cycleWave[wfi] * t.cycleVol
		}
		t.phase++
		if t.phase > t.period {
			t.phase = 0
		}
	}
}
