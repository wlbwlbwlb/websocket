// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"sync"
	"sync/atomic"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	//clients map[*Client]bool
	clients ClientSet

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Exit requests from server shutting down.
	exit chan []byte

	open atomic.Bool
}

func newHub() *Hub {
	return &Hub{
		exit:       make(chan []byte),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		//clients:    make(map[*Client]bool),
		clients: ClientSet{
			members: make(map[*Client]bool),
		},
	}
}

func (p *Hub) run() {
	p.open.Store(true)
	for {
		select {
		case client := <-p.register:
			//h.clients[client] = true
			p.clients.add(client)
		case client := <-p.unregister:
			//if _, ok := h.clients[client]; ok {
			//	delete(h.clients, client)
			//	close(client.send)
			//}
			p.clients.del(client)
		case message := <-p.broadcast:
			//for client := range h.clients {
			//	select {
			//	case client.send <- message:
			//	default:
			//		close(client.send)
			//		delete(h.clients, client)
			//	}
			//}
			p.clients.each(func(client *Client) {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(p.clients.members, client)
				}
			})
		case reason := <-p.exit:
			//for client := range h.clients {
			//	select {
			//	case client.send <- reason:
			//	default:
			//	}
			//	close(client.send)
			//	delete(h.clients, client)
			//}
			p.open.Store(false)
			p.clients.each(func(client *Client) {
				select {
				case client.send <- reason:
				default:
				}
				close(client.send)
				delete(p.clients.members, client)
			})
		}
	}
}

func (p *Hub) closed() bool {
	return !p.open.Load()
}

type ClientSet struct {
	members map[*Client]bool
	sync.RWMutex
}

func (p *ClientSet) add(client *Client) {
	p.Lock()
	defer p.Unlock()

	p.members[client] = true
}

func (p *ClientSet) del(client *Client) {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.members[client]; ok {
		close(client.send)
		delete(p.members, client)
	}
}

func (p *ClientSet) each(f func(*Client)) {
	p.RLock()
	defer p.RUnlock()

	for o, _ := range p.members {
		f(o)
	}
}

func (p *ClientSet) len() int {
	p.RLock()
	defer p.RUnlock()

	return len(p.members)
}
