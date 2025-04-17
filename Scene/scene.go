package Scene

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics/Viewports"
	"github.com/go-gl/mathgl/mgl32"
)

type Scene struct {
	Root     *Node
	Batch    *Graphics.Batch
	Viewport Viewports.Viewport
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
}

func (scene *Scene) RemoveNode(node *Node) {
	scene.Root.RemoveChild(node)
}

func (scene *Scene) Clear() {
	scene.Root.RemoveAllChildren()
}

func hitChildren(parent *Node, x, y float32) *Node {
	nodes := parent.GetChildren()
	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]
		if node.Hit(x, y) {
			return node
		}
		if child := hitChildren(node, x-0.5, y-0.5); child != nil {
			return child
		}
	}
	return nil
}

func (scene *Scene) Hit(x, y float32) *Node {
	return hitChildren(scene.Root, x-0.5, y-0.5)
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

func drawNodeDebug(node *Node, batch *Graphics.Batch) {
	if !node.visible || node.width <= 0 || node.height <= 0 {
		return
	}

	// Get all 4 corners in local space
	corners := []mgl32.Vec2{
		{0, 0},
		{node.width, 0},
		{node.width, node.height},
		{0, node.height},
	}

	// Transform to world space
	transform := node.ComputeTransform()
	for i := range corners {
		corners[i] = transformCoordinate(corners[i].X()-node.originX, corners[i].Y()-node.originY, transform)
	}

	// Draw transformed rectangle
	for i := 0; i < 4; i++ {
		j := (i + 1) % 4
		batch.Line(
			corners[i].X(), corners[i].Y(),
			corners[j].X(), corners[j].Y(),
			2, Graphics.RED,
		)
	}

	// Draw origin point
	origin := transformCoordinate(-node.originX, -node.originY, transform)
	batch.FillRect(origin.X()-0.5, origin.Y()-0.5, 1, 1, Graphics.GREEN)

	// Draw children
	for _, child := range node.children {
		drawNodeDebug(child, batch)
	}
}

func (scene *Scene) DrawDebug() {
	scene.Viewport.Apply(false)
	scene.Batch.SetProjection(scene.Viewport.GetCamera().Matrix)
	scene.Batch.Begin()
	drawNodeDebug(scene.Root, scene.Batch)
	scene.Batch.End()
}

func (scene *Scene) Resize(width, height float32) {
	scene.Viewport.Update(width, height, false)
}

func (scene *Scene) Dispose() {
	scene.Batch.Dispose()
}
