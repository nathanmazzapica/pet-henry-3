package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"pet-henry-3/data"
	"pet-henry-3/models"
	"sync"
)

type Action struct {
	actionType string
	data       map[string]interface{}
}

// fine for development
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients     = make(map[*models.Client]bool)
	clientsLock sync.RWMutex
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("Establishing new connection from: " + r.RemoteAddr)

	log.Println("Fetching user ID for: " + r.RemoteAddr)

	userID, err := GetUserIDCookie(r)

	if err != nil {
		log.Println("Error establishing connection with", r.RemoteAddr, ":", err)
		return
	}

	log.Println("Fetching user data for", r.RemoteAddr)

	user, err := data.GetUserFromDB(userID.Value)

	if err != nil {
		log.Println("Error establishing connection with", r.RemoteAddr, ":", err)
		return
	}

	log.Println("Finalizing websocket connection for", r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error establishing connection with", r.RemoteAddr, ":", err)
		return
	}

	client := &models.Client{Conn: conn, User: user}

	clientsLock.Lock()
	clients[client] = true
	clientsLock.Unlock()

	defer closeConnection(client)

	readMessages(client)

}

func readMessages(client *models.Client) {
	conn := client.Conn
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// handle error properly
			break
		}

		log.Println("Received message:", string(msg))
	}
}

func closeConnection(client *models.Client) {
	clientsLock.Lock()
	delete(clients, client)
	clientsLock.Unlock()
	client.Conn.Close()
	// player leave notification
}
