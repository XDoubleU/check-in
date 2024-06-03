package tests

import (
	"context"
	"time"

	"nhooyr.io/websocket"
)

func DialWebsocket(url string, timeout time.Duration) *websocket.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ws, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil
	}

	return ws
}
