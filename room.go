package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	forward chan []byte      // 他のクライアントに転送するためのメッセージを保持するチャネル
	join    chan *client     // チャットルームに参加を試みるクライアントのためのチャネル
	leave   chan *client     // チャットルームからの退室を試みるクライアントのためのチャネル
	clients map[*client]bool // 在室している全てのクライアントが保持される
	change chan *client
}

type rooms struct {
	rooms map[int]*room
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}
var roomList rooms

func (r *room) run() {

	for {
		fmt.Println("送信されてるよ")
		select {
			// r.joinをclientに入れている
			case client := <-r.join:
				log.Println("入室")
				r.clients[client] = true
			// 退室
			case client := <-r.leave:
				log.Println("退室")
				delete(r.clients, client)

			// msg受取
			case msg := <-r.forward:
				log.Println("文字受取")
				log.Println(msg)
			//	ひとりひとりforでclientに入れて送信
			for client := range r.clients {
				log.Println("読取り")
				select {
				case client.send <- msg:
				default:
					// 失敗した人は退室させる
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request){
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Fatal(err)
	}
	client := &client{
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
		roomId: 0,
	}
	roomList.rooms = make(map[int]*room)

	r.join <- client
	defer func() { r.leave <- client}()

	//client.changeRooms(&roomList)
	go client.write()
	client.read(&roomList, r)
}