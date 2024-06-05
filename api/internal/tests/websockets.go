package tests

import (
	"context"
	"io"
	"time"

	"nhooyr.io/websocket"
)

func DialWebsocket(url string, timeout time.Duration) (*websocket.Conn, string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ws, rs, err := websocket.Dial(ctx, url, nil)

	if err != nil && rs != nil {
		return nil, printJSON(rs.Body)
	}

	return ws, ""
}

func printJSON(body io.Reader) string {
	b, _ := io.ReadAll(body)
	return string(b)
}
