package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	spacewar "github.com/grinova/spacewar-server"
)

var (
	port      = flag.String("p", "3000", "port to serve on")
	directory = flag.String("d", ".", "the directory of static file to host")
)

func main() {
	flag.Parse()

	server := spacewar.Server{}
	server.Start()

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.Handle("/", http.FileServer(http.Dir(*directory)))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		if _, err := server.Connect(c); err != nil {
			c.Close()
			log.Println("upgrade:", err)
			return
		}
	})
	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
