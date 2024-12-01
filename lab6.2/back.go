package main

import (
	"context"
	"fmt"
	"github.com/coder/websocket"
	"net/http"
	"strings"
	"time"
)
import (
	"github.com/jlaffaye/ftp"
	"log"
)

var c *websocket.Conn
var conn *ftp.ServerConn

func handle(text string, cur string) string {
	fmt.Println("executing", text)
	fmt.Println(conn == nil)
	splitted := strings.Split(text, " ")
	req := splitted[0]
	path := "/"
	if req != "END" {
		/*if req == "GET" {
			log.Println("Doing get")
			fmt.Println("Path to retrieve")
			fmt.Scan(&path)
			r, err := conn.Retr(cur + path)
			if err != nil {
				log.Println(err)
				continue
			}
			buf, _ := ioutil.ReadAll(r)
			os.WriteFile(path, buf, 0666)
			log.Println("Done")
		} else if req == "PUSH" {
			log.Println("Doing push")
			fmt.Println("Push where")
			fmt.Scan(&path)
			fmt.Println("push from")
			fmt.Scan(&path2)
			reader, err := os.ReadFile("./" + path2)
			if err != nil {
				log.Println(err)
				continue
			}
			file := bytes.NewBuffer(reader)
			err = conn.Stor(cur+path, file)
			if err != nil {
				log.Println(err)
			}
			log.Println("Done") */
		if req == "MKDIR" {
			path = splitted[1]
			err := conn.MakeDir(cur + path)
			if err != nil {
				log.Println(err)
			}
			log.Println("Done")
		} else if req == "RMF" {
			path = splitted[1]
			err := conn.Delete(path)
			if err != nil {
				log.Println(err)
			}
			log.Println("Done")
		} else if req == "LS" {
			entry, err := conn.List(cur)
			st := ""
			if err != nil {
				log.Println(err)
			}
			for _, elem := range entry {
				st += elem.Name + string(rune(elem.Type)) + elem.Time.String() + "<br>"
			}
			c.Write(context.Background(), websocket.MessageText, []byte(st))
			return cur
		} else if req == "CD" {
			path = splitted[1]
			err := conn.ChangeDir(cur + path)
			if err != nil {
				log.Println(err)
				return cur
			}
			cur = cur + path + "/"
			log.Println("Done", cur)
		} else if req == "RMD" {
			path = splitted[1]
			err := conn.RemoveDir(cur + path)
			cur = cur + path + "/"
			if err != nil {
				log.Println(err)
			}
			log.Println("Done")
		} else if req == "RMDR" {
			path = splitted[1]
			err := conn.RemoveDirRecur(cur + path)
			cur = cur + path + "/"
			if err != nil {
				log.Println(err)
			}
			log.Println("Done")
		}
		c.Write(context.Background(), websocket.MessageText, []byte("Done"))
	}
	return cur
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("went into handler")
	var err error
	c, err = websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println(err)
	}
	_, text, _ := c.Read(context.Background())
	splitted := strings.Split(string(text), " ")
	log.Println(splitted[0], splitted[1])
	result := login(splitted[0], splitted[1])
	c.Write(context.Background(), websocket.MessageText, []byte(fmt.Sprintf("%d", result)))
	var cur string
	cur = "/"
	for {
		_, text, _ = c.Read(context.Background())
		if string(text) == "END" {
			if err := conn.Quit(); err != nil {
				log.Println(err)
				return
			}
		}
		fmt.Println("received", string(text))
		cur = handle(string(text), cur)
	}
}

func login(user string, pass string) int {
	var err error
	conn, err = ftp.Dial("students.yss.su:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Println(err)
		return 1
	}
	err = conn.Login(user, pass)
	if err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func main() {
	http.HandleFunc("/ws", HomeRouterHandler)
	err := http.ListenAndServe("127.0.0.1:9457", nil)
	if err != nil {
		fmt.Println(err)
	}
}
