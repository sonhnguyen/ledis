package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

// mergeSlices merge 2 slices
func mergeSlices(s1, s2 []string) []string {
	hash := make(map[string]bool)
	results := []string{}
	for _, s := range s1 {
		hash[s] = true
	}

	for _, s := range s2 {
		if hash[s] {
			results = append(results, s)
		}
	}
	return results
}

// SaveFile save an interface to file.
func SaveFile(path string, v interface{}) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)

	_, err = io.Copy(f, r)
	return err
}

// RestoreFile restore from a file to a interface
func RestoreFile(path string, v interface{}) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(v)

	return err
}
