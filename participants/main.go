package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shirou/gopsutil/v4/mem"
)

var Cache = make(map[string]int)

func SetCache(key string, value int) {
	Cache[key] = value
}

func GetCache(key string) int {
	return Cache[key]
}

func CheckKeyPresent(key string) bool {
	if _, ok := Cache[key]; ok {
		return true
	} else {
		return false
	}
}

/**
 * User Request Handlers
 */

type UserRequest struct {
	Email    string `json:"email"`
	Password int    `json:"password"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Go Web Server!")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // returning status
	v, _ := mem.VirtualMemory()
	fmt.Fprintf(w, "Total: %v, Available: %v, UsedPercent:%f%%\n", v.Total, v.Available, v.UsedPercent)
	fmt.Println(v)
}

func createNewKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Always close the request body when done
	defer r.Body.Close()

	var user UserRequest

	// 2. Decode the JSON stream directly into the struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// 3. Use the parsed data
	// fmt.Fprintf(w, "Recieved Email %s and Password %d\n", user.Email, user.Password)

	// log the request
	// fmt.Println(user.Email)
	// fmt.Println(user.Password)

	// check user key is alread present
	if CheckKeyPresent(user.Email) {
		fmt.Fprintln(w, "Key is already present")
		return
	}

	// set key
	SetCache(user.Email, user.Password)

	fmt.Fprintln(w, "Key is Stored Sucessfully")
}

func updateUserKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Always close the request body when done
	defer r.Body.Close()

	var user UserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// check user key is alread present
	if !CheckKeyPresent(user.Email) {
		fmt.Fprintln(w, "Key is not present first create")
		return
	}

	SetCache(user.Email, user.Password)

	fmt.Fprintln(w, "User keys updated")
}

func getUserPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Always close the request body when done
	defer r.Body.Close()

	var user UserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// check user key is alread present
	if !CheckKeyPresent(user.Email) {
		fmt.Fprintln(w, "Key is not present first create")
		return
	}

	var password int = GetCache(user.Email)
	fmt.Fprintf(w, "User Password is :%d", password)
}

func main() {
	// 1. Create a new request router (ServeMux)
	mux := http.NewServeMux()

	// 2. Register the handler function for the root URL path ("/")
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/home", homeHandler)
	mux.HandleFunc("/checkhealth", checkHealth)
	mux.HandleFunc("/createUserKey", createNewKey)
	mux.HandleFunc("/updateUserKey", updateUserKey)
	mux.HandleFunc("/getUserKey", getUserPassword)

	// 3. Start the server on port 8080
	fmt.Println("Server starting on http://localhost:8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
