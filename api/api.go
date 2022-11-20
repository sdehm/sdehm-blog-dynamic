package api

import "encoding/json"

// Message interface
type Message interface {
	Marshal() ([]byte, error)
}

type MorphData struct {
	Id   string `json:"id"`
	Html string `json:"html"`
}

func (m *MorphData) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

type Connected struct {
	ConnectionId int `json:"connection_id"`
	Html string `json:"html"`
}

func (c *Connected) Marshal() ([]byte, error) {
	return json.Marshal(c)
}
