package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// defining the struct
type ConsitentHash struct {
	ring        map[int]string
	sorted_keys map[string]struct{}
	replicas    int
}

// get the hash
func (hash ConsitentHash) getHash(data *string) string {
	// hashing
	var hashValue = sha256.Sum256([]byte(*data))
	return hex.EncodeToString(hashValue[:])
}

// initializing the class CreateHashRingInitialize
func CreateHashRingInitialize(replicas int) *ConsitentHash {
	return &ConsitentHash{
		ring:        make(map[int]string),
		sorted_keys: make(map[string]struct{}),
		replicas:    replicas,
	}
}

// functions
func (hash *ConsitentHash) addNodes(val *string) {
	// for i := 0; i < hash.replicas; i++ {
	// }
	var replica_keys = hash.getHash(val)
	fmt.Printf("replica keys creted %s\n", replica_keys)
}

// Create the Nodes and all things required for making this work
func main() {
	consistentHashing := ConsitentHash{replicas: 3}

	node := "SouravServer"
	consistentHashing.addNodes(&node)
}
