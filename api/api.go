package api

import "encoding/json"

// Message interface
type Message interface {
	Marshal() ([]byte, error)
}

type MorphData struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Html string `json:"html"`
}

func (m *MorphData) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

type Connected struct {
	Type         string `json:"type"`
	ConnectionId int    `json:"connection_id"`
	Html         string `json:"html"`
}

func (c *Connected) Marshal() ([]byte, error) {
	c.Type = "connected"
	return json.Marshal(c)
}
