package spacewar

import (
	"github.com/grinova/classic2d-server/dynamic"
	"github.com/grinova/classic2d-server/physics"
	"github.com/grinova/classic2d-server/physics/shapes"
	"github.com/grinova/classic2d-server/vmath"
	physicsnet "github.com/grinova/physicsnet-server"
)

func createArenaBody(w *dynamic.World, v interface{}) interface{} {
	bodyDef := physics.BodyDef{Inverse: true}
	body := w.CreateBody(bodyDef)
	shape := shapes.CircleShape{Radius: 1}
	fixtureDef := physics.FixtureDef{Shape: shape, Density: 1}
	body.SetFixture(fixtureDef)
	body.UserData = UserData{Type: "arena", ID: "arena"}
	body.Type = physics.StaticBody
	return body
}

func createShipBody(w *dynamic.World, v interface{}) interface{} {
	if props, ok := v.(physicsnet.BodyCreateProps); ok {
		RADIUS := 0.05
		bodyDef := physics.BodyDef{
			Position: vmath.Vec2{X: props.Position.X, Y: props.Position.Y},
			Angle:    props.Angle,
		}
		body := w.CreateBody(bodyDef)
		shape := shapes.CircleShape{Radius: RADIUS}
		fixtureDef := physics.FixtureDef{Shape: shape, Density: 1}
		body.SetFixture(fixtureDef)
		body.UserData = UserData{ID: props.ID, Type: "ship"}
		return body
	}
	return nil
}

func createRocketBody(w *dynamic.World, v interface{}) interface{} {
	if props, ok := v.(rocketProps); ok {
		RADIUS := 0.01
		bodyDef := physics.BodyDef{
			Position:       props.position,
			Angle:          props.angle,
			LinearVelocity: props.linearVelocity,
		}
		body := w.CreateBody(bodyDef)
		shape := shapes.CircleShape{Radius: RADIUS}
		fixtureDef := physics.FixtureDef{Shape: shape, Density: 1}
		body.SetFixture(fixtureDef)
		body.UserData = UserData{Type: "rocket", Owner: props.owner}
		return body
	}
	return nil
}

func createShipController(v interface{}) interface{} {
	return &shipController{}
}

func createShipActor(v interface{}) interface{} {
	return shipActor{}
}

func createRocketController(v interface{}) interface{} {
	return &rocketController{}
}

func createRocketActor(v interface{}) interface{} {
	return rocketActor{}
}
