package Scene

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics/Viewports"
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

func (scene *Scene) Hit(x, y float32) *Node {
	x -= 0.5
	y -= 0.5
	nodes := scene.Root.GetChildren()
	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]
		if node.Hit(x, y) {
			return node
		}
	}
	return nil
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

func (scene *Scene) DrawDebug() {
	scene.Viewport.Apply(false)
	scene.Batch.SetProjection(scene.Viewport.GetCamera().Matrix)
	scene.Batch.Begin()
	scene.Batch.SetColor(Graphics.YELLOW)
	nodes := scene.Root.GetChildren()
	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]
		x, y, originX, originY, width, height, scaleX, scaleY, rotation := node.GetWorldTransformEx()
		scene.Batch.LineRectEx(x, y, originX, originY, width, height, scaleX, scaleY, rotation, 1)
	}
	scene.Batch.SetColor(Graphics.WHITE)
	scene.Batch.End()
}

func (scene *Scene) Resize(width, height float32) {
	scene.Viewport.Update(width, height, false)
}

func (scene *Scene) Dispose() {
	scene.Batch.Dispose()
}
