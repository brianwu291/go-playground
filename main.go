package main

import (
	"fmt"
	"html/template"
	"net/http"

	websockethandler "github.com/brianwu291/go-playground/handlers/websocket"
	realtimechat "github.com/brianwu291/go-playground/realtimechat"
)

func serveChat(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/chat.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {

	chat := realtimechat.NewRealTimeChat(10)
	chat.Run()
	defer chat.Stop()

	wsh := websockethandler.NewWebSocketHandler(chat)

	http.HandleFunc("/ws", wsh.HandleRealTimeChat)
	http.HandleFunc("/", serveChat)


	portStr := "8080"
	fmt.Printf("listening port %+v\n", portStr)

	http.ListenAndServe(":"+portStr, nil)
}
