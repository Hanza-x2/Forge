package Scene

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics/Viewports"
	"github.com/go-gl/mathgl/mgl32"
)

type Stage struct {
	Root      Actor
	Batch     *Graphics.Batch
	Viewport  Viewports.Viewport
	sortDirty bool
}

func NewStage(viewport Viewports.Viewport, batch *Graphics.Batch) *Stage {
	root := NewBaseActor()
	stage := &Stage{
		Root:     root,
		Batch:    batch,
		Viewport: viewport,
	}
	root.SetStage(stage)
	return stage
}

func (stage *Stage) AddActor(actor Actor) {
	stage.Root.AddChild(actor)
	stage.sortDirty = true
}

func (stage *Stage) RemoveActor(actor Actor) {
	stage.Root.RemoveChild(actor)
	//stage.sortDirty = true
}

func (stage *Stage) Clear() {
	stage.Root.RemoveAllChildren()
}

func (stage *Stage) Act(delta float32) {
	stage.Root.Act(delta)
}

func (stage *Stage) Draw() {
	stage.Viewport.Apply(false)
	stage.Batch.SetProjection(stage.Viewport.GetCamera().Matrix)
	stage.Batch.Begin()
	stage.Root.Draw(stage.Batch)
	stage.Batch.End()
}

func (stage *Stage) Resize(width, height int32) {
	stage.Viewport.Update(width, height, false)
}

func (stage *Stage) Dispose() {
	stage.Batch.Dispose()
}

func (stage *Stage) ScreenToStageCoordinates(screenX, screenY float32) (float32, float32) {
	camera := stage.Viewport.GetCamera()
	output := camera.Unproject(mgl32.Vec2{screenX, screenY}, camera.Width, camera.Height)
	return output.X(), output.Y()
}

func (stage *Stage) Hit(x, y float32) Actor {
	actors := stage.Root.GetChildren()
	for i := len(actors) - 1; i >= 0; i-- {
		actor := actors[i]
		if actor.Hit(x, y) {
			return actor
		}
	}
	return nil
}
