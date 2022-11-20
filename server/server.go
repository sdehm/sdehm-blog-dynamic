package server

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/sdehm/sdehm-blog-dynamic/api"
	"github.com/sdehm/sdehm-blog-dynamic/data"
)

type Server struct {
	logger            *log.Logger
	repo              data.Repo
	connections       []*connection
	connectionUpdates chan func()
	lastId            int
}

func Start(addr string, logger *log.Logger, repo data.Repo) error {
	server := &Server{
		logger:            logger,
		repo:              repo,
		connectionUpdates: make(chan func()),
	}
	http.Handle("/ws", server.wsHandler())

	go server.startConnectionUpdates()

	return http.ListenAndServe(addr, nil)
}

func (s *Server) wsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Println("new connection")
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		s.addConnection(conn)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Server) addConnection(c net.Conn) {
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

func (s *Server) removeConnection(c *connection) {
	s.connectionUpdates <- func() {
		for i, con := range s.connections {
			if con.id == c.id {
				s.connections = append(s.connections[:i], s.connections[i+1:]...)
				return
			}
		}
	}
}

func (s *Server) startConnectionUpdates() {
	for u := range s.connectionUpdates {
		u()
	}
}

func (s *Server) broadcast(m api.Message) {
	for _, c := range s.connections {
		err := c.send(m)
		if err != nil {
			s.logger.Println(err)
			s.removeConnection(c)
		}
	}
}
