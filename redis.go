package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var keys = make(map[string]Data)

func sinter(sets []string) []string {
	var result []string
	check := make(map[string]int)
	max := 0

	for _, element := range sets {

		if s, ok := keys[element].(*Set); ok {
			members := s.SMembers()
			for _, member := range members {

				check[member]++
				if check[member] > max {
					max = check[member]
				}
			}
		}
	}
	for index, element := range check {
		if element == max {
			result = append(result, index)
		}
	}
	return result
}

func save() {
	//Open file to save
	print("here")
	f, err := os.Create("Snapshot.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	//Save data
	fmt.Println("here2")
	for index, element := range keys {
		fmt.Println(index)
		f.WriteString(index + " " + element.saveData())
	}

	//Close file
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func restore() {
	f, err := os.Open("Snapshot.txt")
	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewScanner(f)
	for r.Scan() {
		data := r.Text()
		fmt.Println(data)
		info := strings.Split(data, " ")
		keyName := info[0]
		switch info[1] {
		case "String":
			keys[keyName] = NewString()
			if str, ok := keys[keyName].(*String); ok {
				str.setData(info[2])
			}
		case "List":
			keys[keyName] = NewList()
			if l, ok := keys[keyName].(*List); ok {
				l.left, err = strconv.Atoi(info[2])
				l.right, err = strconv.Atoi(info[3])
				index := 4
				for i := l.left; i <= l.right; i++ {
					l.s[i] = info[index]
					index++
				}
			}

		case "Set":
			keys[keyName] = NewSet()
			if s, ok := keys[keyName].(*Set); ok {
				//s.count, _ = strconv.Atoi(info[2])
				s.SAdd(info[3:])
			}
		}

	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleFunc(w http.ResponseWriter, r *http.Request) {

	//message := r.URL.Path
	body, err := ioutil.ReadAll(r.Body)
	line := string(body)
	var message string
	command := strings.Split(line, " ")
	fmt.Println(command)
	//Execute request
	if r.Method == "POST" {
		var keyName string
		if command[0] != "SAVE" && command[0] != "RESTORE" {
			keyName = command[1]
		}
		switch command[0] {
		case "GET":
			fmt.Print(keyName)
			if keys[keyName] != nil {
				str, ok := keys[keyName].(*String)

				if ok {

					fmt.Println(str.value)
					message = strings.Join(str.getData(), "")

				}
			} else {
				message = "ERROR: Key does not exist ir is not a String"
			}
		case "SET":
			if keys[keyName] == nil {
				keys[keyName] = NewString()
			}
			message = "NOT SET UP"
			if str, ok := keys[keyName].(*String); ok {
				str.setData(command[2])
				//fmt.Println(str.value)
				//fmt.Println(str.getData())
				//keys[keyName] = str
				message = "SETUP SUCCESSFULLY"
			}

		case "LLEN":
			if l, ok := keys[keyName].(*List); ok {
				message = strconv.Itoa(l.LLen())
				fmt.Println(message)
			} else {
				message = "ERROR: This key datatype is not a list"
			}
		case "RPUSH":
			if keys[keyName] == nil {
				keys[keyName] = NewList()
			}
			if l, ok := keys[keyName].(*List); ok {
				data := command[2:]
				message = strconv.Itoa(l.RPush(data))
			}
		case "LPOP":
			if l, ok := keys[keyName].(*List); ok {
				message, err = l.LPop()
			} else {
				message = "ERROR: This key's datatype is not a list"
			}
		case "RPOP":
			if l, ok := keys[keyName].(*List); ok {
				message, err = l.RPop()
			} else {
				message = "ERROR: This key's datatype is not a list"
			}
		case "LRANGE":
			if l, ok := keys[keyName].(*List); ok {
				startIndex, _ := strconv.Atoi(command[2])
				endIndex, _ := strconv.Atoi(command[3])
				message = strings.Join(l.LRange(startIndex, endIndex), " ")
			} else {
				message = "ERROR: This key's datatype is not a list"
			}
		case "SADD":
			if keys[keyName] == nil {
				keys[keyName] = NewSet()
			}
			if s, ok := keys[keyName].(*Set); ok {
				data := command[2:]
				message = strconv.Itoa(s.SAdd(data))
			} else {
				message = "ERROR: This key's datatype is not a set"
			}

		case "SCARD":
			if s, ok := keys[keyName].(*Set); ok {
				message = strconv.Itoa(s.SCard())
			} else {
				message = "ERROR: This key's datatype is not a set"
			}
		case "SMEMBERS":
			if s, ok := keys[keyName].(*Set); ok {
				message = strings.Join(s.SMembers(), " ")
			} else {
				message = "ERROR: THis key's datatype is not a set"
			}
		case "SREM":
			if s, ok := keys[keyName].(*Set); ok {
				data := command[2:]
				message = strconv.Itoa(s.SRem(data))
			} else {
				message = "ERROR: This key's datatype is not a set"
			}
		case "SINTER":
			data := command[1:]
			message = strings.Join(sinter(data), " ")
		case "SAVE":
			save()
		case "RESTORE":
			restore()
		}

	}

	//Read body

	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s", body)

	message = strings.TrimPrefix(message, "/")

	w.Write([]byte(message))
}

func main() {
	//keys := make(map[string]Data, 100)
	http.HandleFunc("/", handleFunc)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
