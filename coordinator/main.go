package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// defining the struct
type ConsitentHash struct {
	ring        map[uint64]string
	sorted_keys map[uint64]struct{}
	replicas    int
}

// get the hash
func (hash ConsitentHash) getHash(data string) uint64 {
	// hashing
	var hashValue = sha256.Sum256([]byte(data))
	// return hex.EncodeToString(hashValue[:])
	num := binary.BigEndian.Uint64(hashValue[:8])
	return num
}

// initializing the class CreateHashRingInitialize
func CreateHashRingInitialize(replicas int) *ConsitentHash {
	return &ConsitentHash{
		ring:        make(map[uint64]string),
		sorted_keys: make(map[uint64]struct{}),
		replicas:    replicas,
	}
}

// functions
func (hash *ConsitentHash) addNodes(val *string) {
	for i := 0; i < hash.replicas; i++ {
		replicaName := fmt.Sprintf("%s-%d", *val, i)

		replicaKey := hash.getHash(replicaName)
		hash.ring[replicaKey] = *val
		// hash.sorted_keys = append(hash.sorted_keys, replicaKey)
		hash.sorted_keys[replicaKey] = struct{}{}	// sets keys duplicate
		fmt.Printf("Replicated Keys: %d\n", replicaKey)
	}
}

// Create the Nodes and all things required for making this work
func main() {
	consistentHashing := CreateHashRingInitialize(3)

	node := "SouravServer"
	consistentHashing.addNodes(&node)
}
