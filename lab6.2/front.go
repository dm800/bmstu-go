package main

import (
	"context"
	"encoding/json"
	"github.com/coder/websocket"
	"io"
	"log"
	"net/http"
)

var ws *websocket.Conn

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
	} else {
		return
	}
}

func panel(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "panel.html")
}

func web(w http.ResponseWriter, r *http.Request) {
	log.Println("Win")
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
