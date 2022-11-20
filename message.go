package main

import "encoding/json"

// message interface
type message interface {
	marshal() ([]byte, error)
}

type morphData struct {
	Id   string `json:"id"`
	Html string `json:"html"`
}

func (m *morphData) marshal() ([]byte, error) {
	return json.Marshal(m)
}

type connectionId string

func (m connectionId) marshal() ([]byte, error) {
	return json.Marshal( struct {
		Id connectionId `json:"id"`
	}{
		Id: m,
	})
}
