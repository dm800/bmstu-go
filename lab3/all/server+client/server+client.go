package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"example.com/all/proto"

	log "github.com/mgutz/logxi/v1"
	"github.com/skorobogatov/input"

	"github.com/coder/websocket"
)

var addrStr string
var port string
var msg string
var c *websocket.Conn

type Setter struct {
	Tkey string `json:"tkey"`
}

// Client - состояние клиента.
type Client struct {
	logger log.Logger    // Объект для печати логов
	conn   *net.TCPConn  // Объект TCP-соединения
	enc    *json.Encoder // Объект для кодирования и отправки сообщений
	ip     string        // ip
}

// NewClient - конструктор клиента, принимает в качестве параметра
// объект TCP-соединения.
func NewClient(conn *net.TCPConn) *Client {
	return &Client{
		logger: log.New(fmt.Sprintf("client %s", conn.RemoteAddr().String())),
		conn:   conn,
		enc:    json.NewEncoder(conn),
		ip:     "",
	}
}

// serve - метод, в котором реализован цикл взаимодействия с клиентом.
// Подразумевается, что метод serve будет вызаваться в отдельной go-программе.
// server
func (client *Client) serve() {
	defer client.conn.Close()
	decoder := json.NewDecoder(client.conn)
	var req proto.Request
	if err := decoder.Decode(&req); err != nil {
		client.logger.Error("cannot decode message", "reason", err)
	} else {
		client.logger.Info("received command", "command", req.Command)
		if client.handleRequest(&req) {
			client.logger.Info("shutting down connection")
		}
	}
}

// handleRequest - метод обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
// server
func (client *Client) handleRequest(req *proto.Request) bool {
	log.Info("WENT INTO WITH", req.Command)
	switch req.Command {
	case "quit":
		client.respond("ok", nil)
		return true
	case "set":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var msg string
			if err := json.Unmarshal(*req.Data, &msg); err != nil {
				errorMsg = "malformed data field"
			} else {
				log.Info("Msg set")
			}
		}
		if errorMsg == "" {
			client.respond("ok", nil)
		} else {
			client.logger.Error("addition failed", "reason", errorMsg)
			client.respond("failed", errorMsg)
		}
	case "get":
		client.respond("result", msg)
		client.logger.Info("Sended back msg")
	default:
		client.logger.Error("unknown command")
		client.respond("failed", "unknown command")
	}
	return false
}

// client
func interact(conn *net.TCPConn, command string) {
	defer conn.Close()
	encoder, decoder := json.NewEncoder(conn), json.NewDecoder(conn)
	// Отправка запроса.
	switch command {
	case "quit":
		send_request(encoder, "quit", nil)
		return
	case "get":
		send_request(encoder, "get", nil)
	default:
		fmt.Printf("error: unknown command\n")
	}

	// Получение ответа.
	var resp proto.Response
	if err := decoder.Decode(&resp); err != nil {
		fmt.Printf("decode error: %v\n", err)
		return
	}

	// Вывод ответа в стандартный поток вывода.
	switch resp.Status {
	case "ok":
		log.Info("ok\n")
	case "failed":
		if resp.Data == nil {
			log.Info("error: data field is absent in response\n")
		} else {
			var errorMsg string
			if err := json.Unmarshal(*resp.Data, &errorMsg); err != nil {
				log.Info("error: malformed data field in response\n")
			} else {
				log.Info("failed: %s\n", errorMsg)
			}
		}
	case "result":
		if resp.Data == nil {
			fmt.Printf("error: data field is absent in response\n")
		} else {
			var ans string
			if err := json.Unmarshal(*resp.Data, &ans); err != nil {
				fmt.Printf("error: malformed data field in response\n")
			} else {
				fmt.Printf("%s\n", ans)
			}
		}
	default:
		fmt.Printf("error: server reports unknown status %q\n", resp.Status)
	}
}

func send_request(encoder *json.Encoder, command string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	encoder.Encode(&proto.Request{Command: command, Data: &raw})
}

func (client *Client) respond(status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	client.enc.Encode(&proto.Response{Status: status, Data: &raw})
}

