package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// Kontak struct (Model) ...
type Kontak struct {
	ID           string `json:"ID"`
	NamaDepan    string `json:"NamaDepan"`
	NamaBelakang string `json:"NamaBelakang"`
	NoHp         string `json:"NoHp"`
	Email        string `json:"Email"`
	Alamat       string `json:"Alamat"`
}

// Get all kontak
func getContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contacts []Kontak
	sql := `SELECT
				ID,
				IFNULL(NamaDepan,''),
				IFNULL(NamaBelakang,'') NamaBelakang,
				IFNULL(NoHp,'') NoHp,
				IFNULL(Email,'') Email,
				IFNULL(Alamat,'') Alamat,
			FROM kontak`
	result, err := db.Query(sql)
	defer result.Close()
	if err != nil {
		panic(err.Error())
	}
	for result.Next() {
		var kontak Kontak
		err := result.Scan(&kontak.ID, &kontak.NamaDepan, &kontak.NamaBelakang,
			&kontak.NoHp, &kontak.Email, &kontak.Alamat)
		if err != nil {
			panic(err.Error())
		}
		contacts = append(contacts, kontak)
	}
	json.NewEncoder(w).Encode(contacts)
}

func createContact(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ID := r.FormValue("ID")
		NamaDepan := r.FormValue("NamaDepan")
		NamaBelakang := r.FormValue("NamaBelakang")
		NoHp := r.FormValue("NoHp")
		Email := r.FormValue("Email")
		Alamat := r.FormValue("Alamat")
		stmt, err := db.Prepare("INSERT INTO kontak (namadepan, namabelakang, nohp, email, alamat) VALUES (?,?,?,?,?)")
		_, err = stmt.Exec(ID, NamaDepan, NamaBelakang, NoHp, Email, Alamat)
		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}
	}
}

func getContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contacts []Kontak
	params := mux.Vars(r)
	sql := `SELECT
				ID,
				IFNULL(namadepan,'') namadepan,
				IFNULL(namabelakang,'') namabelakang,
				IFNULL(nohp,'') nohp,
				IFNULL(email,'') email,
				IFNULL(alamat,'') alamat,
			FROM kontak WHERE ID = ?`
	result, err := db.Query(sql, params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var kontak Kontak
	for result.Next() {
		err := result.Scan(&kontak.ID, &kontak.NamaDepan, &kontak.NamaBelakang,
			&kontak.NoHp, &kontak.Email, &kontak.Alamat)
		if err != nil {
			panic(err.Error())
		}
		contacts = append(contacts, kontak)
	}
	json.NewEncoder(w).Encode(contacts)
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		params := mux.Vars(r)
		newNamaDepan := r.FormValue("NamaDepan")
		newNamaBelakang := r.FormValue("NamaBelakang")
		newNoHp := r.FormValue("NoHp")
		newEmail := r.FormValue("Email")
		newAlamat := r.FormValue("Alamat")
		stmt, err := db.Prepare("UPDATE kontak SET namadepan = ?, namabelakang = ?, nohp = ?, email = ?, alamat = ? WHERE ID = ?")
		_, err = stmt.Exec(newNamaDepan, newNamaBelakang, newNoHp, newEmail, newAlamat, params["id"])
		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}
		fmt.Fprintf(w, "Contact with ID = %s was updated", params["id"])
	}
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM kontak WHERE ID = ?")
	_, err = stmt.Exec(params["id"])
	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}
	fmt.Fprintf(w, "Contact with ID = %s was deleted", params["id"])
}

// Main function
func main() {
	db, err = sql.Open("mysql", "root:0000@tcp(127.0.0.1:3306)/challenge")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Init router
	r := mux.NewRouter()
	// Route handles & endpoints
	r.HandleFunc("/contacts", getContacts).Methods("GET")
	r.HandleFunc("/contacts/{id}", getContact).Methods("GET")
	r.HandleFunc("/contacts", createContact).Methods("POST")
	r.HandleFunc("/contacts/{id}", updateContact).Methods("PUT")
	r.HandleFunc("/contacts/{id}", deleteContact).Methods("DELETE")
	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
