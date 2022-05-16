package ws

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"mngr/reps"
	"net/http"
	"sync"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id string

	hub *Hub

	mu sync.Mutex

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *Client) Push(s interface{}) error {
	if c == nil {
		log.Println("Something may be wrong with the client side, Client is nil")
		return errors.New("client is nil")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	jb, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		//log.Println("Error while getting next writer. Err: ", err)
		return err
	}
	_, err = w.Write(jb)
	if err != nil {
		log.Println("Error while writing to writer. Err: ", err)
		return err
	}

	if err := w.Close(); err != nil {
		log.Println("Error while closing writer. Err: ", err)
		return err
	}

	return nil
}

func readLoop(client *Client) {
	for {
		conn := client.conn
		if _, _, err := conn.NextReader(); err != nil {
			client.hub.unregister <- client
			err := conn.Close()
			if err != nil {
				log.Println("Error while closing websockets connection. Err: ", err)
				return
			}
			break
		}
	}
}

func CreateClient(hub *Hub, w http.ResponseWriter, r *http.Request) *Client {
	conn, err := WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to websocket connection: ", err)
		return nil
	}
	//conn.SetCloseHandler(func(code int, text string) error {
	//	log.Println("Client connection closed with code: ", code, " and text: ", text)
	//	return nil
	//})
	clientStream := &Client{id: reps.NewId(), hub: hub, conn: conn, send: make(chan []byte, 256)}
	clientStream.hub.register <- clientStream
	go readLoop(clientStream)

	return clientStream
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		log.Println("Error while closing websockets connection. Err: ", err)
	}
	return err
}
