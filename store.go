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
	Expiration *int64
}

// Store represents in-memory db store
type Store struct {
	DefaultExpiration time.Duration
	Items             map[string]Item
	Lock              sync.RWMutex
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
	case "SADD":
		if len(agrs) < 2 {
			return "", fmt.Errorf("SADD expects at least 2 arguments")
		}
		result, err = store.SAdd(agrs[0], agrs[1:], NoExpiration)
		if err != nil {
			return "", fmt.Errorf("error in SADD: %s", err)
		}
	case "SCARD":
		if len(agrs) != 1 {
			return "", fmt.Errorf("SCARD expects 1 arguments")
		}
		result, err = store.SCard(agrs[0])
		if err != nil {
			return "", fmt.Errorf("error in SCARD: %s", err)
		}
	case "SMEMBERS":
		if len(agrs) != 1 {
			return "", fmt.Errorf("SMEMBERS expects 1 arguments")
		}
		result, err = store.SMembers(agrs[0])
		if err != nil {
			return "", fmt.Errorf("error in SMEMBERS: %s", err)
		}
	case "SREM":
		if len(agrs) < 2 {
			return "", fmt.Errorf("SREM expects at least 2 arguments")
		}
		result, err = store.SRem(agrs[0], agrs[1:])
		if err != nil {
			return "", fmt.Errorf("error in SREM: %s", err)
		}
	case "SINTER":
		if len(agrs) < 1 {
			return "", fmt.Errorf("SINTER expects at least 1 arguments")
		}
		result, err = store.SInter(agrs[0:])
		if err != nil {
			return "", fmt.Errorf("error in SINTER: %s", err)
		}
	case "KEYS":
		if len(agrs) != 0 {
			return "", fmt.Errorf("KEYS expects 0 arguments")
		}
		result, err = store.Keys()
		if err != nil {
			return "", fmt.Errorf("error in KEYS: %s", err)
		}
	case "DEL":
		if len(agrs) != 1 {
			return "", fmt.Errorf("DEL expects 1 arguments")
		}
		err = store.Del(agrs[0])
		if err != nil {
			return "", fmt.Errorf("error in DEL: %s", err)
		}
		result = "OK"
	case "FLUSHDB":
		if len(agrs) != 0 {
			return "", fmt.Errorf("FLUSHDB expects 0 arguments")
		}
		err = store.FlushDB()
		if err != nil {
			return "", fmt.Errorf("error in FLUSHDB: %s", err)
		}
		result = "OK"
	case "EXPIRE":
		if len(agrs) != 2 {
			return "", fmt.Errorf("EXPIRE expects 2 arguments")
		}
		seconds, err := strconv.Atoi(agrs[1])
		if err != nil {
			return "", fmt.Errorf("unable to convert seconds to number")
		}

		result, err = store.Expire(agrs[0], time.Duration(seconds)*time.Second)
		if err != nil {
			return "", fmt.Errorf("error in EXPIRE: %s", err)
		}
	case "TTL":
		if len(agrs) != 1 {
			return "", fmt.Errorf("TTL expects 1 arguments")
		}
		result, err = store.TTL(agrs[0])
		if err != nil {
			return "", fmt.Errorf("error in TTL: %s", err)
		}
	case "SAVE":
		if len(agrs) != 0 {
			return "", fmt.Errorf("SAVE expects 0 arguments")
		}
		err = store.Save("snapshot.db")
		if err != nil {
			return "", fmt.Errorf("error in SAVE: %s", err)
		}
	case "RESTORE":
		if len(agrs) != 0 {
			return "", fmt.Errorf("SAVE expects 0 arguments")
		}
		err = store.Restore("snapshot.db")
		if err != nil {
			return "", fmt.Errorf("error in SAVE: %s", err)
		}
	default:
		return "", fmt.Errorf("unknown command")
	}

	return result, nil
}

// Get at key
func (store *Store) Get(key string) interface{} {
	store.Lock.RLock()
	defer store.Lock.RUnlock()

	item, found := store.Items[key]
	if !found {
		return nil
	}

	if *item.Expiration > 0 {
		if time.Now().UnixNano() > *item.Expiration {
			return nil
		}
	}

	return item.Data
}

// Set at key
func (store *Store) Set(key string, data interface{}, t time.Duration) {
	var e int64

	store.Lock.Lock()
	defer store.Lock.Unlock()

	if t == DefaultExpiration {
		t = store.DefaultExpiration
	}
	if t > 0 {
		e = time.Now().Add(t).UnixNano()
	}

	store.Items[key] = Item{
		Data:       data,
		Expiration: &e,
	}
}

// LLen return length of a list
func (store *Store) LLen(key string) (int, error) {
	store.Lock.RLock()
	defer store.Lock.RUnlock()

	slice, ok := store.Get(key).(*List)
	if !ok {
		return 0, fmt.Errorf("the data is not a list")
	}

	return slice.LLen(), nil
}

// RPush append 1 or more values to the list, create list if not exists, return length of list after operation
func (store *Store) RPush(key string, data []string, t time.Duration) (int, error) {
	_, found := store.Items[key]
	if !found {
		store.Set(key, &List{GoList: list.List{}}, t)
	}

	store.Lock.RLock()
	list, ok := store.Get(key).(*List)
	if !ok {
		return 0, fmt.Errorf("the data is not a list")
	}
	store.Lock.RUnlock()

	store.Lock.Lock()
	defer store.Lock.Unlock()

	return list.RPush(data), nil
}

