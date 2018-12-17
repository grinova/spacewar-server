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
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type scoreProps struct {
	PlayerName string `json:"playerName"`
	Amount     int    `json:"amount"`
}

// Server - сервер игры
type Server struct {
	physicsnet.Server
	clients map[string]string
}

// Start - запуск сервера
func (server *Server) Start() {
	server.clients = make(map[string]string)
	listener := physicsnet.ServerListener{
		OnServerStart: func(s *physicsnet.Server) {
			s.Props.NewID = func() (string, error) {
				for i := 'a'; i < 'z'; i++ {
					id := "ship-" + string(i)
					if _, ok := server.clients[id]; !ok {
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
		OnClientConnect: func(s *physicsnet.Server, id string, c *physicsnet.Client) error {
			if props, ok := shipProps[id]; ok {
				s.CreateEntity(id, "ship", props)
				c.SendSystemMessage(systemProps{Type: "user-name", Data: id})
				server.clients[id] = ""
				log.Printf("Client connect: id = %s\n", id)
				return nil
			}
			return fmt.Errorf("onClientConnect: didn't find initial properties for ship id `%s`", id)
		},
		OnClientDisconnect: func(s *physicsnet.Server, id string) {
			opponentName := server.clients[id]
			for clientID := range server.clients {
				if clientID != id {
					c := s.GetClient(clientID)
					c.SendSystemMessage(systemProps{Type: "opponent-leave", Data: opponentName})
				}
			}
			delete(server.clients, id)
			s.DestroyEntity(id)
			log.Printf("Client disconnect: id = %s\n", id)
		},
		OnEventMessage: func(s *physicsnet.Server, id string, m interface{}) bool {
			log.Printf("Event from %s: %s\n", id, m)
			return true
		},
		OnSystemMessage: func(s *physicsnet.Server, id string, m interface{}) bool {
			if data, ok := m.(map[string]interface{}); ok {
				if data["type"] == "player-name" {
					if opponentName, ok := data["data"].(string); ok {
						server.clients[id] = opponentName
						if len(server.clients) > 1 {
							opponent := s.GetClient(id)
							for clientID, clientName := range server.clients {
								if clientID != id {
									c := s.GetClient(clientID)
									if c != nil {
										c.SendSystemMessage(systemProps{Type: "opponent-join", Data: opponentName})
									}
									opponent.SendSystemMessage(systemProps{Type: "opponent-join", Data: clientName})
								}
							}
						}
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

func (server *Server) onDestroyShip(destroyerID string, destroyedID string) {
	var amount int
	if destroyerID == destroyedID {
		amount = -1
	} else {
		amount = 1
	}
	playerName := server.clients[destroyerID]
	for clientID := range server.clients {
		c := server.GetClient(clientID)
		c.SendSystemMessage(systemProps{Type: "score", Data: scoreProps{PlayerName: playerName, Amount: amount}})
	}
}
