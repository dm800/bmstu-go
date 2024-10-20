package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/go-getter/v2"
	log "github.com/mgutz/logxi/v1"
)

func serveClient(response http.ResponseWriter, request *http.Request) {
	os.Remove("./temp/index.html")
	path := request.URL.Path
	path = strings.Replace(path, ":/", "://", 1)
	path = strings.TrimLeft(path, "/")
	got, err := getter.GetFile(context.Background(), "./temp/index.html", path)
	if err != nil {
		log.Error("xd", err)
	} else {
		log.Info(got.Dst)
	}
	log.Info("got request", "Method", request.Method, "Path", path)
	indexHtml := template.Must(template.ParseFiles("./temp/index.html"))
	if err := indexHtml.Execute(response, nil); err != nil {
		log.Error("HTML creation failed", "error", err)
	} else {
		log.Info("response sent to client successfully")
	}
}

func main() {
	http.HandleFunc("/", serveClient)
	log.Info("starting listener")
	log.Error("listener failed", "error", http.ListenAndServe("127.0.0.1:8123", nil))
}
