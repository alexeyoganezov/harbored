package webserver

import (
	"sync"
	"time"
)

type WSMessage struct {
	Command string                 `json:"command"`
	Payload map[string]interface{} `json:"payload"`
}

type WSRoom struct {
	Name       string
	Clients    map[*WSClient]bool
	Broadcast  chan []byte
	Register   chan *WSClient
	Unregister chan *WSClient
}

func NewRoom(name string) *WSRoom {
	return &WSRoom{
		Name:       name,
		Clients:    make(map[*WSClient]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *WSClient),
		Unregister: make(chan *WSClient),
	}
}

func (room *WSRoom) init() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		// Add a client to the room
		case client := <-room.Register:
			room.Clients[client] = true
		// Remove a client from the room
		case client := <-room.Unregister:
			if _, ok := room.Clients[client]; ok {
				delete(room.Clients, client)
				close(client.Send)
			}
		// Remove disconnected clients
		case <-ticker.C:
			for client := range room.Clients {
				if client.Disconnected {
					delete(room.Clients, client)
					close(client.Send)
				}
			}
		// Send a message to all the clients
		case message := <-room.Broadcast:
			for client := range room.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(room.Clients, client)
				}
			}
		}
	}
}

// A list of rooms
type WSRooms struct {
	sync.RWMutex
	items map[string]*WSRoom
}

// Create a room
func (rooms *WSRooms) Set(key string, value *WSRoom) {
	rooms.Lock()
	defer rooms.Unlock()
	rooms.items[key] = value
}

// Get a room
func (rooms *WSRooms) Get(key string) (*WSRoom, bool) {
	rooms.Lock()
	defer rooms.Unlock()
	value, ok := rooms.items[key]
	return value, ok
}

// Remove a room
func (rooms *WSRooms) Delete(key string) {
	rooms.Lock()
	defer rooms.Unlock()
	delete(rooms.items, key)
}

var RoomList = &WSRooms{
	items: make(map[string]*WSRoom),
}
