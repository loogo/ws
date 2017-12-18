package server

import (
	"fmt"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.connected(client)
		case client := <-h.unregister:
			h.disConnected(client)
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) connected(client *Client) {
	h.clients[client] = true
	fmt.Println(client.user.Code, " Connected!")
}

func (h *Hub) disConnected(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		fmt.Println(client.user.Code, " DisConnected!")
	}
}

func (h *Hub) GetUsers() (users []*User) {
	for _, v := range h.GetDistinct() {
		users = append(users, v.user)
	}
	return
}
func (h *Hub) GetDistinct() (clients []*Client) {
	for key, value := range h.clients {
		if value {
			duplicate := false
			for _, v := range clients {
				if v.user.Code == key.user.Code {
					duplicate = true
					break
				}
			}
			if !duplicate {
				clients = append(clients, key)
			}
		}
	}
	return
}