// LPop remove and return the first item of the list
func (store *Store) LPop(key string) (string, error) {
	store.Lock.RLock()
	list, ok := store.Get(key).(*List)
	if !ok {
		return "", fmt.Errorf("the data is not a list")
	}
	store.Lock.RUnlock()

	store.Lock.Lock()
	defer store.Lock.Unlock()

	result, err := list.LPop()
	return result, err
}

// RPop remove and return the last item of the list
func (store *Store) RPop(key string) (string, error) {
	store.Lock.RLock()
	list, ok := store.Get(key).(*List)
	if !ok {
		return "", fmt.Errorf("the data is not a list")
	}
	store.Lock.RUnlock()

	store.Lock.Lock()
	defer store.Lock.Unlock()

	result, err := list.RPop()

	return result, err
}

// LRange return a range of element from the list (zero-based, inclusive of start and stop), start and stop are non-negative integers
func (store *Store) LRange(key string, start, end int) ([]string, error) {
	store.Lock.RLock()
	defer store.Lock.RUnlock()

	list, ok := store.Get(key).(*List)
	if !ok {
		return []string{}, fmt.Errorf("the data is not a list")
	}

	return list.LRange(start, end), nil
}

// SAdd add values to set stored at key
func (store *Store) SAdd(key string, values []string, t time.Duration) (int, error) {
	_, found := store.Items[key]
	if !found {
		store.Set(key, Set{Set: make(map[string]bool)}, t)
	}

	store.Lock.RLock()
	set, ok := store.Get(key).(Set)
	if !ok {
		return 0, fmt.Errorf("the data is not a set")
	}
	store.Lock.RUnlock()

	store.Lock.Lock()
	defer store.Lock.Unlock()

	return set.SAdd(values), nil
}

// SCard return the number of elements of the set stored at key
func (store *Store) SCard(key string) (int, error) {
	store.Lock.RLock()
	defer store.Lock.RUnlock()

	set, ok := store.Get(key).(Set)
	if !ok {
		return 0, fmt.Errorf("the data is not a set")
	}

	return set.SCard(), nil
}

// SMembers return array of all members of set
func (store *Store) SMembers(key string) ([]string, error) {
	store.Lock.RLock()
	defer store.Lock.RUnlock()

	set, ok := store.Get(key).(Set)
	if !ok {
		return []string{}, fmt.Errorf("the data is not a set")
	}

	return set.SMembers(), nil
}

// SRem remove values from set
func (store *Store) SRem(key string, values []string) (int, error) {
	store.Lock.RLock()
	set, ok := store.Get(key).(Set)
	if !ok {
		return 0, fmt.Errorf("the data is not a set")
	}
	store.Lock.RUnlock()

	store.Lock.Lock()
	defer store.Lock.Unlock()

	return set.SRem(values), nil
}

// SInter set intersection among all set stored in specified keys. Return array of members of the result set
func (store *Store) SInter(keys []string) ([]string, error) {
	results := []string{}

	for _, k := range keys {
		store.Lock.RLock()
		set, ok := store.Get(k).(Set)
		if !ok {
			return []string{}, fmt.Errorf("the data is not a set at key: %s", k)
		}
		store.Lock.RUnlock()

		if len(results) == 0 {
			results = set.SMembers()
		}
		results = mergeSlices(results, set.SMembers())
	}

	return results, nil
}

// Keys List all available keys
func (store *Store) Keys() ([]string, error) {
	results := []string{}

	for k := range store.Items {
		store.Lock.RLock()
		if store.Get(k) != nil {
			results = append(results, k)
		}
		store.Lock.RUnlock()
	}

	return results, nil
}

// Del Delete key
func (store *Store) Del(key string) error {
	store.Lock.Lock()
	defer store.Lock.Unlock()

	delete(store.Items, key)

	return nil
}

// FlushDB Delete all keys
func (store *Store) FlushDB() error {
	store.Lock.Lock()
	defer store.Lock.Unlock()

	store.Items = make(map[string]Item)

	return nil
}

// Expire set a timeout on a key, seconds is a positive integer. Return the number of seconds if the timeout is set
func (store *Store) Expire(key string, t time.Duration) (float64, error) {
	var e int64
	store.Lock.Lock()
	defer store.Lock.Unlock()

	item, found := store.Items[key]
	if !found {
		return 0, fmt.Errorf("this key not existed")
	}

	if t == DefaultExpiration {
		t = store.DefaultExpiration
	}
	if t > 0 {
		e = time.Now().Add(t).UnixNano()
	}

	*item.Expiration = e
	return t.Seconds(), nil
}

// TTL key: query the timeout of a key
func (store *Store) TTL(key string) (float64, error) {
	store.Lock.RLock()
	defer store.Lock.RUnlock()
	timeNow := time.Now()

	item, found := store.Items[key]
	if !found || (*item.Expiration > 0 && *item.Expiration < timeNow.UnixNano()) {
		return -2, nil
	}
	if *item.Expiration == 0 {
		return -1, nil
	}

	return time.Unix(0, *item.Expiration).Sub(timeNow).Seconds(), nil
}

// Save saves a snapshot.
func (store *Store) Save(path string) error {
	store.Lock.Lock()
	defer store.Lock.Unlock()

	err := SaveFile(path, store)
	if err != nil {
		return err
	}

	return nil
}

// Restore restore from the last snapshot,
func (store *Store) Restore(path string) error {
	store.Lock.Lock()
	defer store.Lock.Unlock()

	if err := RestoreFile(path, store); err != nil {
		return err
	}
	return nil
}
