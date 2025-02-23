package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"pet-henry-3/data"
	"pet-henry-3/game"
	"pet-henry-3/models"
)

func prepareUserData(w http.ResponseWriter, r *http.Request) *models.User {
	var user *models.User

	userIDCookie, err := r.Cookie("uid_temp")

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
		WS_URL:    "ws://localhost:8080/ws",
	}
}
