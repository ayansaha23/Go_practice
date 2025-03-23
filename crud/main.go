package main

//https://jsonplaceholder.typicode.com/guide/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// https://jsonplaceholder.typicode.com/
// practice of simply parsing crud json responses
/*
1.create
2.index
3.get specific post
4.update specific - PATCH and PUT
5.Delete
6. GET /posts/1/comments
*/

/*
	{
	    "postId": 1,
	    "id": 1,
	    "name": "id labore ex et quam laborum",
	    "email": "Eliseo@gardner.biz",
	    "body": "laudantium enim quasi est quidem magnam voluptate ipsam eos\ntempora quo necessitatibus\ndolor quam autem quasi\nreiciendis et nam sapiente accusantium"
	  },
*/
type post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
type Comment struct {
	Id    int    `json: "id"`
	Name  string `json: "name"`
	Email string `json: "email"`
	Body  string `json: "body"`
}

func main() {
	fmt.Println("Crud operations on json")
	//createPost()
	//getPosts()
	//getPost(1)
	updatePost(2)
	//updatePatchPost(1)
	//DeletePost(1)
}

func createPost() {
	// this will be a post request, body of the request will be an
	const url = "https://jsonplaceholder.typicode.com/posts"
	payload := post{
		UserId: 12,
		Id:     24,
		Title:  "Midnight children",
		Body:   "Welcome to the controversial book",
	}
	// need to convert this into bytes []
	encodedPayload, _ := json.Marshal(payload)

	// need to convert bytes[] into io.Reader
	// strings.newReader(s string) does it
	resp, err := http.Post(url, "application/json", strings.NewReader(string(encodedPayload)))
	if err != nil {
		log.Fatal("error in sending POST request:", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func getPost(id int) {
	// read
	fmt.Printf("going to call post for id %v \n", id)
	var url = "https://jsonplaceholder.typicode.com/posts/"
	url = url + strconv.Itoa(id)
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal("could not send request to %v:", url)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("could not read from body", err)
	}
	defer resp.Body.Close()

	var receivedPost post
	err = json.Unmarshal(body, &receivedPost)
	fmt.Println(string(body))
	if err != nil {
		log.Fatal("Could not unmarshal post", err)
	}
	fmt.Println(receivedPost)

}

func updatePatchPost(id int) {
	// here we are trying to update only a single field
	var url = "https://jsonplaceholder.typicode.com/posts/"
	url = url + strconv.Itoa(id)
	client := &http.Client{}
	updateTitle := make(map[string]string)
	updateTitle["Title"] = "haunting"
	payload, err := json.Marshal(updateTitle)
	if err != nil {
		log.Fatal("Could not marshal map into payload", err)
	}
	data := bytes.NewReader(payload)
	req, err := http.NewRequest(http.MethodPatch, url, data)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("error in sending PATCH request:", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("error parsing body of response:", err)
	}
	defer resp.Body.Close()

	fmt.Println(string(body))
}

func updatePost(id int) {
	// update
	var url = "https://jsonplaceholder.typicode.com/posts/"
	url = url + strconv.Itoa(id)
	client := &http.Client{}
	updatedPost := post{10, 2, "New title", "body"}
	data, err := json.Marshal(updatedPost)
	if err != nil {
		log.Fatal("Unable to marshal post: %v and err : %v", updatedPost, err)
	}
	fmt.Println(string(data))
	r := bytes.NewReader(data)
	fmt.Println(url)

	req, _ := http.NewRequest(http.MethodPut, url, r)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error in update of post:", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Error in response parsing:", err)
	}
	var response post
	err = json.Unmarshal(body, &response)
	fmt.Println(response)
}

func getPosts() {
	// index
	const url = "https://jsonplaceholder.typicode.com/posts"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Could not send request to %v", url)
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Could not read body from response:", err)
	}
	defer resp.Body.Close()

	// declare a slice of posts to hold the parsed json
	var posts []post
	err = json.Unmarshal(body, &posts)

	if err != nil {
		log.Fatal("Could not unmarshal posts", err)
	}
	for _, item := range posts {
		fmt.Println(item)
	}
}

func DeletePost(id int) {
	fmt.Println("Going to delete id:", id)
	var url = "https://jsonplaceholder.typicode.com/posts/"
	url = url + strconv.Itoa(id)
	client := &http.Client{}
	payload := ""
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(payload))
	if err != nil {
		log.Fatal("err in creating http request", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Err in sending DEL req", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Err parsing body", err)
	}
	fmt.Println(string(body))
}
