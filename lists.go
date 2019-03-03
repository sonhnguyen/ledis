package main

import (
	"container/list"
	"fmt"
)

// List type for data
type List struct {
	GoList *list.List
}

// LLen return length of a list.
func (l *List) LLen() int {
	return l.GoList.Len()
}

// RPush append 1 or more values to the list, create list if not exists, return length of list after operation.
func (l *List) RPush(values []string) int {
	for _, v := range values {
		l.GoList.PushBack(v)
	}
	return l.LLen()
}

// LPop remove and return the first item of the list.
func (l *List) LPop() (string, error) {
	e := l.GoList.Front()
	if e == nil {
		return "", fmt.Errorf("no data in the list")
	}
	l.GoList.Remove(e)
	return e.Value.(string), nil
}

// RPop remove and return the last item of the list.
func (l *List) RPop() (string, error) {
	e := l.GoList.Back()
	if e == nil {
		return "", fmt.Errorf("no data in the list")
	}
	l.GoList.Remove(e)
	return e.Value.(string), nil
}

// LRange return a range of element from the list (zero-based, inclusive of start and stop), start and stop are non-negative integers
func (l *List) LRange(start int, end int) []string {
	values := []string{}
	if start > end {
		return values
	}

	e := l.GoList.Front()

	for i := 0; i <= end && e != nil; i++ {
		if i >= start {
			values = append(values, e.Value.(string))
		}
		e = e.Next()
	}
	return values
}
