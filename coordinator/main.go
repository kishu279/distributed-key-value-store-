package main

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"maps"
	"slices"
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
		hash.sorted_keys[replicaKey] = struct{}{} // sets keys duplicate
		fmt.Printf("Replicated Keys: %d\n", replicaKey)
	}
}

func (hash *ConsitentHash) removeNodes(val *string) {
	for i := 0; i < hash.replicas; i++ {
		replicaName := fmt.Sprintf("%s-%d", *val, i)

		replicaKey := hash.getHash(replicaName)

		// remove nodes
		delete(hash.ring, replicaKey)
		delete(hash.sorted_keys, replicaKey)
	}
}

func (hash *ConsitentHash) getNode(key *string) (string, error) {
	if len(hash.ring) == 0 {
		return "", errors.New("Invalid User Id")
	}

	// collect keys and sort the keys
	keys := slices.Collect(maps.Keys(hash.sorted_keys))
	slices.Sort(keys)

	// get the lower bound first keys >= target
	hashValue := hash.getHash(*key)                    // target
	idx, found := slices.BinarySearch(keys, hashValue) // searches for the earliest and also returns idx

	if idx < len(keys) {
		lowerBoundKey := keys[idx]
		fmt.Printf("Exact Match :%v | Lower Bound Key :%d, Value :%s\n", found, lowerBoundKey, hash.sorted_keys[lowerBoundKey])
		return hash.ring[lowerBoundKey], nil
	} else {
		fmt.Printf("No keys >= target found")
		return "", errors.New("Key not found")
	}
}

// Create the Nodes and all things required for making this work
func main() {
	consistentHashing := CreateHashRingInitialize(3)

	nodeA := "SouravServer1"
	nodeB := "SouravServer2"
	nodeC := "SouravServer3"

	consistentHashing.addNodes(&nodeA)
	consistentHashing.addNodes(&nodeB)
	consistentHashing.addNodes(&nodeC)

	newNode := "souravpoddar"
	nodeId, err := consistentHashing.getNode(&newNode)
	if err != nil {
		fmt.Printf("Error %s", err)
	}

	fmt.Printf("Node Id: %s\n", nodeId)

	consistentHashing.printHash()
}

func (hash *ConsitentHash) printHash() {
	// 	for i := 0; i < hash.replicas; i++ {
	// 		replicaName := fmt.Sprintf("%s-%d", *val, i)
	//
	// 		replicaKey := hash.getHash(replicaName)
	// 		fmt.Printf("Ring: %s Sorted Keys: %s\n", hash.ring[replicaKey], hash.sorted_keys[replicaKey])
	// 	}
	// collect keys and sort the keys
	keys := slices.Collect(maps.Keys(hash.sorted_keys))
	slices.Sort(keys)

	for i := 0; i < len(keys); i++ {
		fmt.Printf("%d\n", keys[i])
	}

	// for key, value := range hash.ring {
	// 	fmt.Printf("Values inside the sorted keys Key: %d, Values: %s\n", key, value)
	// }
}
