package main

import (
	"fmt"
	"github.com/distrue/glitter/pkg/engine"
	"net/http"
)

func main() {
	var connMap map[string]*engine.Socket
	onConnection := func(socket *engine.Socket) {
		fmt.Println("connection is created")
	}

	// engine attach example
	mux := http.NewServeMux()
	mux.Handle("/engine.io/", engine.Handler(connMap, onConnection))

	http.ListenAndServe(":8080", mux)
}
