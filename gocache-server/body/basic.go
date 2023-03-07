package body

import (
	"bytes"
	"encoding/binary"
)

type String struct {
	string
}

func (s *String) GetBytes() []byte {
	return []byte(s.string)
}
func (s *String) Type() string {
	return "string"
}

type Integer struct {
	int
}

func (s *Integer) GetBytes() []byte {
	bytesbuff := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesbuff, binary.BigEndian, s.int)
	if err != nil {
		return nil
	}
	return bytesbuff.Bytes()
}
func (s *Integer) Type() string {
	return "integer"
}

type Boolean struct {
	bool
}

func (s *Boolean) GetBytes() []byte {
	bytesbuff := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesbuff, binary.BigEndian, s.bool)
	if err != nil {
		return nil
	}
	return bytesbuff.Bytes()
}
func (s *Boolean) Type() string {
	return "boolean"
}

type Float struct {
	float64
}

func (s *Float) GetBytes() []byte {
	bytesbuff := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesbuff, binary.BigEndian, s.float64)
	if err != nil {
		return nil
	}
	return bytesbuff.Bytes()
}
func (s *Float) Type() string {
	return "float"
}

// accept msg
type Message struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
	Act   int    `json:"act"`
}

// status replay
type ReplayStatus struct {
	Content    []byte `json:"content"`
	StatusCode int    `json:"code"`
}
