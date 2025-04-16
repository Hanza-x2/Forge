package Scene

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"github.com/go-gl/mathgl/mgl32"
)

type Actor struct {
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

func (actor *Actor) AddChild(child *Actor) {
	if child.Parent != nil {
		child.Parent.RemoveChild(child)
	}
	child.Parent = actor
	child.stage = actor.stage
	actor.Children = append(actor.Children, child)
}

func (actor *Actor) RemoveChild(child *Actor) bool {
	for i, c := range actor.Children {
		if c == child {
			actor.Children = append(actor.Children[:i], actor.Children[i+1:]...)
			child.Parent = nil
			child.stage = nil
			return true
		}
	}
	return false
}

func (actor *Actor) SetPosition(x, y float32) {
	actor.X, actor.Y = x, y
	actor.dirty = true
}

func (actor *Actor) SetSize(width, height float32) {
	actor.Width, actor.Height = width, height
	actor.dirty = true
}

func (actor *Actor) SetOrigin(originX, originY float32) {
	actor.OriginX, actor.OriginY = originX, originY
	actor.dirty = true
}

func (actor *Actor) SetScale(scaleX, scaleY float32) {
	actor.ScaleX, actor.ScaleY = scaleX, scaleY
	actor.dirty = true
}

func (actor *Actor) SetRotation(degrees float32) {
	actor.Rotation = degrees
	actor.dirty = true
}

func transformCoordinate(vecX, vecY float32, mat mgl32.Mat3) mgl32.Vec2 {
	x := vecX*mat[0] + vecY*mat[3] + mat[6]
	y := vecX*mat[1] + vecY*mat[4] + mat[7]
	return mgl32.Vec2{x, y}
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

	transform := actor.ComputeTransform()
	batch.PushTransform(transform)
	actor.DrawSelf(batch)

	for _, child := range actor.Children {
		child.Draw(batch)
	}

	batch.PopTransform()
}

func (actor *Actor) DrawSelf(batch *Graphics.Batch) {
	// Base actor doesn't draw anything
	// Override this in child types
}

func (actor *Actor) Act(delta float32) {
	actor.ActSelf(delta)
	for _, child := range actor.Children {
		child.Act(delta)
	}
}

func (actor *Actor) ActSelf(delta float32) {
	// Base actor doesn't act
	// Override this in child types
}

func (actor *Actor) Hit(x, y float32) bool {
	if !actor.Visible {
		return false
	}
	localX, localY := actor.StageToLocalCoordinates(x, y)
	return localX >= 0 && localY >= 0 && localX < actor.Width && localY < actor.Height
}
