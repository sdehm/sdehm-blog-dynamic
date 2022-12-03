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
	s := &Server{
		logger:            logger,
		repo:              repo,
		connectionUpdates: make(chan func()),
	}
	http.Handle("/ws", s.wsHandler())

	go s.startConnectionUpdates()
	s.logger.Printf("Listening on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) wsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		path := queryParams.Get("path")
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			s.logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		go s.addConnection(conn, path)
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

		commentsHtml, err := s.getCommentsHtml(path)
		if err != nil {
			s.logger.Println(err)
			conn.conn.Close()
			return
		}

		go conn.receiver(s)
		s.connections = append(s.connections, conn)
		id := fmt.Sprint(conn.id)
		conn.sendConnected(id, commentsHtml)
		s.logger.Printf("New connection: %s", id)
		go s.updateViewers(path)
	}
}

func (s *Server) removeConnection(c *connection) {
	s.connectionUpdates <- func() {
		for i, con := range s.connections {
			if con.id == c.id {
				s.connections = append(s.connections[:i], s.connections[i+1:]...)
				s.logger.Printf("Connection closed: %d", c.id)
				go s.updateViewers(c.path)
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

func (s *Server) getCommentsHtml(path string) (string, error) {
	post, err := s.repo.GetPost(path)
	if err != nil {
		return "", err
	}
	return api.RenderPostComments(*post), nil
}

func (s *Server) updateViewers(path string) {
	viewers := 0
	for _, c := range s.connections {
		if c.path == path {
			viewers++
		}
	}
	// strip first and last character from path
	id := "views_" + path[1:len(path)-1] + ".md"
	s.broadcast(&api.MorphData{
		Type: "morph",
		Id:   id,
		Html: api.RenderViewers(id, viewers),
	}, path)
}
