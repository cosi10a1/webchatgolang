package main

import (
	"net/http"
	"./chat"
)
func main() {
	server := chat.NewServer("/entry")
	go server.Listen()
	http.Handle("/", http.FileServer(http.Dir("webroot")))
	http.ListenAndServe(":8080",nil)
}
