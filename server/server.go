package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"

	"github.com/gobwas/ws"
	"github.com/sdehm/sdehm-blog-dynamic/api"
	"github.com/sdehm/sdehm-blog-dynamic/data"
)

var postPath = regexp.MustCompile(`^/posts/[^/]+/$`)

type Server struct {
	logger            *log.Logger
	repo              data.Repo
	connections       []*connection
	connectionUpdates chan func()
	connectionCounts	map[string]int
	lastId            int
}

func Start(addr string, logger *log.Logger, repo data.Repo) error {
	s := &Server{
		logger:            logger,
		repo:              repo,
		connectionUpdates: make(chan func()),
		connectionCounts: make(map[string]int),
	}
	http.Handle("/ws", s.wsHandler())

	go s.startConnectionUpdates()
	s.logger.Printf("Listening on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) wsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch path := r.URL.Query().Get("path"); {
		case postPath.MatchString(path) || path == "/":
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				s.logger.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			go s.addConnection(conn, path)
		default:
			s.logger.Printf("Unknown path: %s", path)
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
		s.connectionCounts[path]++

		go conn.receiver(s)
		s.connections = append(s.connections, conn)
		id := fmt.Sprint(conn.id)

		if path == "/" {
			go s.updateAllViewers()
		} else {
			commentsHtml, err := s.getCommentsHtml(path)
			if err != nil {
				s.logger.Println(err)
				conn.conn.Close()
				return
			}
			conn.sendConnected(id, commentsHtml)
		}

		s.logger.Printf("New connection: %s", id)
		go s.updateViewers(path)
	}
}

func (s *Server) removeConnection(c *connection) {
	s.connectionUpdates <- func() {
		for i, con := range s.connections {
			if con.id == c.id {
				s.connections = append(s.connections[:i], s.connections[i+1:]...)
				s.connectionCounts[c.path]--
				if s.connectionCounts[c.path] == 0 {
					delete(s.connectionCounts, c.path)
				}
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
	viewers := s.connectionCounts[path]
	// strip first and last character from path
	id, ok := viewersId(path)
	if !ok {
		// invalid path for the viewer count, don't update
		return
	}
	s.broadcast(&api.MorphData{
		Type: "morph",
		Id:   id,
		Html: api.RenderViewers(id, viewers),
	}, path)
	s.broadcast(&api.MorphData{
		Type: "morph",
		Id:   id,
		Html: api.RenderViewers(id, viewers),
	}, "/")
}

func (s *Server) updateAllViewers() {
	for path := range s.connectionCounts {
		id, ok := viewersId(path)
		if !ok {
			// invalid path for the viewer count, don't update
			continue
		}
		viewers := s.connectionCounts[path]
		if viewers == 0 {
			continue
		}
		s.broadcast(&api.MorphData{
			Type: "morph",
			Id:   id,
			Html: api.RenderViewers(id, s.connectionCounts[path]),
		}, "/")
	}
}

// build the viewers count id from the path
// returns the id and a boolean indicating if the path yielded a valid id
func viewersId(path string) (string, bool) {
	if len(path) < 2 {
		return "", false
	}
	return "views_" + path[1:len(path)-1] + ".md", true
}
