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

type Music struct {
	filePath string
	streamer beep.StreamSeekCloser
	format   beep.Format
	ctrl     *beep.Ctrl
	volume   *effects.Volume
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
		_ = file.Close()
		return nil, errors.New("unsupported audio format")
	}
	if err != nil {
		_ = file.Close()
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

	vol := &effects.Volume{
		Streamer: nil,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	return &Music{
		filePath: filePath,
		streamer: streamer,
		format:   format,
		volume:   vol,
		playing:  false,
	}, nil
}

func (music *Music) Play() {
	if music.playing {
		return
	}

	if err := music.streamer.Seek(0); err != nil {
		return
	}

	var baseStreamer beep.Streamer
	if music.loop {
		baseStreamer, _ = beep.Loop2(music.streamer)
	} else {
		baseStreamer = music.streamer
	}

	music.volume.Streamer = baseStreamer
	music.ctrl = &beep.Ctrl{Streamer: music.volume}

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
	music.volume.Volume = float64(-5 + (5 * volume))
}

func (music *Music) SetLooping(loop bool) {
	music.loop = loop
}

func (music *Music) IsPlaying() bool {
	return music.playing
}

func (music *Music) Close() {
	if music.streamer != nil {
		_ = music.streamer.Close()
	}
}
