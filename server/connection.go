package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net"

	"github.com/gobwas/ws/wsutil"
	"github.com/sdehm/sdehm-blog-dynamic/api"
)

type connection struct {
	id   int
	conn net.Conn
	path string
}

// Send a message to the client to indicate that the connection was successful
func (c *connection) sendConnected(id string, commentsHtml string) {
	c.send(&api.Connected{
		ConnectionId: c.id,
		Html:         commentsHtml,
	})
}

// Serialize the data to JSON and send it to the client
func (c *connection) send(m api.Message) error {
	if c.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	data, err := m.Marshal()
	if err != nil {
		return err
	}

	err = wsutil.WriteServerText(c.conn, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *connection) receiver(s *Server) {
	defer c.conn.Close()
	
	for {
		// wait for semaphore availability
		s.receiveSemaphore <- struct{}{}
		data, _, err := wsutil.ReadClientData(c.conn)
		if err != nil {
			<-s.receiveSemaphore
			continue
		}

		commentData := struct {
			Type    string `json:"type"`
			Author  string `json:"author"`
			Comment string `json:"comment"`
		}{}
		err = json.Unmarshal(data, &commentData)
		if err != nil || commentData.Type != "comment" {
			s.logger.Println("Invalid data received from client")
			<-s.receiveSemaphore
			continue
		}
		comment, err := s.repo.AddComment(c.path, sanitize(commentData.Author), sanitize(commentData.Comment))
		if err != nil {
			s.logger.Println(err)
			<-s.receiveSemaphore
			continue
		}
		go s.broadcast(&api.MorphData{
			Type: "prepend",
			Id:   "comments",
			Html: api.RenderComment(comment),
		}, c.path)
		<-s.receiveSemaphore
	}
}

// Strip out html to sanitize the inputs
func sanitize(s string) string {
	return template.HTMLEscapeString(s)
}
