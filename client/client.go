package main

import (
	"PubsubServer/types"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
)

var (
	c = client{}
)

type client struct {
	MessageChan   chan types.CommunicationMessage   `json:"message_chan"`
	Conn          *websocket.Conn                   `json:"conn"`
}

func (c *client) SetClient() {
	c.MessageChan = make(chan types.CommunicationMessage, 1)
}

func StartClient() {
	//start client
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:3030", Path: "/pubsub"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("got error: %v", err)
	}
	c.SetClient()
	c.Conn = conn
	go c.Receive(conn)
	go c.SendMessage()
	go c.loopSend()

	defer func() {
		err := recover()
		if err != nil {
			log.Fatalf("recover err: %v", err)
		}
	}()

	var stopSignal chan int
	select {
	case _ = <- stopSignal:
		log.Print("exit here")
		os.Exit(0)
	}
}

func (c *client) Receive(conn *websocket.Conn) {
	defer func() {
		err := recover()
		if err != nil {
			log.Fatalf("recover err: %v", err)
		}
	}()
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("got error: %v", err)
			continue
		}
		log.Printf("receive: %v", string(data))
	}
}


func(c *client) SendMessage() {
	for {
		var method string
		var topic string
		fmt.Println("input:")
		_, _ = fmt.Scanf("%s %s", &method, &topic)
		msg := types.NewCommunicationMessage(method, topic)

		c.MessageChan <- msg
	}
}

func (c *client) loopSend() {
	for {
		msg := <- c.MessageChan
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("got error: %v", err)
			continue
		}
	}
}