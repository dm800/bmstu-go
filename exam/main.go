package main

import (
	"context"
	"encoding/json"
	"github.com/coder/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

var con *websocket.Conn

type Request struct {
	A string `json:"a"`
	B string `json:"b"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Websocket")
	var err error
	con, err = websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println(err)
	}
	go update()
}

func logic(a int, b int) {}

func update() {
	for {
		log.Println("updating")
		_, text, _ := con.Read(context.Background())
		var data Request
		err := json.Unmarshal(text, &data)
		if err != nil {
			log.Println("Error", err)
			continue
		}
		log.Println(text)
		time.Sleep(200 * time.Millisecond)
		log.Println("received from socket", data.A, data.B)
		if string(text) == "Done" {
			log.Println("Done")
			continue
		}
		A, _ := strconv.Atoi(data.A)
		B, _ := strconv.Atoi(data.B)
		con.Write(context.Background(), websocket.MessageText, []byte(strconv.Itoa(A+B)))
	}
}

func main() {
	http.HandleFunc("/socket", handler)
	http.ListenAndServe("127.0.0.1:9457", nil)
}
