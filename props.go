package spacewar

import (
	"github.com/grinova/classic2d-server/vmath"
)

type rocketProps struct {
	position       vmath.Vec2
	angle          float64
	linearVelocity vmath.Vec2
	owner          string
}
