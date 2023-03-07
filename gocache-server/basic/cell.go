package basic

import "time"

type Cell struct {
	Value Cache  `json:"value"`
	Type  string `json:"type"`
	Time  string `json:"time"`
}

func NewCell(content CommonType) *Cell {
	return &Cell{Value: newcache(content.GetBytes()), Type: content.Type(), Time: time.Now().Format(time.RFC3339)}
}
func (s *Cell) GetValue() []byte {
	return s.Value.Value
}
