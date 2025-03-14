package server

import "net/http"

func InitRoutes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", ServeHome)
	http.HandleFunc("/error", ServeError)

	http.HandleFunc("/ws", HandleConnections)
}

func InitBroadcasts() {
	go broadcastChat()
	go broadcastEvent()
}

func StartHTTP() error {
	return http.ListenAndServe("0.0.0.0:8081", nil)
}
