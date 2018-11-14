package spacewar

import (
	"github.com/grinova/classic2d-server/dynamic"
	physicsnet "github.com/grinova/physicsnet-server"
)

type contactListener struct {
	server *physicsnet.Server
}

func (l contactListener) BeginContact(c *dynamic.Contact) {
	if userDataA, ok := c.BodyA.UserData.(UserData); ok {
		if userDataB, ok := c.BodyB.UserData.(UserData); ok {
			typeA := userDataA.Type
			typeB := userDataB.Type
			if typeA == "arena" && typeB == "rocket" ||
				typeA == "black-hole" && (typeB == "ship" || typeB == "rocket") {
				l.server.DestroyBody(c.BodyB)
				l.server.DestroyContact(c)
			} else if typeA == "ship" && typeB == "rocket" ||
				typeB == "ship" && typeA == "rocket" ||
				typeA == "rocket" && typeB == "rocket" {
				l.server.DestroyBody(c.BodyA)
				l.server.DestroyBody(c.BodyB)
				l.server.DestroyContact(c)
			}
		}
	}
}

func (l contactListener) EndContact(c *dynamic.Contact) {
}

func (l contactListener) PreSolve(c *dynamic.Contact) {
}
