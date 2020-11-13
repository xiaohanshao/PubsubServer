package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"PubsubServer/types"
	"github.com/gorilla/websocket"
)

var (
	s = Server{}
	ug       = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

)

type Server struct {
	MessageChan   chan types.ServerMesage  `json:"message_chan"`
	Connections   []*websocket.Conn        `json:"connections"`
}

func(s *Server) SetServer() {
	s.MessageChan = make(chan types.ServerMesage, 1)
	s.Connections = make([]*websocket.Conn, 0)
}

func StartServer() {
	// start server listen

	go s.SendMessage()
	go s.loopSend()
	http.HandleFunc("/pubsub", HttpHanlde)
	s.SetServer()

	err := http.ListenAndServe(":3030", nil)
	if err != nil {
		log.Fatalf("got listen_server error: %v", err)
	}

	defer func() {
		err := recover()
		if err != nil {
			log.Fatalf("recover err: %v", err)
		}
	}()

	var stopSignal chan int

	select {
	case _ = <- stopSignal:
		log.Println("exit here")
		os.Exit(0)
	}

}

func HttpHanlde(w http.ResponseWriter, r *http.Request) {
	//get connection.
	conn, err := ug.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("got connect error: %v", err)
		return
	}
	log.Print("here, to receive...")

	go s.Receive(conn)
}

func(s *Server) Receive(con *websocket.Conn) {

	defer func() {
		err := recover()
		switch rt := err.(type) {
		default:
			log.Printf("got err: %s", rt)
		}
	}()
	log.Print("in receive \n")
	for {
		log.Print("in receive \n")
		_, data, err := con.ReadMessage()
		if err != nil {
			log.Printf("read invalid message: %s", err.Error())
			//continue
		}
		log.Print("in receive \n")
		var requestData types.CommunicationMessage
		err = json.Unmarshal(data, &requestData)
		if err != nil {
			log.Printf("read invalid message: %s", err.Error())
			//continue
		}
		s.Connections = append(s.Connections, con)
		log.Printf("got meesage from client :%s", string(data))
	}
}

func(s *Server) SendMessage() {
	for {
		var sendMsg string
		fmt.Println("input:")
		_, _ = fmt.Scanf("%s", &sendMsg)
		msg := types.NewServerMessage(sendMsg)

		s.MessageChan <- msg
	}
}

func (s *Server) loopSend() {
	for {
		message := <- s.MessageChan
		for _, con := range s.Connections {
			err := con.WriteJSON(message)
			if err != nil {
				log.Printf("send message: %s to client failed", err.Error())
				continue
			}
		}
	}
}