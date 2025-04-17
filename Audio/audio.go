package Audio

import (
	"errors"
	"forgejo.max7.fun/m.alkhatib/GoForge/Math"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/gen2brain/malgo"
	"github.com/hajimehoshi/go-mp3"
	"github.com/youpy/go-wav"
)

type Engine struct {
	context *malgo.AllocatedContext
	mu      sync.Mutex
	devices map[*Sound]*malgo.Device
}

// Sound represents a loaded audio resource
type Sound struct {
	data       []byte
	channels   uint32
	sampleRate uint32
	format     malgo.FormatType
	volume     float32
	loop       bool
	pos        int
}

func NewEngine() (*Engine, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}
	return &Engine{
		context: ctx,
		devices: make(map[*Sound]*malgo.Device),
	}, nil
}

func (engine *Engine) Destroy() {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	for _, device := range engine.devices {
		device.Uninit()
	}

	err := engine.context.Uninit()
	if err != nil {
		return
	}
	engine.context.Free()
}

func (engine *Engine) LoadSound(filePath string) (*Sound, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []byte
	var channels, sampleRate uint32
	var format malgo.FormatType

	switch {
	case strings.HasSuffix(strings.ToLower(filePath), ".wav"):
		w := wav.NewReader(file)
		f, err := w.Format()
		if err != nil {
			return nil, err
		}

		channels = uint32(f.NumChannels)
		sampleRate = f.SampleRate
		format = malgo.FormatS16

		data, err = io.ReadAll(w)
		if err != nil {
			return nil, err
		}

	case strings.HasSuffix(strings.ToLower(filePath), ".mp3"):
		m, err := mp3.NewDecoder(file)
		if err != nil {
			return nil, err
		}

		channels = 2
		sampleRate = uint32(m.SampleRate())
		format = malgo.FormatS16

		data, err = io.ReadAll(m)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("unsupported audio format")
	}

	return &Sound{
		data:       data,
		channels:   channels,
		sampleRate: sampleRate,
		format:     format,
		volume:     1.0,
	}, nil
}

func (engine *Engine) Play(sound *Sound) error {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	if device, exists := engine.devices[sound]; exists {
		device.Stop()
		delete(engine.devices, sound)
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = sound.format
	deviceConfig.Playback.Channels = sound.channels
	deviceConfig.SampleRate = sound.sampleRate
	deviceConfig.Alsa.NoMMap = 1

	sound.pos = 0

	onSamples := func(pOutput, pInput []byte, framecount uint32) {
		engine.mu.Lock()
		defer engine.mu.Unlock()

		remaining := len(sound.data) - sound.pos
		if remaining <= 0 {
			if sound.loop {
				sound.pos = 0
				remaining = len(sound.data)
			} else {
				return
			}
		}

		toCopy := min(len(pOutput), remaining)
		copy(pOutput, sound.data[sound.pos:sound.pos+toCopy])
		sound.pos += toCopy

		if sound.volume != 1.0 {
			applyVolume(pOutput[:toCopy], sound.volume)
		}
	}

	device, err := malgo.InitDevice(engine.context.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: onSamples,
	})
	if err != nil {
		return err
	}

	err = device.Start()
	if err != nil {
		device.Uninit()
		return err
	}

	engine.devices[sound] = device
	return nil
}

func (engine *Engine) Stop(sound *Sound) {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	if device, exists := engine.devices[sound]; exists {
		device.Stop()
		device.Uninit()
		delete(engine.devices, sound)
	}
}

func (engine *Engine) SetVolume(sound *Sound, volume float32) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	sound.volume = Math.Clamp(volume, 0.0, 1.0)
}

func (engine *Engine) SetLooping(sound *Sound, loop bool) {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	sound.loop = loop
}

func (engine *Engine) IsPlaying(sound *Sound) bool {
	engine.mu.Lock()
	defer engine.mu.Unlock()
	_, exists := engine.devices[sound]
	return exists
}

func applyVolume(data []byte, volume float32) {
	for i := 0; i < len(data); i += 2 {
		sample := int16(float32(int16(data[i])|int16(data[i+1])<<8) * volume)
		data[i] = byte(sample)
		data[i+1] = byte(sample >> 8)
	}
}
