package Scene

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Behavior interface {
	Draw(actor *Actor, batch *Graphics.Batch)
	Act(actor *Actor, delta float32)
}

type Actor struct {
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
	ZIndex    int
	Visible   bool
	Parent    *Actor
	Children  []*Actor
	UserData  interface{}
	stage     *Stage
	transform mgl32.Mat3
	dirty     bool
}

func NewActor() *Actor {
	return &Actor{
		ScaleX:   1,
		ScaleY:   1,
		Visible:  true,
		Children: make([]*Actor, 0),
		dirty:    true,
	}
}

func (actor *Actor) GetName() string {
	return actor.Name
}

func (actor *Actor) GetBehavior() Behavior {
	return actor.Behavior
}

func (actor *Actor) SetBehavior(behavior Behavior) {
	actor.Behavior = behavior
}

func (actor *Actor) SetName(name string) {
	actor.Name = name
}

func (actor *Actor) GetX() float32 {
	return actor.X
}

func (actor *Actor) GetY() float32 {
	return actor.Y
}

func (actor *Actor) SetPosition(x, y float32) {
	actor.X, actor.Y = x, y
	actor.dirty = true
}

func (actor *Actor) GetWidth() float32 {
	return actor.Width
}

func (actor *Actor) GetHeight() float32 {
	return actor.Height
}

func (actor *Actor) SetSize(width, height float32) {
	actor.Width, actor.Height = width, height
	actor.dirty = true
}

func (actor *Actor) GetOriginX() float32 {
	return actor.OriginX
}

func (actor *Actor) GetOriginY() float32 {
	return actor.OriginY
}

func (actor *Actor) SetOrigin(originX, originY float32) {
	actor.OriginX, actor.OriginY = originX, originY
	actor.dirty = true
}

func (actor *Actor) GetScaleX() float32 {
	return actor.ScaleX
}

func (actor *Actor) GetScaleY() float32 {
	return actor.ScaleY
}

func (actor *Actor) SetScale(scaleX, scaleY float32) {
	actor.ScaleX, actor.ScaleY = scaleX, scaleY
	actor.dirty = true
}

func (actor *Actor) GetRotation() float32 {
	return actor.Rotation
}

func (actor *Actor) SetRotation(degrees float32) {
	actor.Rotation = degrees
	actor.dirty = true
}

func (actor *Actor) GetZIndex() int {
	return actor.ZIndex
}

func (actor *Actor) SetZIndex(zIndex int) {
	actor.ZIndex = zIndex
	actor.dirty = true
}

func (actor *Actor) IsVisible() bool {
	return actor.Visible
}

func (actor *Actor) SetVisible(visible bool) {
	actor.Visible = visible
}

func (actor *Actor) GetUserData() interface{} {
	return actor.UserData
}

func (actor *Actor) SetUserData(data interface{}) {
	actor.UserData = data
}

func (actor *Actor) GetParent() *Actor {
	return actor.Parent
}

func (actor *Actor) SetParent(parent *Actor) {
	if actor.Parent == parent {
		return
	}
	if actor.Parent != nil {
		actor.Parent.RemoveChild(actor)
	}
	actor.Parent = parent
	if parent != nil {
		actor.stage = parent.GetStage()
	} else {
		actor.stage = nil
	}
}

func (actor *Actor) GetStage() *Stage {
	return actor.stage
}

func (actor *Actor) SetStage(stage *Stage) {
	actor.stage = stage
	for _, child := range actor.Children {
		child.SetStage(stage)
	}
}

func (actor *Actor) AddChild(child *Actor) {
	if currentParent := child.GetParent(); currentParent != nil {
		currentParent.RemoveChild(child)
	}
	actor.Children = append(actor.Children, child)
	child.SetParent(actor)
	child.SetStage(actor.stage)
}

func (actor *Actor) RemoveChild(child *Actor) bool {
	for i, c := range actor.Children {
		if c == child {
			actor.Children = append(actor.Children[:i], actor.Children[i+1:]...)
			child.SetParent(nil)
			child.SetStage(nil)
			return true
		}
	}
	return false
}

func (actor *Actor) GetChildren() []*Actor {
	return actor.Children
}

func (actor *Actor) RemoveAllChildren() {
	for _, child := range actor.Children {
		child.SetParent(nil)
		child.SetStage(nil)
	}
	actor.Children = make([]*Actor, 0)
}

func transformCoordinate(vecX, vecY float32, mat mgl32.Mat3) mgl32.Vec2 {
	x := vecX*mat[0] + vecY*mat[3] + mat[6]
	y := vecX*mat[1] + vecY*mat[4] + mat[7]
	return mgl32.Vec2{x, y}
}

func (actor *Actor) GetWorldTransform() (x, y, scaleX, scaleY, rotation float32) {
	// Start with local transform
	x, y = actor.X, actor.Y
	scaleX, scaleY = actor.ScaleX, actor.ScaleY
	rotation = actor.Rotation

	// Apply parent transforms recursively
	if parent := actor.GetParent(); parent != nil {
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

func (actor *Actor) LocalToStageCoordinates(localX, localY float32) (float32, float32) {
	transform := actor.ComputeTransform()
	vec := transformCoordinate(localX, localY, transform)
	return vec.X(), vec.Y()
}

func (actor *Actor) StageToLocalCoordinates(stageX, stageY float32) (float32, float32) {
	transform := actor.ComputeTransform().Inv()
	vec := transformCoordinate(stageX, stageY, transform)
	return vec.X(), vec.Y()
}

func (actor *Actor) ComputeTransform() mgl32.Mat3 {
	if !actor.dirty && actor.Parent == nil {
		return actor.transform
	}

	transform := mgl32.Ident3()
	if actor.Parent != nil {
		transform = actor.Parent.ComputeTransform()
	}

	transform = transform.Mul3(mgl32.Translate2D(actor.X, actor.Y))
	if actor.Rotation != 0 {
		transform = transform.Mul3(mgl32.HomogRotate2D(mgl32.DegToRad(actor.Rotation)))
	}

	if actor.ScaleX != 1 || actor.ScaleY != 1 {
		transform = transform.Mul3(mgl32.Scale2D(actor.ScaleX, actor.ScaleY))
	}

	if actor.OriginX != 0 || actor.OriginY != 0 {
		transform = transform.Mul3(mgl32.Translate2D(-actor.OriginX, -actor.OriginY))
	}

	if actor.Parent == nil {
		actor.transform = transform
		actor.dirty = false
	}

	return transform
}

func (actor *Actor) Draw(batch *Graphics.Batch) {
	if !actor.Visible {
		return
	}
	if actor.Behavior != nil {
		actor.Behavior.Draw(actor, batch)
	}
	for _, child := range actor.Children {
		child.Draw(batch)
	}
}

func (actor *Actor) Act(delta float32) {
	if actor.Behavior != nil {
		actor.Behavior.Act(actor, delta)
	}
	for _, child := range actor.Children {
		child.Act(delta)
	}
}

func (actor *Actor) Hit(x, y float32) bool {
	if !actor.Visible {
		return false
	}
	localX, localY := actor.StageToLocalCoordinates(x, y)
	return localX >= 0 && localY >= 0 && localX < actor.Width && localY < actor.Height
}
