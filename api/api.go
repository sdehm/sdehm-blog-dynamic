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

type ConnectionId string

func (m ConnectionId) Marshal() ([]byte, error) {
	return json.Marshal( struct {
		Id ConnectionId `json:"id"`
	}{
		Id: m,
	})
}
