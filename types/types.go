package types

import "time"

type CommunicationMessage struct {
	//call method
	Method       string     `json:"method"`
	//arguments
	Args         []string   `json:"args"`
}

func NewCommunicationMessage(method string,arg... string) CommunicationMessage {
	var args = make([]string, 0)
	for _, v := range arg {
		args = append(args, v)
	}
	return CommunicationMessage{
		Method: method,
		Args: args,
	}
}


type ServerMesage struct {
	Time    time.Time    	 `json:"time"`
	Message  string          `json:"message"`
}

func NewServerMessage(msg string) ServerMesage {
	return ServerMesage{
		Time: time.Now().Local(),
		Message: msg,
	}
}