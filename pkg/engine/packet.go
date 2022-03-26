package engine

import (
	"strings"
)

func parser(
	p string,
	closeHandler func(),
	pongHandler func(ans string),
	msgHandler func(msg string),
	upgradeHandler func(),
) {
	for separator := strings.Index(p, "\x1e"); separator != -1; p = p[separator+1:] {
		payload := p[:separator]
		switch []byte(payload)[0] - 48 {
		case CLOSE_PACKET: // close
			closeHandler()
			return
		case PING_PACKET: // ping
			if p[1:] == "probe" {
				pongHandler("3probe") // pong
			}
		case MESSAGE_PACKET: // message
			msgHandler(payload[1:])
		case UPGRADE_PACKET: // upgrade
			upgradeHandler()
			// '6' - noop, do not have to care about this
		}
	}
}

const (
	OPEN_PACKET = iota
	CLOSE_PACKET
	PING_PACKET
	PONG_PACKET
	MESSAGE_PACKET
	UPGRADE_PACKET
	NOOP_PACKET
)

// TODO: add error handling, unify to packet logic

const (
	UNKNOWN_TRANSPORT = iota
	UNKNOWN_SID
	BAD_HANDSHAKE_METHOD
	BAD_REQUEST
	FORBIDDEN
	UNSUPPORTED_PROTOCOL_VERSION
)

var errorMessage = map[int]string{
	0: "Transport unknown",
	1: "Session ID unknown",
	2: "Bad handshake method",
	3: "Bad request",
	4: "Forbidden",
	5: "Unsupported protocol version",
}
