package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// Construct
type Customer struct {
	Address struct {
		City   string `json:"city"`
		State  string `json:"state"`
		Street string `json:"street"`
		Zip    string `json:"zip"`
	} `json:"address"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Employee struct {
	Address struct {
		City   string `json:"city"`
		State  string `json:"state"`
		Street string `json:"street"`
		Zip    string `json:"zip"`
	} `json:"address"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

func getCustomers(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var request Customer

	if err = json.Unmarshal(body, &request); err != nil {
		fmt.Println("Failed decoding json message")
	}

	fmt.Fprintln(w, "First Name : "+request.FirstName)
	fmt.Fprintln(w, "City Name : "+request.Address.City)

	//Tugas insert kan ke table Customer
}

func insertEmployees(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var request Employee

	if err = json.Unmarshal(body, &request); err != nil {
		fmt.Println("Failed decoding json message")
	} else {
		if r.Method == "POST" {

			FirstName := request.FirstName
			LastName := request.LastName

			stmt, err := db.Prepare("INSERT INTO employees (LastName, FirstName) VALUES (?,?)")

			_, err = stmt.Exec(LastName, FirstName)

			if err != nil {
				fmt.Fprintf(w, "Data Duplicate")
			} else {
				fmt.Fprintf(w, "Data Created")
			}

		}
	}

}

func insCustomers(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var request Root

	if err = xml.Unmarshal(body, &request); err != nil {
		fmt.Println("Failed decoding json message")
	}

	for i := 0; i < len(request.Customers.Customer); i++ {
		if r.Method == "POST" {

			CustomerID := request.Customers.Customer[i].CustomerID
			CompanyName := request.Customers.Customer[i].CompanyName
			stmt, err := db.Prepare("INSERT INTO customers (CustomerID,CompanyName) VALUES (?,?)")

			_, err = stmt.Exec(CustomerID, CompanyName)

			if err != nil {
				fmt.Fprintf(w, "Data Duplicate")
			} else {
				fmt.Fprintf(w, "Data Created")
			}
			fmt.Fprintln(w, "Customer ID: "+request.Customers.Customer[i].CustomerID)
			fmt.Fprintln(w, "Company Name : "+request.Customers.Customer[i].CompanyName)
		}

	}
}
func insOrders(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var request Root

	if err = xml.Unmarshal(body, &request); err != nil {
		fmt.Println("Failed decoding json message")
	}

	for i := 0; i < len(request.Orders.Order); i++ {
		if r.Method == "POST" {

			CustomerID := request.Orders.Order[i].CustomerID
			EmployeeID := request.Orders.Order[i].EmployeeID
			stmt, err := db.Prepare("INSERT INTO orders (CustomerID,EmployeeID) VALUES (?,?)")

			_, err = stmt.Exec(CustomerID, EmployeeID)

			if err != nil {
				fmt.Fprintf(w, "Data Duplicate")
			} else {
				fmt.Fprintf(w, "Data Created")
			}
			fmt.Fprintln(w, "CustomerID: "+request.Orders.Order[i].CustomerID)
			fmt.Fprintln(w, "EmployeeID : "+request.Orders.Order[i].EmployeeID)
		}

	}
}

func main() {

	db, err = sql.Open("mysql", "root:0000@tcp(127.0.0.1:3306)/northwind")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	fmt.Println("Server on :8181")

	// Route handles & endpoints
	r.HandleFunc("/customers", getCustomers).Methods("POST")
	r.HandleFunc("/employee", insertEmployees).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":8181", r))

}
