package main

import "strings"

// String data structure
type String struct {
	value   string
	keyName string
}

func (s String) getData() []string {
	var result []string
	result = append(result, s.value)
	return result
}

func (s *String) setData(data string) {
	s.value = data
}

//NewString to init new String data structure
func NewString() *String {
	return &String{value: " "}
}

func (s String) getType() string {
	return "String"
}

func (s String) saveData() string {
	var result string
	result = result + s.getType()
	result = result + " " + strings.Join(s.getData(), "")
	result = result + "\r\n"
	return result
}
