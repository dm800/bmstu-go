package main

import (
	"context"
	"encoding/json"
	"github.com/coder/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

var ws *websocket.Conn
var con *websocket.Conn

type Request struct {
	User string `json:"user"`
	Pass string `json:"password"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Access-Control-Allow-Origin", "*")
	log.Println("served index")
	http.ServeFile(w, r, "index.html")
}

func button(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Access-Control-Allow-Origin", "*")
	log.Println("Button clicked")
	log.Println(r.Method)
	if r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		log.Println("unparsed", string(body))
		var ans Request
		err := json.Unmarshal(body, &ans)
		if err != nil {
			log.Fatal("error", err)
		}
		log.Println("parsed ", ans.Pass, ans.User)
		cred := ans.User + " " + ans.Pass
		ws, _, _ = websocket.Dial(context.Background(), "ws://127.0.0.1:9457/ws", nil)
		ws.Write(context.Background(), websocket.MessageText, []byte(cred))
		time.Sleep(200 * time.Millisecond)
		_, code, _ := ws.Read(context.Background())
		if string(code) == "1" {
			log.Println("Failed")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "Unauthorized"}`))
		} else {
			log.Println("Success")
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"message": "Success"}`))
		}
	} else {
		return
	}
}

func panel(w http.ResponseWriter, r *http.Request) {
	log.Println("Panel")
	http.ServeFile(w, r, "panel.html")
}

func update() {
	for {
		_, text, _ := con.Read(context.Background())
		command := string(text)
		log.Println("send to socket", command)
		ws.Write(context.Background(), websocket.MessageText, []byte(command))
		time.Sleep(200 * time.Millisecond)
		_, text, _ = ws.Read(context.Background())
		log.Println("received from socket", string(text))
		if string(text) == "Done" {
			log.Println("Done")
			continue
		}
		con.Write(context.Background(), websocket.MessageText, text)
	}
}

func web(w http.ResponseWriter, r *http.Request) {
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

func bhelp() {
	http.HandleFunc("/button", button)
	//http.ListenAndServe("127.0.0.1:9455", nil)
}

func main() {
	http.HandleFunc("/panel", panel)
	http.HandleFunc("/ws", web)
	http.HandleFunc("/", handler)
	go bhelp()
	http.ListenAndServe("127.0.0.1:9456", nil)
}
