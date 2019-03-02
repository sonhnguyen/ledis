package main

import "time"

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
