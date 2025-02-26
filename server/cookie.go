package server

import (
	"net/http"
	"strings"
	"time"
)

const UserIDCookieName = "user_id_daisy"

func NewUserCookie(userID string, r *http.Request) *http.Cookie {
	domain := ""

	if strings.Contains(r.Host, "pethenry.com") {
		domain = ".pethenry.com"
	}

	cookie := http.Cookie{
		Name:     UserIDCookieName,
		Value:    userID,
		HttpOnly: true,
		Expires:  time.Now().AddDate(10, 0, 0),
		Domain:   domain,
	}

	return &cookie
}

func GetCookie(r *http.Request, name string) (*http.Cookie, error) {
	cookie, err := r.Cookie(name)
	return cookie, err
}

func GetUserIDCookie(r *http.Request) (*http.Cookie, error) {
	return GetCookie(r, UserIDCookieName)
}
