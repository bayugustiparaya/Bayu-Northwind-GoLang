package main

import (
	"database/sql"
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

type Root struct {
	XMLName   xml.Name `xml:"Root"`
	Text      string   `xml:",chardata"`
	Xmlns     string   `xml:"xmlns,attr"`
	Customers struct {
		Text     string `xml:",chardata"`
		Customer []struct {
			Text         string `xml:",chardata"`
			CustomerID   string `xml:"CustomerID,attr"`
			CompanyName  string `xml:"CompanyName"`
			ContactName  string `xml:"ContactName"`
			ContactTitle string `xml:"ContactTitle"`
			Phone        string `xml:"Phone"`
			FullAddress  struct {
				Text       string `xml:",chardata"`
				Address    string `xml:"Address"`
				City       string `xml:"City"`
				Region     string `xml:"Region"`
				PostalCode string `xml:"PostalCode"`
				Country    string `xml:"Country"`
			} `xml:"FullAddress"`
			Fax string `xml:"Fax"`
		} `xml:"Customer"`
	} `xml:"Customers"`
	Orders struct {
		Text  string `xml:",chardata"`
		Order []struct {
			Text         string `xml:",chardata"`
			CustomerID   string `xml:"CustomerID"`
			EmployeeID   string `xml:"EmployeeID"`
			OrderDate    string `xml:"OrderDate"`
			RequiredDate string `xml:"RequiredDate"`
			ShipInfo     struct {
				Text           string `xml:",chardata"`
				ShippedDate    string `xml:"ShippedDate,attr"`
				ShipVia        string `xml:"ShipVia"`
				Freight        string `xml:"Freight"`
				ShipName       string `xml:"ShipName"`
				ShipAddress    string `xml:"ShipAddress"`
				ShipCity       string `xml:"ShipCity"`
				ShipRegion     string `xml:"ShipRegion"`
				ShipPostalCode string `xml:"ShipPostalCode"`
				ShipCountry    string `xml:"ShipCountry"`
			} `xml:"ShipInfo"`
		} `xml:"Order"`
	} `xml:"Orders"`
}

func getCustomers(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var request Root

	if err = xml.Unmarshal(body, &request); err != nil {
		fmt.Println("Failed decoding json message")
	}

	for i := 0; i < len(request.Customers.Customer); i++ {

		fmt.Fprintln(w, "Customer ID: "+request.Customers.Customer[i].CustomerID)
		fmt.Fprintln(w, "Company Name : "+request.Customers.Customer[i].CompanyName)
	}
}

//Tugas insert kan ke table Customer
func insCustomer(w http.ResponseWriter, r *http.Request) {

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
			// fmt.Fprintln(w, "Customer ID: "+request.Customers.Customer[i].CustomerID)
			// fmt.Fprintln(w, "Company Name : "+request.Customers.Customer[i].CompanyName)
		}

	}
}

//Tugas insert kan ke table Order
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
			// fmt.Fprintln(w, "CustomerID: "+request.Orders.Order[i].CustomerID)
			// fmt.Fprintln(w, "EmployeeID : "+request.Orders.Order[i].EmployeeID)
		}

	}
}

func getOrders(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var request Root

	if err = xml.Unmarshal(body, &request); err != nil {
		fmt.Println("Failed decoding json message")
	}

	for i := 0; i < len(request.Orders.Order); i++ {

		fmt.Fprintln(w, "Customer ID: "+request.Orders.Order[i].CustomerID)
		fmt.Fprintln(w, "Employee ID: "+request.Orders.Order[i].EmployeeID)
		fmt.Fprintln(w, "Order Date: "+request.Orders.Order[i].OrderDate)
		fmt.Fprintln(w, "Require Date: "+request.Orders.Order[i].RequiredDate)
		// fmt.Fprintln(w, "Ship Info: "+request.Orders.Order[i].ShipInfo.Text)
		fmt.Fprintln(w, "==========================")

	}
	//Tugas insert kan ke table Customer

	//Tugas insert kan ke table Order

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
	r.HandleFunc("/orders", getOrders).Methods("POST")
	r.HandleFunc("/insertcustomers", insCustomer).Methods("POST")
	r.HandleFunc("/insertorders", insOrders).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":8181", r))

}
