package server

import (
	"net/http"
	"strings"
	"time"
)

func NewUserCookie(userID string, r *http.Request) *http.Cookie {
	domain := ""

	if strings.Contains(r.Host, "pethenry.com") {
		domain = ".pethenry.com"
	}

	cookie := http.Cookie{
		Name:     "uid_temp",
		Value:    userID,
		HttpOnly: true,
		Expires:  time.Now().AddDate(10, 0, 0),
		Domain:   domain,
	}

	return &cookie
}
