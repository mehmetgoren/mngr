package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"mngr/utils"
	"net/http"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *Client) Push(s interface{}) error { //s is StreamingEvent
	if c == nil {
		log.Println("Something may be wrong with the client side, Client is nil")
		return errors.New("client is nil")
	}
	json, err := utils.SerializeJson(s)
	if err != nil {
		return err
	}
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		log.Println("Error while getting next writer. Err: ", err)
		return err
	}
	_, err = w.Write([]byte(json))
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

func CreateClient(hub *Hub, w http.ResponseWriter, r *http.Request) *Client {
	conn, err := WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to websocket connection: ", err)
		return nil
	}
	clientStreaming := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	clientStreaming.hub.register <- clientStreaming

	return clientStreaming
}