// server part
func main_cycle(listener net.TCPListener) {
	for {
		if conn, err := listener.AcceptTCP(); err != nil {
			log.Error("cannot accept connection", "reason", err)
		} else {
			log.Info("accepted connection", "address", conn.RemoteAddr().String())

			// Запуск go-программы для обслуживания клиентов.
			go NewClient(conn).serve()
		}
	}
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("went into handler")
	var err error
	c, err = websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		fmt.Println(err)
	}
	msgsend()
}

func msgsend() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	if c == nil {
		log.Info("there is no websocket connection")
		return
	}
	wr, err := c.Writer(ctx, websocket.MessageType(websocket.MessageText))
	if err != nil {
		log.Error("sending to websocket err:", err)
		return
	}
	log.Info("trying to write msg")
	wr.Write([]byte(msg))
	wr.Close()
	log.Info("i believe that we sent msg")
}

func runsocket() {
	http.HandleFunc("/", HomeRouterHandler) // установим роутер
	sport, _ := strconv.Atoi(port)
	log.Info("Startig listening ws at", addrStr+":"+strconv.Itoa(sport+1))
	err := http.ListenAndServe(addrStr+":"+strconv.Itoa(sport+1), nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func hhelper(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Access-Control-Allow-Origin", "*")
	log.Info("went into hhandler")
	if r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		log.Info("unparsed", string(body))
		var ans Setter
		err := json.Unmarshal(body, &ans)
		if err != nil {
			log.Error("error", err)
		}
		log.Info("parsed ", ans.Tkey)
		msg = ans.Tkey
		log.Info("set text to ", msg)
	}
	msgsend()
}

func hhandler() {
	http.HandleFunc("/h", hhelper) // установим роутер
	sport, _ := strconv.Atoi(port)
	err := http.ListenAndServe(addrStr+":"+strconv.Itoa(sport-1), nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//9456 - LISTENING PORT FOR P2P
//9457 - LISTENING PORT FOR WS
//9455 - LISTENING PORT FOR POST

func main() {
	c = nil
	go runsocket()
	go hhandler()
	fl := true
	msg = " "
	flag.StringVar(&addrStr, "addr", "192.168.31.116", "ip address")
	flag.StringVar(&port, "p", "9456", "port")
	flag.Parse()
	file, _ := os.Open("config.txt")
	data := make([]byte, 128)
	n, _ := file.Read(data)
	str := string(data[:n])
	ips := strings.Split(str, "\n")
	// Разбор адреса, строковое представление которого находится в переменной addrStr.
	// server part
	if addr, err := net.ResolveTCPAddr("tcp", addrStr+":"+port); err != nil {
		log.Error("address resolution failed", "address", addrStr+":"+port)
	} else {
		log.Info("resolved TCP address", "address", addr.String())

		// Инициация слушания сети на заданном адресе.
		if listener, err := net.ListenTCP("tcp", addr); err != nil {
			log.Error("listening failed", "reason", err)
		} else {
			// Цикл приёма входящих соединений.
			go main_cycle(*listener)
		}
	}
	for fl {
		// Чтение команды из стандартного потока ввода
		command := input.Gets()
		log.Info("executing", command)

		if command == "set" {
			fmt.Printf("type your text: ")
			msg = input.Gets()
			fmt.Printf("your text: %s\n", msg)
			msgsend()
			continue
		}

		if command == "get" {
			fmt.Printf("Your: %s \n", msg)
		}

		for _, ip := range ips {
			if !(fl) {
				break
			}
			ip = strings.TrimSpace(ip)
			if addrStr+":"+port == ip {
				continue
			}
			fmt.Printf("ip %s: ", ip)
			if addr, err := net.ResolveTCPAddr("tcp", ip); err != nil {
				log.Error("error: ", err)
			} else if conn, err := net.DialTCP("tcp", nil, addr); err != nil {
				log.Info(ip, " is not available")
			} else {
				if command == "quit" {
					interact(conn, "quit")
					fl = false
				}
				switch command {
				case "get":
					interact(conn, "get")
				default:
					fmt.Printf("error: unknown command\n")
					continue
				}
			}
		}
	}
	file.Close()
}
