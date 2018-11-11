package spacewar

import (
	"time"

	"github.com/grinova/classic2d-server/physics"
	"github.com/grinova/classic2d-server/vmath"
)

const (
	shipMaxForce         = 0.1
	shipMaxTorque        = 5
	shipDumpRotationCoef = 0.97
)

type shipController struct {
	thrust float64
	torque float64
	// FIXME: Костыль
	body *physics.Body
}

func (c *shipController) Step(body *physics.Body, d time.Duration) {
	c.body = body
	force := vmath.Vec2{X: 0, Y: c.thrust}.Rotate(body.GetRot())
	body.ApplyForce(force.Mul(shipMaxForce))
	body.SetTorque(c.torque * shipMaxTorque)
	body.AngularVelocity *= shipDumpRotationCoef
}

func (c *shipController) getNewRocketProps() rocketProps {
	ship := c.body
	position := vmath.Vec2{X: 0, Y: rocketStartDistance}.Rotate(ship.GetRot()).Add(ship.GetPosition())
	angle := ship.GetAngle()
	linearVelocity := vmath.Vec2{X: 0, Y: rocketStartVelocity}.Rotate(ship.GetRot()).Add(ship.LinearVelocity)
	return rocketProps{position: position, angle: angle, linearVelocity: linearVelocity}
}

const (
	rocketForce = 1
)

type rocketController struct {
}

func (c rocketController) Step(body *physics.Body, d time.Duration) {
	force := vmath.Vec2{X: 0, Y: rocketForce}.Rotate(body.GetRot())
	body.ApplyForce(force)
}
