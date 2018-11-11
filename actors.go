package spacewar

import (
	"github.com/grinova/actors"
	physicsnet "github.com/grinova/physicsnet-server"
)

const (
	rocketStartDistance = 0.075
	rocketStartVelocity = 0.2
)

type shipActor struct {
}

func (a shipActor) OnInit(c physicsnet.Controller, selfID actors.ActorID, send actors.Send, spawn physicsnet.Spawn, exit actors.Exit) {
}

func (a shipActor) OnMessage(c physicsnet.Controller, m actors.Message, send actors.Send, spawn physicsnet.Spawn, exit actors.Exit) {
	if ship, ok := c.(*shipController); ok {
		if msg, ok := m.(map[string]interface{}); ok {
			switch msg["type"] {
			case "thrust":
				if thrust, ok := msg["amount"].(float64); ok {
					ship.thrust = thrust
				}
			case "torque":
				if torque, ok := msg["amount"].(float64); ok {
					ship.torque = torque
				}
			case "fire":
				spawn("rocket", ship.getNewRocketProps())
			}
		}
	}
}

type rocketActor struct {
}

func (a rocketActor) OnInit(c physicsnet.Controller, selfID actors.ActorID, send actors.Send, spawn physicsnet.Spawn, exit actors.Exit) {
}

func (a rocketActor) OnMessage(c physicsnet.Controller, m actors.Message, send actors.Send, spawn physicsnet.Spawn, exit actors.Exit) {
}
