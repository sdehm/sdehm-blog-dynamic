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
		queryParams := r.URL.Query()
		path := queryParams.Get("path")
		s.logger.Println("path:", path)
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		s.addConnection(conn, path)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Server) addConnection(c net.Conn, path string) {
	s.connectionUpdates <- func() {
		s.lastId++
		conn := &connection{
			id:   s.lastId,
			conn: c,
			path: path,
		}
		go conn.receiver(s)
		s.connections = append(s.connections, conn)
		id := fmt.Sprint(conn.id)
		comments, err := s.repo.GetPost(path)
		if err != nil {
			s.logger.Println(err)
			return
		}
		commentsHtml := api.RenderPostComments(comments)
		conn.sendConnected(id, commentsHtml)
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

func (s *Server) broadcast(m api.Message, path string) {
	for _, c := range s.connections {
		if c.path != path {
			continue
		}
		err := c.send(m)
		if err != nil {
			s.logger.Println(err)
			s.removeConnection(c)
		}
	}
}
