package rooms

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Token    string
	pingConn *websocket.Conn
	msgConn  *websocket.Conn
	Alive    bool
}

func (c *Client) close() {
	c.pingConn.Close()
	c.msgConn.Close()
	c.Alive = false
}

/**
pings client and tries to recieve message
*/
func (c Client) closedConnection() bool {
	if c.pingConn == nil {
		return false
	}
	c.pingConn.WriteMessage(websocket.BinaryMessage, []byte{1})
	_, _, err := c.pingConn.ReadMessage()
	if err != nil {
		return true
	}
	return false
}

func (c Client) send(msg string) error {
	return c.msgConn.WriteMessage(websocket.TextMessage, []byte(msg))
}
