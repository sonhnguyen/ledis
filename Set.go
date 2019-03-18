package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Set struct {
	set   map[string]bool
	count int
}

func NewSet() *Set {
	return &Set{set: make(map[string]bool), count: 0}
}
func (s *Set) SAdd(values []string) int {
	add_count := 0

	for _, element := range values {
		_, found := s.set[element]
		if !found {
			s.set[element] = true
			add_count++
		}
	}
	s.count = s.count + add_count
	return add_count
}

func (s Set) SCard() int {
	return s.count
}

func (s Set) SMembers() []string {
	var a []string
	for k, v := range s.set {
		if v == true {
			a = append(a, k)
		}
	}
	return a
}

func (s *Set) SRem(values []string) int {
	remove_count := 0
	for _, element := range values {
		_, found := s.set[element]
		if found {
			s.set[element] = false
			remove_count++
		}
	}
	s.count = s.count - remove_count
	return remove_count
}

func (s Set) getData() []string {
	result := s.SMembers()
	return result
}

func (s Set) getType() string {
	return "Set"
}
func (s Set) saveData() string {
	var result string
	result = result + s.getType() + " " + strconv.Itoa(s.count) + " " + strings.Join(s.getData(), " ") + "\r\n"
	return result
}
func Task4() {
	reader := bufio.NewReader(os.Stdin)
	s := NewSet()
	for {
		line, _ := reader.ReadString('\n')
		words := strings.Fields(line)
		if words[0] == "SAdd" {
			data := words[1:]
			fmt.Println(s.SAdd(data))
		} else if words[0] == "SCard" {
			fmt.Println(s.SCard())
		} else if words[0] == "SMembers" {
			fmt.Println(s.SMembers())
		} else if words[0] == "SRem" {
			data := words[1:]
			fmt.Println(s.SRem(data))
		}
	}
}
