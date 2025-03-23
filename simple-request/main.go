package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type user struct {
	UserId    int8   `json: userId`
	Id        int8   `json: id`
	Title     string `json: title`
	Completed bool   `json: completed`
}

func main() {
	// here make simple http get request and parse json data
	const URL = "https://jsonplaceholder.typicode.com/todos/1"
	resp, err := http.Get(URL)

	if err != nil {
		log.Fatal("Could not send request:", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Could not read response body:", err)
	}

	defer resp.Body.Close()

	fmt.Println(string(body))
	getTodo(body)
}

func getTodo(data []byte) {
	var receivedUser user
	// for parsing json , unmarshall into a struct
	err := json.Unmarshal(data, &receivedUser)
	if err != nil {
		fmt.Println("unable to parse json response")
	}
	fmt.Println(receivedUser)
}
