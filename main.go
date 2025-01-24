package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

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

func manageRoomLifecycle(rt *realtimechat.RealTimeChat) {
	for {
		// sleep for 1.5 hours between cleanup cycles
		time.Sleep(90 * time.Minute)
		rooms := rt.ListRooms()
		for _, roomName := range rooms {
			if room, err := rt.GetRoom(roomName); err == nil {
				room.Stop()
			}
		}
	}
}

func main() {
	// init without max clients as it's per room now
	chat := realtimechat.NewRealTimeChat()

	// start room lifecycle management
	go manageRoomLifecycle(chat)

	wsh := websockethandler.NewWebSocketHandler(chat)

	http.HandleFunc("/ws", wsh.HandleRealTimeChat)
	http.HandleFunc("/", serveChat)

	portStr := "8080"
	fmt.Printf("listening port %+v\n", portStr)

	http.ListenAndServe(":"+portStr, nil)
}
