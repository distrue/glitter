package engine

import (
	"encoding/json"
	"net/http"
)

type HandShake struct {
	Sid          string   `json:"sid"`
	Upgrades     []string `json:"upgrades"`
	PingTimeout  int      `json:"pingTimeout"`
	PingInterval int      `json:"pingInterval"`
}

func Handler(connMap map[string]*SocketConn, onConnection func(conn *SocketConn)) http.HandlerFunc {
	// only support websocket connection initiation; doest not support sending message

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			query := r.URL.Query()
			eioVersion := query.Get("EIO")
			if eioVersion != "4" {
				w.WriteHeader(400)
				w.Write([]byte("EIO version 4 is only supported"))
				return
			}

			transport := query.Get("transport")
			if transport == "websocket" {
				socketConn := New(w, r)
				connMap[socketConn.Sid] = socketConn
				onConnection(socketConn)
				// TODO: remove SocketConn from connMap on close channel
			}

			// open
			handShake := &HandShake{
				Sid: "temporary",
				Upgrades: []string{"websocket"},
				PingInterval: 25000,
				PingTimeout: 5000,
			}
			res, err := json.Marshal(handShake)
			if err != nil {
				panic(err)
			}
			w.WriteHeader(200)
			w.Write([]byte{'0'})
			w.Write(res)
		} else if r.Method == "POST" {
			var body string
			_, err := r.Body.Read([]byte(body))
			if err != nil {
				panic(err)
			}
			parser(
				body,
				func() {},
				func (msg string) {
					if msg == "" {
						w.WriteHeader(204)
						return
					}
					w.WriteHeader(200)
					w.Write([]byte(msg))
				},
				func(msg string) {
					// only support initialize connection; ignore message
					w.WriteHeader(200)
					w.Write([]byte("ok"))
				},
				func() {
					socketConn := New(w, r)
					connMap[socketConn.Sid] = socketConn
					onConnection(socketConn)
				},
			)
		}
	}
}
