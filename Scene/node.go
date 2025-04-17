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
	name     string
	behavior Behavior
	x, y     float32
	width    float32
	height   float32
	originX  float32
	originY  float32
	scaleX   float32
	scaleY   float32
	rotation float32
	visible  bool
	parent   *Node
	children []*Node
	userData interface{}
	scene    *Scene
	dirty    bool
}

func NewNode() *Node {
	return &Node{
		scaleX:   1,
		scaleY:   1,
		visible:  true,
		children: make([]*Node, 0),
		dirty:    true,
	}
}

func (node *Node) GetName() string {
	return node.name
}

func (node *Node) SetName(name string) {
	node.name = name
}

func (node *Node) GetBehavior() Behavior {
	return node.behavior
}

func (node *Node) SetBehavior(behavior Behavior) {
	node.behavior = behavior
}

func (node *Node) GetX() float32 {
	return node.x
}

func (node *Node) GetY() float32 {
	return node.y
}

func (node *Node) SetPosition(x, y float32) {
	node.x, node.y = x, y
	node.dirty = true
}

func (node *Node) GetWidth() float32 {
	return node.width
}

func (node *Node) GetHeight() float32 {
	return node.height
}

func (node *Node) SetSize(width, height float32) {
	node.width, node.height = width, height
	node.dirty = true
}

func (node *Node) GetOriginX() float32 {
	return node.originX
}

func (node *Node) GetOriginY() float32 {
	return node.originY
}

func (node *Node) SetOrigin(originX, originY float32) {
	node.originX, node.originY = originX, originY
	node.dirty = true
}

func (node *Node) GetScaleX() float32 {
	return node.scaleX
}

func (node *Node) GetScaleY() float32 {
	return node.scaleY
}

func (node *Node) SetScale(scaleX, scaleY float32) {
	node.scaleX, node.scaleY = scaleX, scaleY
	node.dirty = true
}

func (node *Node) GetRotation() float32 {
	return node.rotation
}

func (node *Node) SetRotation(degrees float32) {
	node.rotation = degrees
	node.dirty = true
}

func (node *Node) IsVisible() bool {
	return node.visible
}

func (node *Node) SetVisible(visible bool) {
	node.visible = visible
}

func (node *Node) GetUserData() interface{} {
	return node.userData
}

func (node *Node) SetUserData(data interface{}) {
	node.userData = data
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
	x, y = node.x, node.y
	scaleX, scaleY = node.scaleX, node.scaleY
	rotation = node.rotation
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
	return x, y, node.originX, node.originY, node.width, node.height, scaleX, scaleY, rotation
}

func (node *Node) Act(delta float32) {
	if node.behavior != nil {
		node.behavior.Act(node, delta)
	}
	for _, child := range node.children {
		child.Act(delta)
	}
}

func (node *Node) Draw(batch *Graphics.Batch) {
	if !node.visible {
		return
	}
	if node.behavior != nil {
		node.behavior.Draw(node, batch)
	}
	for _, child := range node.children {
		child.Draw(batch)
	}
}

func (node *Node) LocalToParentCoordinates(localX, localY float32) (float32, float32) {
	rotation := -node.rotation
	if rotation == 0 {
		if node.scaleX == 1 && node.scaleY == 1 {
			localX += node.x
			localY += node.y
		} else {
			localX = (localX-node.originX)*node.scaleX + node.originX + node.x
			localY = (localY-node.originY)*node.scaleY + node.originY + node.y
		}
	} else {
		rad := float64(rotation * math.Pi / 180)
		cos := float32(math.Cos(rad))
		sin := float32(math.Sin(rad))
		toX := (localX - node.originX) * node.scaleX
		toY := (localY - node.originY) * node.scaleY
		localX = (toX*cos + toY*sin) + node.originX + node.x
		localY = (toX*-sin + toY*cos) + node.originY + node.y
	}
	return localX, localY
}

func (node *Node) ParentToLocalCoordinates(parentX, parentY float32) (float32, float32) {
	if node.rotation == 0 {
		if node.scaleX == 1 && node.scaleY == 1 {
			parentX -= node.x
			parentY -= node.y
		} else {
			parentX = (parentX-node.x-node.originX)/node.scaleX + node.originX
			parentY = (parentY-node.y-node.originY)/node.scaleY + node.originY
		}
	} else {
		rad := float64(node.rotation * math.Pi / 180)
		cos := float32(math.Cos(rad))
		sin := float32(math.Sin(rad))
		toX := parentX - node.x - node.originX
		toY := parentY - node.y - node.originY
		parentX = (toX*cos+toY*sin)/node.scaleX + node.originX
		parentY = (toX*-sin+toY*cos)/node.scaleY + node.originY
	}
	return parentX, parentY
}

func (node *Node) Hit(x, y float32) bool {
	if !node.visible {
		return false
	}
	return x >= 0 && y >= 0 && x < node.width && y < node.height
}
