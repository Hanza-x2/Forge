package Scene

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"github.com/go-gl/mathgl/mgl32"
)

type Actor interface {
	GetName() string
	SetName(name string)

	GetX() float32
	GetY() float32
	SetPosition(x, y float32)

	GetWidth() float32
	GetHeight() float32
	SetSize(width, height float32)

	GetOriginX() float32
	GetOriginY() float32
	SetOrigin(originX, originY float32)

	GetScaleX() float32
	GetScaleY() float32
	SetScale(scaleX, scaleY float32)

	GetRotation() float32
	SetRotation(degrees float32)

	GetZIndex() int
	SetZIndex(zIndex int)

	IsVisible() bool
	SetVisible(visible bool)

	GetUserData() interface{}
	SetUserData(data interface{})

	GetParent() Actor
	SetParent(parent Actor)

	GetStage() *Stage
	SetStage(stage *Stage)

	AddChild(child Actor)
	RemoveChild(child Actor) bool
	GetChildren() []Actor
	RemoveAllChildren()

	LocalToStageCoordinates(localX, localY float32) (float32, float32)
	StageToLocalCoordinates(stageX, stageY float32) (float32, float32)
	ComputeTransform() mgl32.Mat3

	Draw(batch *Graphics.Batch)
	DrawSelf(batch *Graphics.Batch)

	Act(delta float32)
	ActSelf(delta float32)

	Hit(x, y float32) bool
}

type BaseActor struct {
	Name      string
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
	Parent    Actor
	Children  []Actor
	UserData  interface{}
	stage     *Stage
	transform mgl32.Mat3
	dirty     bool
}

func NewBaseActor() *BaseActor {
	return &BaseActor{
		ScaleX:   1,
		ScaleY:   1,
		Visible:  true,
		Children: make([]Actor, 0),
		dirty:    true,
	}
}

func (actor *BaseActor) GetName() string {
	return actor.Name
}

func (actor *BaseActor) SetName(name string) {
	actor.Name = name
}

func (actor *BaseActor) GetX() float32 {
	return actor.X
}

func (actor *BaseActor) GetY() float32 {
	return actor.Y
}

func (actor *BaseActor) SetPosition(x, y float32) {
	actor.X, actor.Y = x, y
	actor.dirty = true
}

func (actor *BaseActor) GetWidth() float32 {
	return actor.Width
}

func (actor *BaseActor) GetHeight() float32 {
	return actor.Height
}

func (actor *BaseActor) SetSize(width, height float32) {
	actor.Width, actor.Height = width, height
	actor.dirty = true
}

func (actor *BaseActor) GetOriginX() float32 {
	return actor.OriginX
}

func (actor *BaseActor) GetOriginY() float32 {
	return actor.OriginY
}

func (actor *BaseActor) SetOrigin(originX, originY float32) {
	actor.OriginX, actor.OriginY = originX, originY
	actor.dirty = true
}

func (actor *BaseActor) GetScaleX() float32 {
	return actor.ScaleX
}

func (actor *BaseActor) GetScaleY() float32 {
	return actor.ScaleY
}

func (actor *BaseActor) SetScale(scaleX, scaleY float32) {
	actor.ScaleX, actor.ScaleY = scaleX, scaleY
	actor.dirty = true
}

func (actor *BaseActor) GetRotation() float32 {
	return actor.Rotation
}

func (actor *BaseActor) SetRotation(degrees float32) {
	actor.Rotation = degrees
	actor.dirty = true
}

func (actor *BaseActor) GetZIndex() int {
	return actor.ZIndex
}

func (actor *BaseActor) SetZIndex(zIndex int) {
	actor.ZIndex = zIndex
	actor.dirty = true
}

func (actor *BaseActor) IsVisible() bool {
	return actor.Visible
}

func (actor *BaseActor) SetVisible(visible bool) {
	actor.Visible = visible
}

func (actor *BaseActor) GetUserData() interface{} {
	return actor.UserData
}

func (actor *BaseActor) SetUserData(data interface{}) {
	actor.UserData = data
}

func (actor *BaseActor) GetParent() Actor {
	return actor.Parent
}

func (actor *BaseActor) SetParent(parent Actor) {
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

func (actor *BaseActor) GetStage() *Stage {
	return actor.stage
}

func (actor *BaseActor) SetStage(stage *Stage) {
	actor.stage = stage
	for _, child := range actor.Children {
		child.SetStage(stage)
	}
}

func (actor *BaseActor) AddChild(child Actor) {
	if currentParent := child.GetParent(); currentParent != nil {
		currentParent.RemoveChild(child)
	}
	actor.Children = append(actor.Children, child)
	child.SetParent(actor)
	child.SetStage(actor.stage)
}

func (actor *BaseActor) RemoveChild(child Actor) bool {
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

func (actor *BaseActor) GetChildren() []Actor {
	return actor.Children
}

func (actor *BaseActor) RemoveAllChildren() {
	for _, child := range actor.Children {
		child.SetParent(nil)
		child.SetStage(nil)
	}
	actor.Children = make([]Actor, 0)
}

func transformCoordinate(vecX, vecY float32, mat mgl32.Mat3) mgl32.Vec2 {
	x := vecX*mat[0] + vecY*mat[3] + mat[6]
	y := vecX*mat[1] + vecY*mat[4] + mat[7]
	return mgl32.Vec2{x, y}
}

func (actor *BaseActor) LocalToStageCoordinates(localX, localY float32) (float32, float32) {
	transform := actor.ComputeTransform()
	vec := transformCoordinate(localX, localY, transform)
	return vec.X(), vec.Y()
}

func (actor *BaseActor) StageToLocalCoordinates(stageX, stageY float32) (float32, float32) {
	transform := actor.ComputeTransform().Inv()
	vec := transformCoordinate(stageX, stageY, transform)
	return vec.X(), vec.Y()
}

func (actor *BaseActor) ComputeTransform() mgl32.Mat3 {
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

func (actor *BaseActor) Draw(batch *Graphics.Batch) {
	if !actor.Visible {
		return
	}

	transform := actor.ComputeTransform()
	batch.PushTransform(transform)
	actor.DrawSelf(batch)

	for _, child := range actor.Children {
		child.Draw(batch)
	}

	batch.PopTransform()
}

func (actor *BaseActor) DrawSelf(batch *Graphics.Batch) {
	// Base actor doesn't draw anything
	// Override this in child types
}

func (actor *BaseActor) Act(delta float32) {
	actor.ActSelf(delta)
	for _, child := range actor.Children {
		child.Act(delta)
	}
}

func (actor *BaseActor) ActSelf(delta float32) {
	// Base actor doesn't act
	// Override this in child types
}

func (actor *BaseActor) Hit(x, y float32) bool {
	if !actor.Visible {
		return false
	}
	localX, localY := actor.StageToLocalCoordinates(x, y)
	return localX >= 0 && localY >= 0 && localX < actor.Width && localY < actor.Height
}
