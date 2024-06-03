package tests

import (
	"context"
	"time"

	"nhooyr.io/websocket"
)

func DialWebsocket(url string, timeout time.Duration) *websocket.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ws, rs, err := websocket.Dial(ctx, url, nil)
	if rs != nil || err != nil {
		return nil
	}

	return ws
}
