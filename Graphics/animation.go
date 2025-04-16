package Graphics

type Animation struct {
	Frames        []*TextureRegion
	FrameDuration float32
	Loop          bool
}

func NewAnimation(frames []*TextureRegion, frameDuration float32, loop bool) *Animation {
	return &Animation{
		Frames:        frames,
		FrameDuration: frameDuration,
		Loop:          loop,
	}
}

func (a *Animation) GetFrame(time float32) *TextureRegion {
	if len(a.Frames) == 1 {
		return a.Frames[0]
	}

	frameIndex := int(time / a.FrameDuration)
	if a.Loop {
		return a.Frames[frameIndex%len(a.Frames)]
	}

	if frameIndex >= len(a.Frames)-1 {
		return a.Frames[len(a.Frames)-1]
	}
	return a.Frames[frameIndex]
}

func (a *Animation) IsFinished(time float32) bool {
	if a.Loop {
		return false
	}
	return int(time/a.FrameDuration) >= len(a.Frames)-1
}
