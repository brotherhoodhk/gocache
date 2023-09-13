package body

import (
	"fmt"
	"gocache/basic"
	"strconv"
	"sync"
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
	origin := strconv.Itoa(s.int)
	return []byte(origin)
}
func (s *Integer) Type() string {
	return "integer"
}

type Boolean struct {
	bool
}

func (s *Boolean) GetBytes() []byte {
	return []byte(strconv.FormatBool(s.bool))
}
func (s *Boolean) Type() string {
	return "boolean"
}

type Float struct {
	float64
}

func (s *Float) GetBytes() []byte {
	return []byte(fmt.Sprintf("%f", s.float64))
}
func (s *Float) Type() string {
	return "float"
}

type CustomDb struct {
	Cellmap         map[string]*basic.Cell
	MapContaierSize int
	Name            string
	Mutex           sync.RWMutex
}

func NewCustomDB(dbname string) *CustomDb {
	return &CustomDb{Cellmap: make(map[string]*basic.Cell), MapContaierSize: 5, Name: dbname}
}

// accept msg
type Message struct {
	DB    string `json:"db"`
	Key   string `json:"key"`
	Value []byte `json:"value"`
	Act   int    `json:"act"`
}

// status replay
type ReplayStatus struct {
	Content    []byte `json:"content"`
	StatusCode int    `json:"code"`
	Type       string `json:"type"`
}

// 应用于批量数据
type ReplayStatusVtwo struct {
	Content    map[string]any `json:"content"`
	StatusCode int            `json:"code"`
}

// set default fuzzy match method
var Default_Fuzzy_Match func(string, *CustomDb) []byte = FuzzyMatch
