package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
		if len(agrs) != 1 {
			return "", fmt.Errorf("GET expects 1 argument")
		}
		result = store.Get(agrs[0])

	case "SET":
		if len(agrs) != 2 {
			return "", fmt.Errorf("SET expects 2 arguments")
		}
		store.Set(agrs[0], agrs[1], NoExpiration)
		result = "OK"
	case "LLEN":
		if len(agrs) != 1 {
			return "", fmt.Errorf("LLEN expects 1 arguments")
		}
		result, err = store.LLen(agrs[0])
		if err != nil {
			return "", fmt.Errorf("error in LLEN: %s", err)
		}
	case "RPUSH":
		if len(agrs) < 2 {
			return "", fmt.Errorf("RPUSH expects at least 2 arguments")
		}
		result, err = store.RPush(agrs[0], agrs[1:], NoExpiration)
		if err != nil {
			return "", fmt.Errorf("error in rpush: %s", err)
		}
	case "LPOP":
		if len(agrs) != 1 {
			return "", fmt.Errorf("LPOP expects 1 arguments")
		}
		result, err = store.LPop(agrs[0])
		if err != nil {
			return "", fmt.Errorf("error in lpop: %s", err)
		}
	case "RPOP":
		if len(agrs) != 1 {
			return "", fmt.Errorf("RPOP expects 1 arguments")
		}
		result, err = store.RPop(agrs[0])
		if err != nil {
			return "", fmt.Errorf("error in rpop: %s", err)
		}
	case "LRANGE":
		if len(agrs) != 3 {
			return "", fmt.Errorf("RPOP expects 3 arguments")
		}
		start, err := strconv.Atoi(agrs[1])
		if err != nil {
			return "", fmt.Errorf("unable to convert start to number")
		}
		end, err := strconv.Atoi(agrs[2])
		if err != nil {
			return "", fmt.Errorf("unable to convert end to number")
		}
		result, err = store.LRange(agrs[0], start, end)
		if err != nil {
			return "", fmt.Errorf("error in rpop: %s", err)
		}
	default:
		return "", fmt.Errorf("unknown command")
	}

	return result, nil
}

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

// LLen return length of a list
func (store *Store) LLen(key string) (int, error) {
	store.lock.RLock()
	// "Inlining" of get and Expired
	item, found := store.items[key]
	if !found {
		store.lock.RUnlock()
		return 0, nil
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			store.lock.RUnlock()
			return 0, nil
		}
	}

	slice, ok := item.Data.(*List)
	if !ok {
		store.lock.RUnlock()
		return 0, fmt.Errorf("the data is not a list")
	}

	store.lock.RUnlock()
	return slice.LLen(), nil
}

// RPush append 1 or more values to the list, create list if not exists, return length of list after operation
func (store *Store) RPush(key string, data []string, t time.Duration) (int, error) {

	_, found := store.items[key]
	if !found {
		store.Set(key, &List{GoList: list.New()}, t)
	}

	list, ok := store.items[key].Data.(*List)
	if !ok {
		return 0, fmt.Errorf("the data is not a list")
	}

	store.lock.Lock()
	defer store.lock.Unlock()

	return list.RPush(data), nil
}

// LPop remove and return the first item of the list
func (store *Store) LPop(key string) (string, error) {
	store.lock.Lock()
	defer store.lock.Unlock()

	list, ok := store.items[key].Data.(*List)
	if !ok {
		return "", fmt.Errorf("the data is not a list")
	}

	return list.LPop()
}

// RPop remove and return the last item of the list
func (store *Store) RPop(key string) (string, error) {
	store.lock.Lock()
	defer store.lock.Unlock()

	list, ok := store.items[key].Data.(*List)
	if !ok {
		return "", fmt.Errorf("the data is not a list")
	}

	return list.RPop()
}

// LRange return a range of element from the list (zero-based, inclusive of start and stop), start and stop are non-negative integers
func (store *Store) LRange(key string, start, end int) ([]string, error) {
	store.lock.Lock()
	defer store.lock.Unlock()

	list, ok := store.items[key].Data.(*List)
	if !ok {
		return []string{}, fmt.Errorf("the data is not a list")
	}

	return list.LRange(start, end), nil
}
