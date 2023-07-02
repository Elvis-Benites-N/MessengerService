package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	ServerID string `json:"serverId"`
	Username string `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	ServerID string `json:"serverId"`
	Username string `json:"username"`
}

// writeMessage writes messages from the channel to the client's websocket connection.
func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Message { // Use range to iterate over channel, eliminates need for checking ok
		err := c.Conn.WriteJSON(message)
		if err != nil {
			log.Printf("error writing JSON message: %v", err) // Add more descriptive error message
			break                                             // Exit loop on error to close the connection
		}
	}
}

// readMessage reads messages from the client's websocket connection and broadcasts them to the hub.
func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break // Exit loop on error to close the connection
		}

		msg := &Message{
			Content:  string(m),
			ServerID: c.ServerID,
			Username: c.Username,
		}

		hub.Broadcast <- msg
	}
}
