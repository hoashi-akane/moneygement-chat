package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

type Message struct {
	Id       int    `json:"id"`
	LedgerId int    `json:"ledgerId"`
	UserId   int    `json:"UserId"`
	Message  string `json:"Message"`
	RoomId   int    `json:"roomId"`
	NickName string `json:"nickName"`
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	roomList.rooms = make(map[int]*room)

	r := newRoom()
	http.Handle("/room", r)
	go r.run()

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}