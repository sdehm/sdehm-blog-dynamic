package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
)

type server struct {
	logger            *log.Logger
	connections       []*connection
	connectionUpdates chan func()
	lastId            int
}

func start(addr string, logger *log.Logger) error {
	server := &server{
		logger:            logger,
		connectionUpdates: make(chan func()),
	}
	http.Handle("/ws", server.wsHandler())

	go server.startConnectionUpdates()

	return http.ListenAndServe(addr, nil)
}

func (s *server) wsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Println("new connection")
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		s.addConnection(conn)
		if err != nil {
			panic(err)
		}
	}
}

func (s *server) addConnection(c net.Conn) {
	s.connectionUpdates <- func() {
		s.lastId++
		conn := &connection{
			id:   s.lastId,
			conn: c,
		}
		go conn.receiver(s)
		s.connections = append(s.connections, conn)
		id := fmt.Sprint(conn.id)
		conn.sendConnected(id)
	}
}

func (s *server) removeConnection(c *connection) {
	s.connectionUpdates <- func() {
		for i, con := range s.connections {
			if con.id == c.id {
				s.connections = append(s.connections[:i], s.connections[i+1:]...)
				return
			}
		}
	}
}

func (s *server) startConnectionUpdates() {
	for u := range s.connectionUpdates {
		u()
	}
}

func (s *server) broadcast(m *morphData) {
	for _, c := range s.connections {
		err := c.send(m)
		if err != nil {
			s.logger.Println(err)
			s.removeConnection(c)
		}
	}
}