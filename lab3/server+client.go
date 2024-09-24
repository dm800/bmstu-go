package main

import (
	"net"
	"net/http"
	"os"
	"strings"

	log "github.com/mgutz/logxi/v1"
)

func serveClient(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	log.Info("got request", "Method", request.Method, "Path", path)
	if path != "/" {
		log.Error("invalid path", "Path", path)
		response.WriteHeader(http.StatusNotFound)
	} else {
		log.Info("response sent to client successfully")
	}
}

func ping(conn *net.TCPConn) {

}

func main() {
	http.HandleFunc("/", serveClient)
	file, _ := os.Open("config.txt")
	data := make([]byte, 128)
	n, _ := file.Read(data)
	str := string(data[:n])
	ips := strings.Split(str, "\n")
	/*for _, ip := range(ips) {
		ips
	}*/
	for _, ip := range ips {
		str, _ := net.ResolveTCPAddr("tcp", ip)
		conn, _ := net.DialTCP("tcp", nil, str)
		ping(conn)
	}
	log.Info("starting listener")
	log.Error("listener failed", "error", http.ListenAndServe("127.0.0.1:8123", nil))
	file.Close()
}
