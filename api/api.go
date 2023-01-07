package api

import "encoding/json"

type Message interface {
	Marshal() ([]byte, error)
}

type messageData struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Html string `json:"html"`
}

type Morph struct {
	messageData
}

func NewMorph(id, html string) *Morph {
	return &Morph{
		messageData: messageData{
			Id:   id,
			Html: html,
		},
	}
}

func (m *Morph) Marshal() ([]byte, error) {
	m.Type = "morph"
	return json.Marshal(m)
}

type Prepend struct {
	messageData
}

func NewPrepend(id, html string) *Prepend {
	return &Prepend{
		messageData: messageData{
			Id:   id,
			Html: html,
		},
	}
}

func (p *Prepend) Marshal() ([]byte, error) {
	p.Type = "prepend"
	return json.Marshal(p)
}

type Connected struct {
	messageData
}

func NewConnected(id, html string) *Connected {
	return &Connected{
		messageData: messageData{
			Id:   id,
			Html: html,
		},
	}
}

func (c *Connected) Marshal() ([]byte, error) {
	c.Type = "connected"
	return json.Marshal(c)
}
