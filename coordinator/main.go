package main

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	// "log"
	// "net/http"
	// "net/http/httputil"
	// "net/url"
)

// defining the struct
type ConsitentHash struct {
	ring        map[uint64]string
	sorted_keys map[uint64]struct{}
	servers     map[string]string
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
		servers:     make(map[string]string),
		replicas:    replicas,
	}
}

func (hash ConsitentHash) getIpFromNodeId(nodeId *string) string {
	return hash.servers[*nodeId]
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
	hashValue := hash.getHash(*key)                // target
	idx, _ := slices.BinarySearch(keys, hashValue) // searches for the earliest and also returns idx

	if idx == len(keys) {
		idx = 0
	}

	loweBoundKey := keys[idx]
	return hash.ring[loweBoundKey], nil
}

// Create the Nodes and all things required for making this work
func (hash *ConsitentHash) ConfiguringHashingServers(servers map[string]string) {

	// put all the servers ip with node id name and calculate the hash ring
	for key, value := range servers {
		hash.addNodes(&key)
		hash.servers[key] = value
	}

	fmt.Printf("Server Configured\n")
}

// calculate the each requst hash and get node id
func (hash ConsitentHash) CalculateRoute(requestId *string) string {
	nodeId, err := hash.getNode(requestId)
	if err != nil {
		fmt.Printf("Error %s\n", err)
	}

	return hash.getIpFromNodeId(&nodeId)
}

func main() {
	consistentHashing := CreateHashRingInitialize(3)

	// configuring the consistent hashing
	servers := map[string]string{
		"Node1": "http://localhost:8080",
		"Node2": "http://localhost:8081",
		"Node3": "http://localhost:8082",
	}

	consistentHashing.ConfiguringHashingServers(servers) // configuration done

	// incomind request
	requestId := "PoddarKishu279"

	target := consistentHashing.CalculateRoute(&requestId) // target ip address
	fmt.Printf("Target IP Address: %s\n", target)

	// reverse proxy -> direct the request to target
	targetIpAddress, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetIpAddress)
	log.Fatal(http.ListenAndServe(":3000", proxy))
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
