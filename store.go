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
	var err error
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
	case "LLEN":
		if len(agrs) > 1 {
			return "", fmt.Errorf("LLEN expects 1 arguments")
		}
		result, err = store.Llen(agrs[0])
		if err != nil {
			return "", fmt.Errorf("error in LLEN: %s", err)
		}
	case "RPUSH":
		if len(agrs) < 2 {
			return "", fmt.Errorf("RPUSH expects at least 2 arguments")
		}
		result, err = store.Rpush(agrs[0], agrs[1:], NoExpiration)
		if err != nil {
			return "", fmt.Errorf("error in rpush: %s", err)
		}
	default:
		return "", fmt.Errorf("unknown command")
	}

	return result, nil
}
