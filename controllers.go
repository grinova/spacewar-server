package spacewar

import (
	"time"

	"github.com/grinova/classic2d-server/vmath"
	physicsnet "github.com/grinova/physicsnet-server"
)

const (
	shipMaxForce         = 0.1
	shipMaxTorque        = 5
	shipDumpRotationCoef = 0.97
)

type shipController struct {
	physicsnet.BaseController
	thrust float64
	torque float64
}

func (c *shipController) Step(d time.Duration) {
	body := c.GetBody()
	force := vmath.Vec2{X: 0, Y: c.thrust}.Rotate(body.GetRot())
	body.ApplyForce(force.Mul(shipMaxForce))
	body.SetTorque(c.torque * shipMaxTorque)
	body.AngularVelocity *= shipDumpRotationCoef
}

func (c *shipController) getNewRocketProps() rocketProps {
	body := c.GetBody()
	position := vmath.Vec2{X: 0, Y: rocketStartDistance}.Rotate(body.GetRot()).Add(body.GetPosition())
	angle := body.GetAngle()
	linearVelocity := vmath.Vec2{X: 0, Y: rocketStartVelocity}.Rotate(body.GetRot()).Add(body.LinearVelocity)
	return rocketProps{position: position, angle: angle, linearVelocity: linearVelocity}
}

const (
	rocketForce = 1
)

type rocketController struct {
	physicsnet.BaseController
}

func (c *rocketController) Step(d time.Duration) {
	body := c.GetBody()
	force := vmath.Vec2{X: 0, Y: rocketForce}.Rotate(body.GetRot())
	body.ApplyForce(force)
}
