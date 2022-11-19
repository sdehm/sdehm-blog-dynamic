package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/gobwas/ws/wsutil"
)


type connection struct {
	id   int
	conn net.Conn
}

func (c *connection) sendConnected(id string) {
	c.send(connectionId(id))
}

// Serialize the data to JSON and send it to the client
func (c *connection) send(m message) error {
	if c.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	data, err := m.marshal()
	if err != nil {
		return err
	}

	err = wsutil.WriteServerText(c.conn, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *connection) receiver(s *server) {
	for {
		data, _, err := wsutil.ReadClientData(c.conn)
		if err != nil {
			return
		}

		xy := struct {
			X  int    `json:"x"`
			Y  int    `json:"y"`
			Id string `json:"id"`
		}{}
		err = json.Unmarshal(data, &xy)
		if err != nil {
			return
		}

		// s.broadcast(morphData{
		// 	Id:   "cursor_" + xy.Id,
		// 	Html: fmt.Sprintf("<div id=\"cursor_%s\" class=\"cursor\" style=\"--x: %d; --y: %d;\">%[1]s</div>", xy.Id, xy.X, xy.Y),
		// })
	}
}