package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"pet-henry-3/data"
	"pet-henry-3/game"
	"pet-henry-3/models"
	"sync"
)

type Event struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// fine for development
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	clients     = make(map[*models.Client]bool)
	clientsLock sync.RWMutex
	chats       = make(chan Event)
	events      = make(chan Event)
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("Establishing new connection from: " + r.RemoteAddr)

	log.Println("Fetching user ID for: " + r.RemoteAddr)

	userID, err := GetUserIDCookie(r)

	if err != nil {
		log.Println("Error establishing connection with", r.RemoteAddr, ":", err)
		return
	}

	// check if user in cache
	// if yes, retrieve from cache
	// else, retrieve from db & cache
	log.Println("TESTING USER FETCH CACHE")
	user, err := data.GetFromHash(userID.Value)

	if user == nil {
		log.Println("Fetching user data for", r.RemoteAddr, "from DB")
		user, err = data.GetUserFromDB(userID.Value)
		// reflect on error handling .. .
		data.InsertIntoHash(user)
	} else {
		log.Printf("Retrieved user data for %s{%v} from cache", r.RemoteAddr, user.DisplayName)
	}

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

		log.Printf("Received message from %s: %v", client.User.DisplayName, string(msg))
		handleMessage(client, msg)
	}
}

func handleMessage(client *models.Client, msg []byte) {
	var clientEvent Event
	if err := json.Unmarshal(msg, &clientEvent); err != nil {
		log.Println("Invalid JSON:", err)
		return
	}

	switch clientEvent.Type {
	case "pet":
		handlePet(client, clientEvent)
	case "chat":
		handleChatMessage(client, clientEvent)
	default:
		log.Println("Unknown clientEvent:", clientEvent.Type)
	}
}

func handlePet(client *models.Client, petEvent Event) {
	game.IncrementPetCounter()
	petEvent.Data["count"] = game.Counter
	events <- petEvent
	client.User.PetDaisy()
	data.IncrementUserPet(client.User.UserID)
}

func handleChatMessage(client *models.Client, chat Event) {
	_, ok := chat.Data["message"].(string)
	if !ok {
		log.Println("Invalid message:", chat.Data["message"])
		return
	}
	// check for potty language later
	chat.Data["sender"] = client.User.DisplayName
	chats <- chat
}

func broadcastChat() {
	for {
		newChat := <-chats

		for client := range clients {
			networkEvent(client, prepareEventJSON(newChat))
		}
	}
}

func broadcastEvent() {
	for {
		newEvent := <-events
		for client := range clients {
			networkEvent(client, prepareEventJSON(newEvent))
		}
	}
}

func closeConnection(client *models.Client) {
	clientsLock.Lock()
	delete(clients, client)
	clientsLock.Unlock()
	client.Conn.Close()
	// player leave notification
}
