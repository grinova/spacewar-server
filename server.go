package spacewar

import (
	"fmt"
	"log"
	"math"

	physicsnet "github.com/grinova/physicsnet-server"
)

var shipProps = map[string]physicsnet.BodyCreateProps{
	"ship-a": physicsnet.BodyCreateProps{
		ID:       "ship-a",
		Position: physicsnet.Point{X: -0.5, Y: -0.5},
		Angle:    0,
	},
	"ship-b": physicsnet.BodyCreateProps{
		ID:       "ship-b",
		Position: physicsnet.Point{X: 0.5, Y: 0.5},
		Angle:    math.Pi,
	},
}

type systemProps struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// Server - сервер игры
type Server struct {
	physicsnet.Server
}

// Start - запуск сервера
func (s *Server) Start() {
	listener := physicsnet.ServerListener{
		OnServerStart: func(s *physicsnet.Server) {
			s.GetWorld().SetContactListener(contactListener{server: s})
			s.GetBodyRegistrator().Register("arena", func(v interface{}) interface{} {
				return createArenaBody(s.GetWorld(), v)
			})
			s.GetBodyRegistrator().Register("ship", func(v interface{}) interface{} {
				return createShipBody(s.GetWorld(), v)
			})
			s.GetBodyRegistrator().Register("rocket", func(v interface{}) interface{} {
				return createRocketBody(s.GetWorld(), v)
			})
			s.GetControllerRegistrator().Register("ship", createShipController)
			s.GetControllerRegistrator().Register("rocket", createRocketController)
			s.GetActorRegistrator().Register("ship", createShipActor)
			s.GetActorRegistrator().Register("rocket", createRocketActor)

			s.CreateEntity("arena", "arena", physicsnet.BodyCreateProps{})
			log.Println("Server start")
		},
		OnServerStop: func(s *physicsnet.Server) {
			log.Println("Server stop")
		},
		OnClientConnect: func(s *physicsnet.Server, id string, c *physicsnet.Client) error {
			if props, ok := shipProps[id]; ok {
				s.CreateEntity(id, "ship", props)
				c.SendSystemMessage(systemProps{Type: "user-name", Data: id})
				log.Printf("Client connect: id = %s\n", id)
				return nil
			}
			return fmt.Errorf("onClientConnect: didn't find initial properties for ship id `%s`", id)
		},
		OnClientDisconnect: func(s *physicsnet.Server, id string) {
			s.DestroyEntity(id)
			log.Printf("Client disconnect: id = %s\n", id)
		},
		OnEventMessage: func(s *physicsnet.Server, id string, m interface{}) bool {
			log.Printf("Event from %s: %s\n", id, m)
			return true
		},
		OnSystemMessage: func(s *physicsnet.Server, id string, m interface{}) bool {
			log.Printf("System from %s: %s\n", id, m)
			return true
		},
	}
	s.SetListener(listener)
	go s.Loop()
}
