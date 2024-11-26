package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)
import (
	"github.com/jlaffaye/ftp"
	"log"
)

func handle(conn *ftp.ServerConn) {
	req := "temp"
	path := "/"
	path2 := "."
	cur := "/"
	for req != "END" {
		fmt.Scan(&req)
		if req == "GET" {
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
			log.Println("Done")
		} else if req == "MKDIR" {
			fmt.Println("Name of folder")
			fmt.Scan(&path)
			err := conn.MakeDir(cur + path)
			if err != nil {
				log.Println(err)
			}
			log.Println("Done")
		} else if req == "RMF" {
			fmt.Println("Name of file")
			fmt.Scan(&path)
			err := conn.Delete(path)
			if err != nil {
				log.Println(err)
			}
			log.Println("Done")
		} else if req == "LS" {
			entry, err := conn.List(cur)
			if err != nil {
				log.Println(err)
			}
			for _, elem := range entry {
				fmt.Println(elem.Name, elem.Type, elem.Time)
			}
		} else if req == "CD" {
			fmt.Println("Directory name")
			fmt.Scan(&path)
			err := conn.ChangeDir(cur + path)
			cur = cur + path + "/"
			if err != nil {
				log.Println(err)
			}
			log.Println("Done")
		} else if req == "RMD" {
			fmt.Println("Directory name")
			fmt.Scan(&path)
			err := conn.RemoveDir(cur + path)
			cur = cur + path + "/"
			if err != nil {
				log.Println(err)
			}
			log.Println("Done")
		} else if req == "RMDR" {
			fmt.Println("Directory name")
			fmt.Scan(&path)
			err := conn.RemoveDirRecur(cur + path)
			cur = cur + path + "/"
			if err != nil {
				log.Println(err)
			}
			log.Println("Done")
		}
	}
}

func main() {
	c, err := ftp.Dial("", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login("", "")
	if err != nil {
		log.Fatal(err)
	}

	handle(c)

	if err := c.Quit(); err != nil {
		log.Fatal(err)
	}
}
