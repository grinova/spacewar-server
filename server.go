package spacewar

import (
	"fmt"
	"log"
	"math"

	"github.com/grinova/classic2d-server/vmath"
	physicsnet "github.com/grinova/physicsnet-server"
)

var shipProps = map[string]physicsnet.BodyCreateProps{
	"ship-a": physicsnet.BodyCreateProps{
		ID:              "ship-a",
		Position:        physicsnet.Point{X: -0.5, Y: -0.5},
		Angle:           0,
		LinearVelocity:  physicsnet.Point{X: 0, Y: 0},
		AngularVelocity: 0,
	},
	"ship-b": physicsnet.BodyCreateProps{
		ID:              "ship-b",
		Position:        physicsnet.Point{X: 0.5, Y: 0.5},
		Angle:           math.Pi,
		LinearVelocity:  physicsnet.Point{X: 0, Y: 0},
		AngularVelocity: 0,
	},
}

type systemProps struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type scoreProps struct {
	PlayerName string `json:"playerName"`
	Amount     int    `json:"amount"`
}

type player struct {
	name  string
	score int
}

// Server - сервер игры
type Server struct {
	physicsnet.Server
	players map[string]*player
}

// Start - запуск сервера
func (server *Server) Start() {
	server.players = make(map[string]*player)
	listener := physicsnet.ServerListener{
		OnServerStart: func(s *physicsnet.Server) {
			s.Props.NewID = func() (string, error) {
				for i := 'a'; i < 'z'; i++ {
					id := "ship-" + string(i)
					if _, ok := server.players[id]; !ok {
						return id, nil
					}
				}
				return "", fmt.Errorf("genNewID: can't generate new id")
			}
			s.GetWorld().SetContactListener(contactListener{server: server})
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
		OnClientConnect: func(s *physicsnet.Server, id string, client *physicsnet.Client) error {
			if ok := server.createShip(id); !ok {
				return fmt.Errorf("OnClientConnect: can't create ship with id = `%s`", id)
			}
			client.SendSystemMessage(systemProps{Type: "user-name", Data: id})
			server.players[id] = &player{}
			log.Printf("Client connect: id = %s\n", id)
			return nil
		},
		OnClientDisconnect: func(s *physicsnet.Server, id string) {
			s.DestroyEntity(id)
			player := server.players[id]
			delete(server.players, id)
			for playerID, player := range server.players {
				if client := s.GetClient(playerID); client != nil {
					client.SendSystemMessage(systemProps{Type: "opponent-leave", Data: player.name})
				}
			}
			log.Printf("Client disconnect: name = %s, id = %s\n", player.name, id)
		},
		OnEventMessage: func(s *physicsnet.Server, id string, m interface{}) bool {
			log.Printf("Event from %s: %s\n", id, m)
			return true
		},
		OnSystemMessage: func(s *physicsnet.Server, id string, m interface{}) bool {
			if data, ok := m.(map[string]interface{}); ok {
				if data["type"] == "player-name" {
					if playerName, ok := data["data"].(string); ok {
						server.onReceivePlayerName(id, playerName)
					}
				}
			}
			log.Printf("System from %s: %s\n", id, m)
			return true
		},
	}
	server.SetListener(listener)
	go server.Loop()
}

func (server *Server) createShip(id string) bool {
	if props, ok := shipProps[id]; ok {
		server.CreateEntity(id, "ship", props)
		return true
	}
	return false
}

func (server *Server) resetShip(id string) bool {
	if props, ok := shipProps[id]; ok {
		body := server.GetBody(id)
		body.Sweep.A = props.Angle
		body.Sweep.C = vmath.Vec2{X: props.Position.X, Y: props.Position.Y}
		body.LinearVelocity = vmath.Vec2{X: props.LinearVelocity.X, Y: props.LinearVelocity.Y}
		body.AngularVelocity = props.AngularVelocity
		return true
	}
	return false
}

func (server *Server) onDestroyShip(destroyerID string, destroyedID string) {
	var amount int
	if destroyerID == destroyedID {
		amount = -1
	} else {
		amount = 1
	}
	if player, ok := server.players[destroyerID]; ok {
		for playerID := range server.players {
			c := server.GetClient(playerID)
			player.score++
			c.SendSystemMessage(systemProps{Type: "score", Data: scoreProps{PlayerName: player.name, Amount: amount}})
		}
	}
	for playerID, player := range server.players {
		if playerID == destroyedID {
			server.createShip(destroyedID)
		} else {
			server.resetShip(playerID)
		}
		player.score = 0
	}
}

func (server *Server) onReceivePlayerName(playerID string, playerName string) {
	if player, ok := server.players[playerID]; ok {
		player.name = playerName
		if len(server.players) > 1 {
			if client := server.GetClient(playerID); client != nil {
				for id, player := range server.players {
					if playerID != id {
						if c := server.GetClient(id); c != nil {
							c.SendSystemMessage(systemProps{Type: "opponent-join", Data: playerName})
						}
						client.SendSystemMessage(systemProps{Type: "opponent-join", Data: player.name})
					}
				}
			}
		}
	}
}
