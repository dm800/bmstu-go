package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/coder/websocket"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strings"
	"time"
)

var db *sql.DB
var c int
var conn *websocket.Conn

func update() []string {
	ans, err := db.Query("SELECT * FROM `iu9Limonov`")
	defer ans.Close()
	if err != nil {
		log.Fatal(err)
	}
	var feed []string
	feed = make([]string, 0)
	for ans.Next() {
		var val string
		var id int
		var date string
		err := ans.Scan(&id, &val, &date)
		if err != nil {
			fmt.Println("Error occured")
			return nil
		}
		feed = append(feed, val)
	}
	return feed
}

func msg(val string) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	wr, err := conn.Writer(ctx, websocket.MessageType(websocket.MessageText))
	if err != nil {
		log.Println("sending to websocket err:", err)
		return
	}
	log.Println("trying to write msg")
	wr.Write([]byte(val))
	wr.Close()
	log.Println("i believe that we sent msg")
}

func maincycle() {
	for {
		log.Println("Serving a client")
		k := update()
		log.Println("Updated")
		if c != len(k) {
			c = len(k)
			if conn != nil {
				log.Println("Sending")
				msg(strings.Join(k, "<br>"))
				log.Println("Sent")
			}
		}
		time.Sleep(time.Second * 10)
	}
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	go maincycle()
	log.Println("went into handler")
	var err error
	conn, err = websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	msg("Connected")
	if err != nil {
		fmt.Println(err)
	}
}

func RouterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("went into HTTP handler")
	http.ServeFile(w, r, "websocket_client.html")
}

func runsocket() {
	http.HandleFunc("/socket", HomeRouterHandler) // установим роутер
	log.Println("Starting listening ws at", "127.0.0.1:9457")
	err := http.ListenAndServe("127.0.0.1:9457", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	c = 0
	db, _ = sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su:3306)/iu9networkslabs")
	fmt.Println(db == nil)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	http.HandleFunc("/", RouterHandler)
	go http.ListenAndServe("176.124.206.238:9456", nil)
	runsocket()
}
