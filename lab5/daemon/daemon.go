package main

import (
	"database/sql"
	"fmt" // пакет для форматированного ввода вывода
	"github.com/SlyMarbo/rss"
	_ "github.com/go-sql-driver/mysql"
	"log" // пакет для логирования
	"time"
)

var db *sql.DB

func update() int {
	log.Println("Went into update")
	feed, errrss := rss.Fetch("https://neftegaz.ru/export/yandex.php")
	if errrss != nil {
		panic(errrss)
	}
	log.Println("Fetched rss")
	db.Exec("DELETE FROM `iu9Limonov` WHERE 1")
	c := 0
	for ind, val := range feed.Items {
		if ind >= 10 {
			break
		}
		tit := val.Title
		text := val.Date
		log.Println(text)
		log.Println("Starting exec")
		fmt.Println(db == nil)
		fmt.Println(db.Ping())
		query := "INSERT INTO `iu9Limonov` (`title`, `date`) VALUES(?, ?);"
		log.Println(query)
		_, err := db.Exec(query, tit, text)
		if err != nil {
			panic(err)
		}
		log.Println("Executed")
		c = c + 1
	}
	log.Println("Execution complete")
	return c
}

func main() {
	db, _ = sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su:3306)/iu9networkslabs")
	fmt.Println(db == nil)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	for {
		update()
		time.Sleep(time.Second * 100)
	}
}
