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
				break
			}
		}
	}
	if r.Method == "PUT" {
		vars := mux.Vars(r)
		id := vars["id"]

		for i, asset := range payload {
			if asset.AssetSerial == id {
				payload[i].Name = r.FormValue("name")
				payload[i].Serial = r.FormValue("serial")
				payload[i].Email = r.FormValue("email")
				payload[i].Country = r.FormValue("country")
				payload[i].AssetSerial = r.FormValue("asset_serial")
				payload[i].AssetType = r.FormValue("asset_type")
				payload[i].ExpirationDate = r.FormValue("expiration_date")
				break
			}
		}
	}
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		vars := mux.Vars(r)
		id := vars["id"]

		var row Data
		for i, asset := range payload {
			if asset.AssetSerial == id {
				row = payload[i]
				break
			}
		}
		tmpl, _ := template.ParseFiles("../insert.html")
		tmpl.Execute(w, row)
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
	router.HandleFunc("/assets/{id}", AssetHandler)
	router.HandleFunc("/assets/{id}/edit", EditHandler)

	http.Handle("/", router)

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
