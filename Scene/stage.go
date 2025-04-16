package Scene

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics/Viewports"
	"github.com/go-gl/mathgl/mgl32"
	"sort"
)

type Stage struct {
	Root      *Actor
	Actors    []*Actor
	Batch     *Graphics.Batch
	Viewport  Viewports.Viewport
	sortDirty bool
}

func NewStage(viewport Viewports.Viewport, batch *Graphics.Batch) *Stage {
	stage := &Stage{
		Root:     NewActor(),
		Batch:    batch,
		Viewport: viewport,
		Actors:   make([]*Actor, 0),
	}
	stage.Root.stage = stage
	return stage
}

func (stage *Stage) AddActor(actor *Actor) {
	stage.Root.AddChild(actor)
	stage.Actors = append(stage.Actors, actor)
	stage.sortDirty = true
}

func (stage *Stage) RemoveActor(actor *Actor) bool {
	if stage.Root.RemoveChild(actor) {
		for i, a := range stage.Actors {
			if a == actor {
				stage.Actors = append(stage.Actors[:i], stage.Actors[i+1:]...)
				break
			}
		}
		return true
	}
	return false
}

func (stage *Stage) Clear() {
	stage.Root.Children = make([]*Actor, 0)
	stage.Actors = make([]*Actor, 0)
}

func (stage *Stage) Act(delta float32) {
	stage.Root.Act(delta)
}

func (stage *Stage) Draw() {
	stage.Viewport.Apply(false)
	stage.Batch.SetProjection(stage.Viewport.GetCamera().Matrix)

	stage.Batch.Begin()
	if stage.sortDirty {
		sort.SliceStable(stage.Actors, func(i, j int) bool {
			return stage.Actors[i].ZIndex < stage.Actors[j].ZIndex
		})
		stage.sortDirty = false
	}

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

func (stage *Stage) Hit(x, y float32) *Actor {
	for i := len(stage.Actors) - 1; i >= 0; i-- {
		actor := stage.Actors[i]
		if actor.Hit(x, y) {
			return actor
		}
	}
	return nil
}
