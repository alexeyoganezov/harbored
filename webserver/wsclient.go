package webserver

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type WSClient struct {
	Connection   *websocket.Conn
	Send         chan []byte
	Disconnected bool
}

func NewClient(ws *websocket.Conn) *WSClient {
	return &WSClient{
		Connection:   ws,
		Send:         make(chan []byte),
		Disconnected: false,
	}
}

// Start to listen for client messages
func (client *WSClient) read() {
	defer client.close()
	for {
		var message WSMessage
		err := client.Connection.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("IsUnexpectedCloseError error: %v", err)
			}
			break
		}
		if message.Command == "presentation:subscribe" {
			room, _ := RoomList.Get("global")
			room.Register <- client
		}
	}
}

// Close WebSocket connection for a client
func (client *WSClient) close() {
	client.Connection.Close()
	client.Disconnected = true
}
