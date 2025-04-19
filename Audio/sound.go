package Audio

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

type Sound struct {
	buffer  *beep.Buffer
	format  beep.Format
	ctrl    *beep.Ctrl
	volume  *effects.Volume
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

	if !speakerInitialized {
		err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/20))
		if err != nil {
			_ = streamer.Close()
			return nil, err
		}
		speakerInitialized = true
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	_ = streamer.Close()

	vol := &effects.Volume{
		Streamer: nil,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	return &Sound{
		buffer:  buffer,
		format:  format,
		volume:  vol,
		playing: false,
	}, nil
}

func (sound *Sound) Play(volume float32) {
	sound.Stop()

	streamer := sound.buffer.Streamer(0, sound.buffer.Len())
	sound.volume.Streamer = streamer
	sound.SetVolume(volume)
	sound.ctrl = &beep.Ctrl{Streamer: sound.volume}
	sound.playing = true

	speaker.Play(beep.Seq(sound.ctrl, beep.Callback(func() {
		sound.playing = false
	})))
}

func (sound *Sound) Stop() {
	if sound.ctrl != nil {
		speaker.Lock()
		sound.ctrl.Paused = true
		speaker.Unlock()
	}
	sound.playing = false
}

func (sound *Sound) SetVolume(volume float32) {
	sound.volume.Volume = float64(-5 + (5 * volume))
}

func (sound *Sound) IsPlaying() bool {
	return sound.playing
}
