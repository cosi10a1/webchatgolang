package chat

import (
"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

type Server struct {
	name string
	messages []*Message
	clients map[int]*Client
	addCh chan *Client
	delCh chan *Client
	doneCh chan bool
	sendallCh chan *Message
}

func  NewServer(name string) *Server  {
	messages :=[]*Message{}
	clients :=make(map[int]*Client)
	addCh:=make(chan *Client)
	delCh:=make(chan *Client)
	doneCh:=make(chan bool)
	senallChanel:=make(chan *Message)
	return  &Server{name,
		messages,
		clients,
		addCh,
		delCh,
		doneCh,
		senallChanel,
	}
	
}
func (server *Server) AddClient(client *Client)  {

}
func (server *Server) Add(client *Client){
	server.addCh<-client
}
func (server *Server) Del(client *Client){
	server.delCh<-client

}
func (server *Server)SendAllMessage(message *Message){
server.sendallCh<-message
}
func (server * Server) sendall(message *Message){
	for _,client:=range  server.clients {
		client.write(message)
	}
}
func (server *Server) sendpast(client *Client){
	for _,message:=range server.messages  {
		client.write(message)
	}



}
func (server *Server) Listen(){
	fmt.Println("Listenning to server")
	onConnected:= func(ws *websocket.Conn) {
		defer func() {
			err:=ws.Close()
			if err!=nil {
				fmt.Println("Error while close websocket")

			}
		}()
		client:=NewClient(ws,server)
		server.Add(client)
		client.Listen()

	}
	http.Handle(server.name,websocket.Handler(onConnected))
	for  {
		select {
		case client:=<-server.addCh :
			server.clients[client.id]=client
			fmt.Println("New client")
		case client:=<-server.delCh :
			delete(server.clients,client.id)
			fmt.Print("Delete client")
		case message:=<-server.sendallCh :
			server.messages=append(server.messages,message)
			server.sendall(message)
		case <-server.doneCh :
			return
		}

	}

}