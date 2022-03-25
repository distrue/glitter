package main

import (
	"fmt"
	"github.com/distrue/glitter/pkg/engine"
	"net/http"
)

func main() {
	var connMap map[string]*engine.SocketConn
	onConnection := func(conn *engine.SocketConn) {
		fmt.Println("connection is created")
	}

	mux := http.NewServeMux()
	mux.Handle("/engine.io/", engine.Handler(connMap, onConnection))

	http.ListenAndServe(":8080", mux)
}
