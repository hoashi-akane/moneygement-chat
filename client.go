package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type client struct {
	socket *websocket.Conn // クライアントのWebSocket
	send   chan []byte     // メッセージが送られるチャネル
	room   *room           // 参加しているチャットルーム
	roomId int
}

func (c *client) read(rs *rooms, r *room) {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil{
			var message Message
			c.socket.ReadJSON(&message)

			// クライアントから送信された部屋番号と部屋が違うならtrue
			if message.RoomId != c.roomId {
				if rs.rooms[message.RoomId] == nil{
					// 部屋ができていなければ部屋を作る
					log.Println("部屋作成")
					newr := newRoom()
					rs.rooms[message.RoomId] = newr
				}
				// 部屋に入る。clientのルームidを新しく設定
				r.join <- c
				c.roomId = message.RoomId
			}else{
				// ルームidが同じならメッセージを配信
				c.room.forward <- msg
			}
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send{
		if err := c.socket.WriteMessage(websocket.TextMessage, msg);
			err != nil {
				break
		}
	}
	c.socket.Close()
}