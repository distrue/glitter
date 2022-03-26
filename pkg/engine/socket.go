package engine

import (
	"context"
	"fmt"
	"github.com/distrue/glitter/pkg/util"
	"github.com/gorilla/websocket"
	"net/http"
)

type Socket struct {
	Sid        string
	OnMessage  func(msg string)
	Send       func(msg string) error
}

func New(w http.ResponseWriter, r *http.Request) *Socket {
	conn := &Socket{
		Sid:       util.RandomSid(),
		OnMessage: func(msg string) {},
	}
	conn.connectWebSocket(w, r)

	return conn
}

func (c *Socket) connectWebSocket(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	webSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	go c.webSocketListener(cancel, webSocket)
	go func() {
		select {
		case <-ctx.Done():
			webSocket.Close()
		}
	}()
	c.Send = func(msg string) error {
		return webSocket.WriteMessage(websocket.TextMessage, []byte("4" + msg))
	}
}

func (c *Socket) webSocketListener(cancel context.CancelFunc, conn *websocket.Conn) {
	closeMe := false
	for !closeMe {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		if messageType != 1 {
			continue
		}

		parser(
			string(p),
			func() {
				closeMe = true
			},
			func(ans string) {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(ans)); err != nil {
					fmt.Println(err)
				}
			},
			func(msg string) {
				c.OnMessage(msg)
			},
			func() {},
		)
	}
	cancel()
}

// TODO: consider multiple types on websocket messageType
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)
