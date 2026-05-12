package main

import (
	"log"
	"net/http"

	"github.com/licheng1013/rocket-cat/modules/chat"
	"github.com/licheng1013/rocket-cat/modules/room"
	"github.com/licheng1013/rocket-cat/modules/user"
	"github.com/licheng1013/rocket-cat/ws"
)

func main() {
	router := ws.NewRouter()

	user.Register(router)
	chat.Register(router)
	room.Register(router)

	http.Handle("/ws", ws.NewServer(router))

	log.Println("rocket-cat listening on :10100")
	log.Fatal(http.ListenAndServe(":10100", nil))
}
