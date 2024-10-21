package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/go-getter/v2"
	log "github.com/mgutz/logxi/v1"
)

func replacer(path string, addr string) {
	content, _ := os.ReadFile(path)
	text := string(content)
	text = strings.ReplaceAll(text, "/static", "static")
	text = strings.ReplaceAll(text, "\n", "")
	prev := 0
	for {
		start1 := strings.Index(text[prev:], "href") + 4
		start2 := strings.Index(text[prev:], "src") + 3
		start := min(start1, start2)
		if start == 2 {
			break
		}
		substr := text[start+prev : start+prev+30]
		fmt.Println("start", start, "prev", prev)
		if strings.Contains(substr, global_ip) {
			prev = prev + start + 3 + len(global_ip)
			continue
		}
		if strings.Contains(substr, "://") {
			text = text[:start+prev+2] + global_ip + "/" + text[start+prev+2:]
			prev = prev + start + 3 + len(global_ip)
		} else {
			text = text[:start+prev+2] + global_ip + "/" + addr + "/" + text[start+prev+2:]
			prev = prev + start + 3 + len(global_ip) + len(addr) + 2
		}
	}
	err := os.WriteFile(path, []byte(text), 0644)
	if err != nil {
		fmt.Println("replacer failed")
	} else {
		fmt.Println("replacer succeeded")
	}
}

func serveClient(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	path = strings.Trim(path, "/")
	path = strings.Replace(path, ":/", "://", 1)
	fmt.Println(path)
	if path == "favicon.ico" {
		_, err := getter.GetFile(context.Background(), "./temp/favicon.ico", path)
		if err != nil {
			log.Error("icon error", err)
		} else {
			log.Info("served an icon")
		}
		return
	}
	pathname := strings.Split(path, "/")[2:]
	predres := strings.Split(pathname[len(pathname)-1], ".")
	res := "." + predres[len(predres)-1]
	is_html := false
	if res == ".html" {
		is_html = true
	} else if res != ".css" && res != ".js" && res != ".png" && res != ".jpg" && res != ".jpeg" && res != ".ico" && res != ".gif" {
		pathname = append(pathname, "index.html")
		is_html = true
	}
	pathres := "./temp/" + strings.Join(pathname, "/")
	os.Remove(pathres)
	got, err := getter.GetFile(context.Background(), pathres, path)
	if err != nil {
		log.Error("get error", err)
	} else {
		log.Info(got.Dst)
	}
	log.Info("got request", "Method", request.Method, "Path", pathres)
	if is_html {
		replacer(pathres, path)
		indexHtml := template.Must(template.ParseFiles(pathres))
		if err := indexHtml.Execute(response, nil); err != nil {
			log.Error("HTML creation failed", "error", err)
		} else {
			log.Info("response sent to client successfully")
		}
	} else {
		http.ServeFile(response, request, pathres)
	}
}

var ip = "192.168.31.116:9456"
var global_ip = "http://46.138.248.69:9456"

func main() {
	http.HandleFunc("/", serveClient)
	log.Info("starting listener")
	log.Error("listener failed", "error", http.ListenAndServe(ip, nil))
}
