package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Item represents our values
type Item struct {
	Data       interface{}
	Expiration int64
}

// Store represents in-memory db store
type Store struct {
	defaultExpiration time.Duration
	items             map[string]Item
	lock              sync.RWMutex
}

const (
	// NoExpiration is for use to init items without expiration, if it is -1, the item never expired.
	NoExpiration time.Duration = -1
	// DefaultExpiration is for use to init items with expiration, if using 0, it used the configed expiration time.
	DefaultExpiration time.Duration = 0
)

// Get at key
func (store *Store) Get(key string) interface{} {
	store.lock.RLock()
	// "Inlining" of get and Expired
	item, found := store.items[key]
	if !found {
		store.lock.RUnlock()
		return nil
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			store.lock.RUnlock()
			return nil
		}
	}
	store.lock.RUnlock()
	return item.Data
}

// Set at key
func (store *Store) Set(key string, data interface{}, t time.Duration) {
	var e int64

	store.lock.Lock()
	defer store.lock.Unlock()

	if t == DefaultExpiration {
		t = store.defaultExpiration
	}
	if t > 0 {
		e = time.Now().Add(t).UnixNano()
	}

	store.items[key] = Item{
		Data:       data,
		Expiration: e,
	}
}

func (store *Store) requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		command := strings.Split(string(body), " ")
		result, err := store.commandHandler(command[0], command[1:])
		if err != nil {
			log.Printf("error when processing command: %s", err)
			http.Error(w, err.Error(), 400)
		} else {
			err = json.NewEncoder(w).Encode(result)
			if err != nil {
				log.Printf("error when return json: %s", err)
				http.Error(w, err.Error(), 500)
			}
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func (store *Store) commandHandler(name string, agrs []string) (interface{}, error) {
	var result interface{}

	switch name {
	case "GET":
		if len(agrs) < 1 {
			return "", fmt.Errorf("GET expects 1 argument")
		}
		result = store.Get(agrs[0])

	case "SET":
		if len(agrs) < 2 {
			return "", fmt.Errorf("SET expects 2 arguments")
		}
		store.Set(agrs[0], agrs[1], NoExpiration)
		result = "OK"
	default:
		return "", fmt.Errorf("unknown command")
	}

	return result, nil
}

func main() {
	store := Store{defaultExpiration: 100, items: make(map[string]Item), lock: sync.RWMutex{}}
	http.HandleFunc("/", store.requestHandler)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
