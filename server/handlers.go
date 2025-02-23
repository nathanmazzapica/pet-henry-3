package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {

	user := prepareUserData(w, r)

	if user == nil {
		log.Println("Something went wrong fetching user")
	}

	fmt.Printf("USER: {%s} CONNECTED\n", user.DisplayName)

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	err := tmpl.Execute(w, prepareHomeData(user))

	if err != nil {
		fmt.Println("error sending html", err)
	}
}

func ServeError(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/error.html"))

	_ = tmpl.Execute(w, nil)
}
