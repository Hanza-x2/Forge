package Audio

import (
	"errors"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type Sound struct {
	buffer  *beep.Buffer
	format  beep.Format
	ctrl    *beep.Ctrl
	gain    *effects.Gain
	mu      sync.Mutex
	playing bool
}

func NewSound(filePath string) (*Sound, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var streamer beep.StreamSeekCloser
	var format beep.Format

	switch {
	case strings.HasSuffix(strings.ToLower(filePath), ".wav"):
		streamer, format, err = wav.Decode(file)
	case strings.HasSuffix(strings.ToLower(filePath), ".mp3"):
		streamer, format, err = mp3.Decode(file)
	default:
		return nil, errors.New("unsupported audio format")
	}
	if err != nil {
		return nil, err
	}

	if speakerInitialized == false {
		err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			return nil, err
		}
		speakerInitialized = true
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	return &Sound{
		buffer:  buffer,
		format:  format,
		gain:    &effects.Gain{Gain: 0},
		playing: false,
	}, nil
}

func (sound *Sound) Play(volume float32) {
	sound.mu.Lock()
	sound.Stop()

	streamer := sound.buffer.Streamer(0, sound.buffer.Len())
	sound.gain.Streamer = streamer
	sound.SetVolume(volume)
	sound.ctrl = &beep.Ctrl{Streamer: sound.gain}
	sound.playing = true
	sound.mu.Unlock()

	speaker.Play(beep.Seq(sound.ctrl, beep.Callback(func() {
		sound.mu.Lock()
		sound.playing = false
		sound.mu.Unlock()
	})))
}

func (sound *Sound) Stop() {
	sound.mu.Lock()
	defer sound.mu.Unlock()

	if sound.ctrl != nil {
		speaker.Lock()
		sound.ctrl.Paused = true
		speaker.Unlock()
	}
	sound.playing = false
}

func (sound *Sound) SetVolume(volume float32) {
	sound.mu.Lock()
	defer sound.mu.Unlock()
	if volume <= 0 {
		sound.gain.Gain = -100
	} else {
		sound.gain.Gain = 20 * math.Log10(float64(volume))
	}
}

func (sound *Sound) IsPlaying() bool {
	sound.mu.Lock()
	defer sound.mu.Unlock()
	return sound.playing
}
