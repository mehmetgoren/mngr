package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func (c *Client) writeStreamServerInitialized() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			//fmt.Println(message)
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//	w.Write(newline)
			//	w.Write(<-c.send)
			//}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func HandlerStreaming(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writeStreamServerInitialized()
}
