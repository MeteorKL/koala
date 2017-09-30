package koala

import (
	"fmt"
	"html/template"
	"net/http"
	"log"
	"encoding/json"
)

var RenderPath string

func Render(w http.ResponseWriter, file string, data interface{}) {
	t, err := template.New(file).ParseFiles(RenderPath + file)
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, data)
}

func Relocation(w http.ResponseWriter, URL string) {
	// w.Header().Set("Location", URL)
	t, err := template.New("x").Parse("<script>window.location.href='" + URL + "';</script>")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

func Back(w http.ResponseWriter) {
	t, err := template.New("x").Parse("<script>history.go(-1);</script>")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

func WriteJSON(w http.ResponseWriter, data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(json)
}
