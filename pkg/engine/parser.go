package engine

import "strings"

func parser(
	p string,
	closeHandler func(),
	pongHandler func(ans string),
	msgHandler func(msg string),
	upgradeHandler func(),
) {
	for separator := strings.Index(p, "\x1e"); separator != -1; p = p[separator+1:] {
		payload := p[:separator]
		switch payload[0] {
		case '1': // close
			closeHandler()
			return
		case '2': // ping
			if p[1:] == "probe" {
				pongHandler("3probe") // pong
			}
		case '4': // message
			msgHandler(payload[1:])
		case '5': // upgrade
			upgradeHandler()
			// '6' - noop, do not have to care about this
		}
	}
}
