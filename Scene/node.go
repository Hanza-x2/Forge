package Scene

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"math"
)

type Behavior interface {
	Draw(node *Node, batch *Graphics.Batch)
	Act(node *Node, delta float32)
}

type Node struct {
	Name     string
	Behavior Behavior
	X, Y     float32
	Width    float32
	Height   float32
	OriginX  float32
	OriginY  float32
	ScaleX   float32
	ScaleY   float32
	Rotation float32
	Visible  bool
	UserData interface{}
	parent   *Node
	children []*Node
	scene    *Scene
}

func NewNode() *Node {
	return &Node{
		ScaleX:   1,
		ScaleY:   1,
		Visible:  true,
		children: make([]*Node, 0),
	}
}

func (node *Node) SetPosition(x, y float32) {
	node.X, node.Y = x, y
}

func (node *Node) SetSize(width, height float32) {
	node.Width, node.Height = width, height
}

func (node *Node) SetOrigin(originX, originY float32) {
	node.OriginX, node.OriginY = originX, originY
}

func (node *Node) SetScale(scaleX, scaleY float32) {
	node.ScaleX, node.ScaleY = scaleX, scaleY
}

func (node *Node) GetParent() *Node {
	return node.parent
}

func (node *Node) SetParent(parent *Node) {
	if node.parent == parent {
		return
	}
	if node.parent != nil {
		node.parent.RemoveChild(node)
	}
	node.parent = parent
	if parent != nil {
		node.scene = parent.GetScene()
	} else {
		node.scene = nil
	}
}

func (node *Node) GetScene() *Scene {
	return node.scene
}

func (node *Node) SetScene(scene *Scene) {
	node.scene = scene
	for _, child := range node.children {
		child.SetScene(scene)
	}
}

func (node *Node) AddChild(child *Node) {
	if currentParent := child.GetParent(); currentParent != nil {
		currentParent.RemoveChild(child)
	}
	node.children = append(node.children, child)
	child.SetParent(node)
	child.SetScene(node.scene)
}

func (node *Node) RemoveChild(child *Node) bool {
	for i, c := range node.children {
		if c == child {
			node.children = append(node.children[:i], node.children[i+1:]...)
			child.SetParent(nil)
			child.SetScene(nil)
			return true
		}
	}
	return false
}

func (node *Node) GetChildren() []*Node {
	return node.children
}

func (node *Node) RemoveAllChildren() {
	for _, child := range node.children {
		child.SetParent(nil)
		child.SetScene(nil)
	}
	node.children = make([]*Node, 0)
}

func (node *Node) GetWorldTransform() (x, y, scaleX, scaleY, rotation float32) {
	x, y = node.X, node.Y
	scaleX, scaleY = node.ScaleX, node.ScaleY
	rotation = node.Rotation
	if parent := node.GetParent(); parent != nil {
		parentX, parentY, parentScaleX, parentScaleY, parentRotation := parent.GetWorldTransform()
		rad := float64(parentRotation * math.Pi / 180)
		cos := float32(math.Cos(rad))
		sin := float32(math.Sin(rad))
		rotatedX := x*cos - y*sin
		rotatedY := x*sin + y*cos
		x = parentX + rotatedX*parentScaleX
		y = parentY + rotatedY*parentScaleY
		scaleX *= parentScaleX
		scaleY *= parentScaleY
		rotation += parentRotation
	}
	return x, y, scaleX, scaleY, rotation
}

func (node *Node) GetWorldTransformEx() (x, y, originX, originY, width, height, scaleX, scaleY, rotation float32) {
	x, y, scaleX, scaleY, rotation = node.GetWorldTransform()
	return x, y, node.OriginX, node.OriginY, node.Width, node.Height, scaleX, scaleY, rotation
}

func (node *Node) Act(delta float32) {
	if node.Behavior != nil {
		node.Behavior.Act(node, delta)
	}
	for _, child := range node.children {
		child.Act(delta)
	}
}

func (node *Node) Draw(batch *Graphics.Batch) {
	if !node.Visible {
		return
	}
	if node.Behavior != nil {
		node.Behavior.Draw(node, batch)
	}
	for _, child := range node.children {
		child.Draw(batch)
	}
}

func (node *Node) LocalToParentCoordinates(localX, localY float32) (float32, float32) {
	rotation := -node.Rotation
	if rotation == 0 {
		if node.ScaleX == 1 && node.ScaleY == 1 {
			localX += node.X
			localY += node.Y
		} else {
			localX = (localX-node.OriginX)*node.ScaleX + node.OriginX + node.X
			localY = (localY-node.OriginY)*node.ScaleY + node.OriginY + node.Y
		}
	} else {
		rad := float64(rotation * math.Pi / 180)
		cos := float32(math.Cos(rad))
		sin := float32(math.Sin(rad))
		toX := (localX - node.OriginX) * node.ScaleX
		toY := (localY - node.OriginY) * node.ScaleY
		localX = (toX*cos + toY*sin) + node.OriginX + node.X
		localY = (toX*-sin + toY*cos) + node.OriginY + node.Y
	}
	return localX, localY
}

func (node *Node) ParentToLocalCoordinates(parentX, parentY float32) (float32, float32) {
	if node.Rotation == 0 {
		if node.ScaleX == 1 && node.ScaleY == 1 {
			parentX -= node.X
			parentY -= node.Y
		} else {
			parentX = (parentX-node.X-node.OriginX)/node.ScaleX + node.OriginX
			parentY = (parentY-node.Y-node.OriginY)/node.ScaleY + node.OriginY
		}
	} else {
		rad := float64(node.Rotation * math.Pi / 180)
		cos := float32(math.Cos(rad))
		sin := float32(math.Sin(rad))
		toX := parentX - node.X - node.OriginX
		toY := parentY - node.Y - node.OriginY
		parentX = (toX*cos+toY*sin)/node.ScaleX + node.OriginX
		parentY = (toX*-sin+toY*cos)/node.ScaleY + node.OriginY
	}
	return parentX, parentY
}

func (node *Node) Hit(x, y float32) bool {
	if !node.Visible {
		return false
	}
	return x >= 0 && y >= 0 && x < node.Width && y < node.Height
}
