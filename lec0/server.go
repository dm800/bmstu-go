package main

import (
	"fmt"      // пакет для форматированного ввода вывода
	"log"      // пакет для логирования
	"net/http" // пакет для поддержки HTTP протокола
	"strings"  // пакет для работы с  UTF-8 строками

	"github.com/SlyMarbo/rss"

	"github.com/coder/websocket"
)

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer c.CloseNow()

	r.ParseForm()       //анализ аргументов,
	fmt.Println(r.Form) // ввод информации о форме на стороне сервера
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	feed, _ := rss.Fetch("https://neftegaz.ru/export/yandex.php")
	for v := range feed.Items {
		val := feed.Items[v]
		fmt.Fprintf(w, "<b>"+val.Title+"<\b><br>")

	}
	c.Close(websocket.StatusNormalClosure, "all good")
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)  // установим роутер
	err := http.ListenAndServe(":9000", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
