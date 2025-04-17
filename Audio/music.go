package Audio

import (
	"errors"
	"math"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type Music struct {
	filePath string
	streamer beep.StreamSeekCloser
	format   beep.Format
	ctrl     *beep.Ctrl
	gain     *effects.Gain
	loop     bool
	playing  bool
}

func NewMusic(filePath string) (*Music, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	var streamer beep.StreamSeekCloser
	var format beep.Format

	switch {
	case strings.HasSuffix(strings.ToLower(filePath), ".wav"):
		streamer, format, err = wav.Decode(file)
	case strings.HasSuffix(strings.ToLower(filePath), ".mp3"):
		streamer, format, err = mp3.Decode(file)
	default:
		file.Close()
		return nil, errors.New("unsupported audio format")
	}
	if err != nil {
		file.Close()
		return nil, err
	}

	if speakerInitialized == false {
		err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			streamer.Close()
			return nil, err
		}
		speakerInitialized = true
	}

	return &Music{
		filePath: filePath,
		streamer: streamer,
		format:   format,
		gain:     &effects.Gain{Gain: 0},
		playing:  false,
	}, nil
}

func (music *Music) Play() {
	if music.playing {
		return
	}
	err := music.streamer.Seek(0)
	if err != nil {
		return
	}

	var baseStreamer beep.Streamer
	if music.loop {
		baseStreamer = beep.Loop(-1, music.streamer)
	} else {
		baseStreamer = music.streamer
	}

	music.gain.Streamer = baseStreamer
	music.ctrl = &beep.Ctrl{Streamer: music.gain}

	go func() {
		speaker.Play(beep.Seq(music.ctrl, beep.Callback(func() {
			music.playing = false
		})))
	}()

	music.playing = true
}

func (music *Music) Stop() {
	if music.ctrl != nil {
		speaker.Lock()
		music.ctrl.Paused = true
		speaker.Unlock()
	}
	music.playing = false
}

func (music *Music) SetVolume(volume float32) {
	if volume <= 0 {
		music.gain.Gain = -100
	} else {
		music.gain.Gain = 20 * float64(math.Log10(float64(volume)))
	}
}

func (music *Music) SetLooping(loop bool) {
	music.loop = loop
}

func (music *Music) IsPlaying() bool {
	return music.playing
}

func (music *Music) Close() {
	if music.streamer != nil {
		music.streamer.Close()
	}
}
