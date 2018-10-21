package main

import (
	"net"
)

type Lobby struct {
	clients  []*Client
	joins    chan net.Conn
	incoming chan []byte
	outgoing chan []byte
}

func NewLobby() *Lobby {
	lobby := &Lobby{
		clients:  make([]*Client, 0),
		joins:    make(chan net.Conn),
		incoming: make(chan []byte),
		outgoing: make(chan []byte),
	}

	lobby.Listen()

	return lobby
}

func (l *Lobby) Broadcast(data []byte) {
	for _, client := range l.clients {
		client.outgoing <- data
	}
}

func (l *Lobby) Join(conn net.Conn) {
	client := NewClient(conn)
	l.clients = append(l.clients, client)
	go func() {
		for {
			l.incoming <- <-client.incoming
		}
	}()
}

func (l *Lobby) Listen() {
	go func() {
		for {
			select {
			case data := <-l.incoming:
				l.Broadcast(data)
			case conn := <-l.joins:
				l.Join(conn)
			}
		}
	}()
}
