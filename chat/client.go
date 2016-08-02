package chat

import (
	"golang.org/x/net/websocket"
	"fmt"
	"io"
)
var  MaxId int
type Client struct {
	id int
	ws *websocket.Conn
	server *Server
	ch chan *Message
	doneCh chan bool
}

func NewClient(ws *websocket.Conn,server *Server)  *Client {
	if ws==nil {
		fmt.Println("invalid websocket")

	}
	if server==nil {
		fmt.Println("invalid server")
	}
	MaxId++
	ch:=make(chan *Message)
	doneCh:=make(chan bool)
	return &Client{
		MaxId,
		ws,
		server,
		ch,
		doneCh,
	}




}
func (client *Client) Listen(){
	go client.ListenWrite()
	client.ListenRead()
}
func (client *Client)ListenWrite()  {
fmt.Println("Listen write")
	for {
		select {
		case message:=<-client.ch :
			websocket.JSON.Send(client.ws,message)
		case <-client.doneCh :
			client.server.Del(client)
			client.doneCh<-true
		}
	}

}
func (client *Client) write(message *Message)  {
select {
	case client.ch<-message:
	default:
	client.server.Del(client)
	fmt.Println("Client is disconnected");

	}
}
func (client *Client)ListenRead()  {

	for{
		select {
		case <-client.doneCh :
			client.server.Del(client)
			fmt.Println("Delete client")
		default :
		var message Message
		  error:=websocket.JSON.Receive(client.ws,&message)
			if error==io.EOF {
				fmt.Println("EOF")
				client.doneCh<-true
			}else if error!=nil {
				fmt.Print(error)
			}else {
				client.server.SendAllMessage(&message)
			}

		}
	}
}