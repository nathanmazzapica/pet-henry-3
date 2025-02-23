package server

import "net/http"

func InitRoutes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", ServeHome)
	http.HandleFunc("/error", ServeError)
}

func StartHTTP() error {
	return http.ListenAndServe("127.0.0.1:8080", nil)
}
