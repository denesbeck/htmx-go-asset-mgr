package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Data struct {
	Name           string
	Serial         string
	Email          string
	Country        string
	AssetSerial    string `json:"asset_serial"`
	AssetType      string `json:"asset_type"`
	ExpirationDate string `json:"expiration_date"`
}

var payload []Data

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("../index.html")
	tmpl.Execute(w, payload)
}

func AssetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		vars := mux.Vars(r)
		id := vars["id"]

		for i, asset := range payload {
			if asset.AssetSerial == id {
				payload = append(payload[:i], payload[i+1:]...)
			}
		}
	}
	if r.Method == "PUT" {
		tmpl, _ := template.ParseFiles("../insert.html")
		tmpl.Execute(w, payload)
	}
}

func main() {
	data, err := os.ReadFile("../data.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		log.Fatal("Error when parsing data: ", err)
	}
	router := mux.NewRouter()

	router.HandleFunc("/", HomeHandler)

	router.HandleFunc("/assets/{id}/", AssetHandler)

	http.Handle("/", router)

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
