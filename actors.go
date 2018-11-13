package spacewar

import (
	net "github.com/grinova/physicsnet-server"
)

const (
	rocketStartDistance = 0.075
	rocketStartVelocity = 0.2
)

type shipActor struct {
}

func (a shipActor) OnInit(c net.Controller, selfID net.ActorID, send net.Send, spawn net.Spawn, exit net.Exit) {
}

func (a shipActor) OnMessage(c net.Controller, m net.Message, send net.Send, spawn net.Spawn, exit net.Exit) {
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

func (a rocketActor) OnInit(c net.Controller, selfID net.ActorID, send net.Send, spawn net.Spawn, exit net.Exit) {
}

func (a rocketActor) OnMessage(c net.Controller, m net.Message, send net.Send, spawn net.Spawn, exit net.Exit) {
}
