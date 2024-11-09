package main

import (
	"database/sql"
	"fmt" // пакет для форматированного ввода вывода
	"github.com/SlyMarbo/rss"
	_ "github.com/go-sql-driver/mysql"
	"log"      // пакет для логирования
	"net/http" // пакет для поддержки HTTP протокола
	"strings"  // пакет для работы с  UTF-8 строками
	"time"
)

var db *sql.DB

func update() {
	log.Println("Went into update")
	feed, errrss := rss.Fetch("https://neftegaz.ru/export/yandex.php")
	if errrss != nil {
		panic(errrss)
	}
	log.Println("Fetched rss")
	for ind, val := range feed.Items {
		if ind > 10 {
			return
		}
		tit := val.Title
		text := val.Content
		log.Println("Starting exec")
		fmt.Println(db == nil)
		fmt.Println(db.Ping())
		query := "INSERT INTO `iu9Limonov` (`title`, `content`) VALUES(?, ?);"
		log.Println(query)
		_, err := db.Exec(query, tit, text)
		if err != nil {
			panic(err)
		}
		log.Println("Executed")
	}
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //анализ аргументов,
	fmt.Println(r.Form) // ввод информации о форме на стороне сервера
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	update()
	fmt.Fprintf(w, "Update completed")
}

// INSERT INTO `iu9Limonov` (`title`, `content`)
// VALUES ('Test2', 'Testing');

func main() {
	fmt.Println(db == nil)
	db, _ = sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su:3306)/iu9networkslabs")
	fmt.Println(db == nil)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	http.HandleFunc("/", HomeRouterHandler)  // установим роутер
	err := http.ListenAndServe(":9000", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
