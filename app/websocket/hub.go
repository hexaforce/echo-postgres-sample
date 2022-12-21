// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package websocket

import "sync"

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {

	// Registered clients.
	clients map[string][]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// clients mutex
	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string][]*Client),
		mu:         sync.RWMutex{},
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			// h.clients[client] = true
			if _, ok := h.clients[client.userName]; !ok {
				h.clients[client.userName] = make([]*Client, 0)
			}
			// ｄelete if there is already one connected with the same userName
			for _, client := range h.clients[client.userName] {
				h.release(client)
			}
			h.clients[client.userName] = append(h.clients[client.userName], client)

		case client := <-h.unregister:
			h.release(client)

		case message := <-h.broadcast:
			for _, clients := range h.clients {
				h.send(clients, message)
			}
		}

	}
}

// client release
func (h *Hub) release(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, c := range h.clients[client.userName] {
		if c == client {
			// log.Println("release <<<<<<<<< :" + client.userName)
			slice := h.clients[client.userName]
			for i, c := range slice {
				if c == client {
					slice = append(slice[:i], slice[i+1:]...)
					temp := make([]*Client, len(slice))
					copy(temp, slice)
					client.close()
					h.clients[client.userName] = temp
				}
			}
		}
	}
}

// send chanel
func (h *Hub) send(clients []*Client, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range clients {
		select {
		case client.send <- message:
		default:
			h.release(client)
		}
	}
}
