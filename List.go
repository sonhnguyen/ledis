package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//
type List struct {
	s     [100]string
	left  int
	right int
}

func NewList() *List {
	return &List{left: 0, right: -1}
}

func (l List) LLen() int {

	return l.right - l.left + 1
}

func (l *List) RPush(values []string) int {

	for _, element := range values {
		l.s[l.right+1] = element
		l.right = l.right + 1

	}

	return l.right - l.left + 1
}

func (l *List) LPop() (string, error) {

	var x string = l.s[l.left]
	l.left = l.left + 1
	return x, nil
}

func (l *List) RPop() (string, error) {
	var x string = l.s[l.right]
	l.right = l.right - 1

	return x, nil
}

func (l *List) LRange(start int, end int) []string {
	x := make([]string, 0)

	for i := start; i <= end; i++ {
		if l.left+i > l.right {
			break
		} else {
			x = append(x, l.s[l.left+i])
		}
	}

	return x
}

func (l List) getData() []string {
	result := l.LRange(l.left, l.right)
	return result
}

func (l List) getType() string {
	return "List"
}

func (l List) saveData() string {
	var result string
	result = result + l.getType() + " " + strconv.Itoa(l.left) + " " + strconv.Itoa(l.right) + " " + strings.Join(l.getData(), " ")
	result = result + "\r\n"
	return result
}

func Task3() {
	reader := bufio.NewReader(os.Stdin)
	l := NewList()
	for {
		line, _ := reader.ReadString('\n')
		words := strings.Fields(line)
		if words[0] == "RPush" {
			data := words[1:]
			fmt.Println(l.RPush(data))
		} else if words[0] == "LLen" {
			fmt.Println(l.LLen())
		} else if words[0] == "RPop" {
			fmt.Println(l.RPop())
		} else if words[0] == "LPop" {
			fmt.Println(l.LPop())
		} else if words[0] == "LRange" {
			start_index, _ := strconv.Atoi(words[1])
			end_index, _ := strconv.Atoi(words[2])
			fmt.Println(l.LRange(start_index, end_index))
		}
	}
}
