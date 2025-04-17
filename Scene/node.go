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
	name      string
	behavior  Behavior
	x, y      float32
	width     float32
	height    float32
	originX   float32
	originY   float32
	scaleX    float32
	scaleY    float32
	rotation  float32
	visible   bool
	parent    *Node
	children  []*Node
	userData  interface{}
	scene     *Scene
	transform mgl32.Mat3
	dirty     bool
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

func transformCoordinate(vecX, vecY float32, mat mgl32.Mat3) mgl32.Vec2 {
	x := vecX*mat[0] + vecY*mat[3] + mat[6]
	y := vecX*mat[1] + vecY*mat[4] + mat[7]
	return mgl32.Vec2{x, y}
}

func (node *Node) GetWorldTransform() (x, y, scaleX, scaleY, rotation float32) {
	// Start with local transform
	x, y = node.x, node.y
	scaleX, scaleY = node.scaleX, node.scaleY
	rotation = node.rotation

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
	return x, y, node.originX, node.originY, node.width, node.height, scaleX, scaleY, rotation
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

	transform = transform.Mul3(mgl32.Translate2D(node.x, node.y))
	if node.rotation != 0 {
		transform = transform.Mul3(mgl32.HomogRotate2D(mgl32.DegToRad(node.rotation)))
	}

	if node.scaleX != 1 || node.scaleY != 1 {
		transform = transform.Mul3(mgl32.Scale2D(node.scaleX, node.scaleY))
	}

	if node.originX != 0 || node.originY != 0 {
		transform = transform.Mul3(mgl32.Translate2D(-node.originX, -node.originY))
	}

	if node.parent == nil {
		node.transform = transform
		node.dirty = false
	}

	return transform
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

func (node *Node) Act(delta float32) {
	if node.behavior != nil {
		node.behavior.Act(node, delta)
	}
	for _, child := range node.children {
		child.Act(delta)
	}
}

func (node *Node) Hit(x, y float32) bool {
	if !node.visible {
		return false
	}
	localX, localY := node.SceneToLocalCoordinates(x, y)
	return localX >= 0 && localY >= 0 && localX < node.width && localY < node.height
}
