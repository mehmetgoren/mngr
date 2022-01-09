package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"mngr/utils"
	"net/http"
)

func PushStreamServerInfo(s interface{}) { //s is StreamingEvent
	if clientStreaming == nil {
		log.Println("Something may be wrong with the client side, clientStreaming is nil")
		return
	}
	json, err := utils.SerializeJson(s)
	if err != nil {
		return
	}
	w, err := clientStreaming.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	w.Write([]byte(json))

	if err := w.Close(); err != nil {
		return
	}
}

var clientStreaming *Client

func HandlerStreaming(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	clientStreaming = &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	clientStreaming.hub.register <- clientStreaming
}
