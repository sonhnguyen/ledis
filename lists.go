package main

import (
	"fmt"
	"time"
)

// Llen return length of a list
func (store *Store) Llen(key string) (int, error) {
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

	slice, ok := item.Data.([]string)
	if !ok {
		store.lock.RUnlock()
		return 0, fmt.Errorf("the data is not a list")
	}

	store.lock.RUnlock()
	return len(slice), nil
}

// Rpush append 1 or more values to the list, create list if not exists, return length of list after operation
func (store *Store) Rpush(key string, data []string, t time.Duration) (int, error) {
	var e int64

	store.lock.Lock()
	defer store.lock.Unlock()

	if t == DefaultExpiration {
		t = store.defaultExpiration
	}
	if t > 0 {
		e = time.Now().Add(t).UnixNano()
	}

	item, found := store.items[key]
	if !found {
		store.items[key] = Item{
			Data:       []string{},
			Expiration: e,
		}
		item = store.items[key]
		fmt.Println("not found in key, creating new")
	}
	fmt.Println("current:", store.items[key].Data)

	_, ok := item.Data.([]string)
	if !ok {
		return 0, fmt.Errorf("the data is not a list")
	}
	item.Data = append(item.Data, interface{ data })

	return len(item.Data), nil
}
