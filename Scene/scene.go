package Scene

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics/Viewports"
	"github.com/go-gl/mathgl/mgl32"
)

type Scene struct {
	Root      *Node
	Batch     *Graphics.Batch
	Viewport  Viewports.Viewport
	sortDirty bool
}

func NewScene(viewport Viewports.Viewport, batch *Graphics.Batch) *Scene {
	root := NewNode()
	scene := &Scene{
		Root:     root,
		Batch:    batch,
		Viewport: viewport,
	}
	root.SetScene(scene)
	return scene
}

func (scene *Scene) AddNode(node *Node) {
	scene.Root.AddChild(node)
	scene.sortDirty = true
}

func (scene *Scene) RemoveNode(node *Node) {
	scene.Root.RemoveChild(node)
}

func (scene *Scene) Clear() {
	scene.Root.RemoveAllChildren()
}

func (scene *Scene) Act(delta float32) {
	scene.Root.Act(delta)
}

func (scene *Scene) Draw() {
	scene.Viewport.Apply(false)
	scene.Batch.SetProjection(scene.Viewport.GetCamera().Matrix)
	scene.Batch.Begin()
	scene.Root.Draw(scene.Batch)
	scene.Batch.End()
}

func (scene *Scene) Resize(width, height int32) {
	scene.Viewport.Update(width, height, false)
}

func (scene *Scene) Dispose() {
	scene.Batch.Dispose()
}

func (scene *Scene) ScreenToSceneCoordinates(screenX, screenY float32) (float32, float32) {
	camera := scene.Viewport.GetCamera()
	output := camera.Unproject(mgl32.Vec2{screenX, screenY}, camera.Width, camera.Height)
	return output.X(), output.Y()
}

func (scene *Scene) Hit(x, y float32) *Node {
	nodes := scene.Root.GetChildren()
	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]
		if node.Hit(x, y) {
			return node
		}
	}
	return nil
}
