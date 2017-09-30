package session

import (
	"net/http"
)

// key 不能带  ','  ';'  ' '
func newCookie(w http.ResponseWriter, key string, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:    key,
		Value:   value,
		Path:    "/",
	})
}

func delCookie(w http.ResponseWriter, key string) {
	http.SetCookie(w, &http.Cookie{
		Name:   key,
		Path:   "/",
		MaxAge: -1,
	})
}
