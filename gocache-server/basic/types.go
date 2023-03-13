package basic

// accept msg
type Message struct {
	DB    string `json:"db"`
	Key   string `json:"key"`
	Value []byte `json:"value"`
	Act   int    `json:"act"`
}
