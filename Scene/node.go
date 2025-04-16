package Scene

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Behavior interface {
	Draw(node *Node, batch *Graphics.Batch)
	Act(node *Node, delta float32)
}

type Node struct {
	Name      string
	Behavior  Behavior
	X, Y      float32
	Width     float32
	Height    float32
	OriginX   float32
	OriginY   float32
	ScaleX    float32
	ScaleY    float32
	Rotation  float32
	Visible   bool
	UserData  interface{}
	children  []*Node
	parent    *Node
	scene     *Scene
	transform mgl32.Mat3
	dirty     bool
}

func NewNode() *Node {
	return &Node{
		ScaleX:   1,
		ScaleY:   1,
		Visible:  true,
		children: make([]*Node, 0),
		dirty:    true,
	}
}

func (node *Node) SetPosition(x, y float32) {
	node.X, node.Y = x, y
	node.dirty = true
}

func (node *Node) SetSize(width, height float32) {
	node.Width, node.Height = width, height
	node.dirty = true
}

func (node *Node) SetOrigin(originX, originY float32) {
	node.OriginX, node.OriginY = originX, originY
	node.dirty = true
}

func (node *Node) SetScale(scaleX, scaleY float32) {
	node.ScaleX, node.ScaleY = scaleX, scaleY
	node.dirty = true
}

func (node *Node) SetRotation(degrees float32) {
	node.Rotation = degrees
	node.dirty = true
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

func transformCoordinate(vecX, vecY float32, mat mgl32.Mat3) mgl32.Vec2 {
	x := vecX*mat[0] + vecY*mat[3] + mat[6]
	y := vecX*mat[1] + vecY*mat[4] + mat[7]
	return mgl32.Vec2{x, y}
}

func (node *Node) GetWorldTransform() (x, y, scaleX, scaleY, rotation float32) {
	// Start with local transform
	x, y = node.X, node.Y
	scaleX, scaleY = node.ScaleX, node.ScaleY
	rotation = node.Rotation

	// Apply parent transforms recursively
	if parent := node.GetParent(); parent != nil {
		parentX, parentY, parentScaleX, parentScaleY, parentRot := parent.GetWorldTransform()

		rad := float64(parentRot * math.Pi / 180)
		cos := float32(math.Cos(rad))
		sin := float32(math.Sin(rad))

		// Calculate combined position
		rotatedX := x*cos - y*sin
		rotatedY := x*sin + y*cos
		x = parentX + rotatedX*parentScaleX
		y = parentY + rotatedY*parentScaleY

		// Combine scale
		scaleX *= parentScaleX
		scaleY *= parentScaleY

		// Combine rotation
		rotation += parentRot
	}

	return x, y, scaleX, scaleY, rotation
}

func (node *Node) GetWorldTransformEx() (x, y, originX, originY, width, height, scaleX, scaleY, rotation float32) {
	x, y, scaleX, scaleY, rotation = node.GetWorldTransform()
	return x, y, node.OriginX, node.OriginY, node.Width, node.Height, scaleX, scaleY, rotation
}

func (node *Node) LocalToSceneCoordinates(localX, localY float32) (float32, float32) {
	transform := node.ComputeTransform()
	vec := transformCoordinate(localX, localY, transform)
	return vec.X(), vec.Y()
}

func (node *Node) SceneToLocalCoordinates(sceneX, sceneY float32) (float32, float32) {
	transform := node.ComputeTransform().Inv()
	vec := transformCoordinate(sceneX, sceneY, transform)
	return vec.X(), vec.Y()
}

func (node *Node) ComputeTransform() mgl32.Mat3 {
	if !node.dirty && node.parent == nil {
		return node.transform
	}

	transform := mgl32.Ident3()
	if node.parent != nil {
		transform = node.parent.ComputeTransform()
	}

	transform = transform.Mul3(mgl32.Translate2D(node.X, node.Y))
	if node.Rotation != 0 {
		transform = transform.Mul3(mgl32.HomogRotate2D(mgl32.DegToRad(node.Rotation)))
	}

	if node.ScaleX != 1 || node.ScaleY != 1 {
		transform = transform.Mul3(mgl32.Scale2D(node.ScaleX, node.ScaleY))
	}

	if node.OriginX != 0 || node.OriginY != 0 {
		transform = transform.Mul3(mgl32.Translate2D(-node.OriginX, -node.OriginY))
	}

	if node.parent == nil {
		node.transform = transform
		node.dirty = false
	}

	return transform
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

func (node *Node) Act(delta float32) {
	if node.Behavior != nil {
		node.Behavior.Act(node, delta)
	}
	for _, child := range node.children {
		child.Act(delta)
	}
}

func (node *Node) Hit(x, y float32) bool {
	if !node.Visible {
		return false
	}
	localX, localY := node.SceneToLocalCoordinates(x, y)
	return localX >= 0 && localY >= 0 && localX < node.Width && localY < node.Height
}
