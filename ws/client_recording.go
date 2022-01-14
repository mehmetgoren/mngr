package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"mngr/utils"
	"net/http"
)

func PushRecordServerInfo(s interface{}) { //s is RecordingEvent
	if clientRecording == nil {
		log.Println("Something may be wrong with the client side, clientRecording is nil")
		return
	}
	json, err := utils.SerializeJson(s)
	if err != nil {
		return
	}
	w, err := clientRecording.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	w.Write([]byte(json))

	if err := w.Close(); err != nil {
		return
	}
}

var clientRecording *Client

func HandlerRecording(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	clientRecording = &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	clientRecording.hub.register <- clientRecording
}
