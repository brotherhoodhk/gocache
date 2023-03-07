package basic

import (
	"math/rand"
	"time"
)

// the bottom of all cache system
type Cache struct {
	Value []byte
	Key   uint32
}

// use this method to create a basic cache
func newcache(value []byte) Cache {
	rand.Seed(time.Now().UnixNano())
	return Cache{Value: value, Key: rand.Uint32()}
}
