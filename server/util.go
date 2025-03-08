package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"pet-henry-3/data"
	"pet-henry-3/game"
	"pet-henry-3/models"
)

func prepareUserData(w http.ResponseWriter, r *http.Request) *models.User {
	var user *models.User

	userIDCookie, err := GetUserIDCookie(r)

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			user = data.CreateNewUser()
			http.SetCookie(w, NewUserCookie(user.UserID, r))
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
			return nil
		}
	} else {
		user, err = data.GetUserFromDB(userIDCookie.Value)
		if err != nil {
			fmt.Println(err)
		}

	}

	return user
}

func prepareHomeData(user *models.User) models.HomeData {
	return models.HomeData{
		User:      user.DisplayName,
		SyncCode:  user.SyncCode,
		UserPets:  user.PetCount,
		TotalPets: game.Counter,
		WS_URL:    "ws://localhost:8081/ws",
	}
}

func prepareEventJSON(event Event) []byte {
	json, err := json.Marshal(event)

	if err != nil {
		log.Println("Error preparing ActionJSON:", err)
		return nil
	}

	return json
}

func networkEvent(client *models.Client, event []byte) {

	if event == nil {
		log.Println("Error preparing ActionJSON")
		return
	}

	err := client.Conn.WriteMessage(websocket.TextMessage, event)

	if err != nil {
		log.Printf("Error writing to client: %v DROPPING CONNECTION\n", err)
		client.Conn.Close()
		clientsLock.Lock()
		delete(clients, client)
		clientsLock.Unlock()
	}
}
