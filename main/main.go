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

type Params struct {
	Type string
	Data Data
}

var payload []Data

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("../index.html")
	tmpl.Execute(w, payload)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
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
		tmpl.Execute(w, Params{Type: "edit", Data: row})
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		vars := mux.Vars(r)
		id := vars["id"]

		for i, asset := range payload {
			if asset.AssetSerial == id {
				payload[i].Name = r.FormValue("name")
				payload[i].Serial = r.FormValue("serial")
				payload[i].Email = r.FormValue("email")
				payload[i].Country = r.FormValue("country")
				payload[i].AssetSerial = r.FormValue("assetserial")
				payload[i].AssetType = r.FormValue("assettype")
				payload[i].ExpirationDate = r.FormValue("expirationdate")

				tmpl := template.Must(template.ParseFiles("../index.html"))
				tmpl.ExecuteTemplate(w, "assets-rows", payload[i])
				break
			}
		}
	}
}

func NewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("../insert.html")
		tmpl.Execute(w, Params{Type: "new"})
	}
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		payload = append(payload, Data{})
		i := len(payload) - 1
		payload[i].Name = r.FormValue("name")
		payload[i].Serial = r.FormValue("serial")
		payload[i].Email = r.FormValue("email")
		payload[i].Country = r.FormValue("country")
		payload[i].AssetSerial = r.FormValue("assetserial")
		payload[i].AssetType = r.FormValue("assettype")
		payload[i].ExpirationDate = r.FormValue("expirationdate")

		tmpl := template.Must(template.ParseFiles("../index.html"))
		tmpl.ExecuteTemplate(w, "assets-rows", payload[i])
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
	router.HandleFunc("/assets/new", NewHandler)
	router.HandleFunc("/assets/create", CreateHandler)
	router.HandleFunc("/assets/{id}", DeleteHandler)
	router.HandleFunc("/assets/{id}/edit", EditHandler)
	router.HandleFunc("/assets/{id}/update", UpdateHandler)

	http.Handle("/", router)

	fmt.Println("Server is listening at http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
