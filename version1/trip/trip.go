package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Passenger struct {
	ID         string `json:"ID"`
	PhoneNo    string `json:"PhoneNo"`
	Email      string `json:"Email"`
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	RegisterDT time.Time
}

type CarOwner struct {
	UserAcc   Passenger
	LicenseNo string
	PlateNo   string
}

type Trip struct {
	ID                 string
	PickUpAddress      string
	AlternativeAddress string
	StartTime          time.Time
	Destination        string
	PaxLimit           int
	Status             string
	Driver             CarOwner
}

type TripPassenger struct {
	TripID      string
	PassengerID string
}

var (
	db  *sql.DB
	err error
)

func generateTripID() string {
	id := uuid.New()
	return id.String()
}

func openDB() {
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/version1")

	if err != nil {
		panic(err.Error())
	}
}

func main() {
	// openDB()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/trip/{id}", trip).Methods("GET", "DELETE", "POST", "PATCH", "PUT", "OPTIONS")
	router.HandleFunc("/api/v1/tripPassenger/{id}", tripPassenger).Methods("GET", "DELETE", "POST", "PATCH", "PUT", "OPTIONS")
	router.HandleFunc("/api/v1/listOfTrip/{query}", listOfTrip).Methods("GET", "DELETE", "POST", "PATCH", "PUT", "OPTIONS")
	router.HandleFunc("/api/v1/alltrip/{searchType}/{search}", alltrip).Methods("GET", "DELETE", "POST", "PATCH", "PUT", "OPTIONS")
	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))
}

// REST API - /api/v1/trip/{id}
func trip(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Println(r.Method)
	switch r.Method {
	case http.MethodPost: //create trip
		if body, err := io.ReadAll(r.Body); err == nil {
			var trip Trip
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &trip); err == nil {
				if created := createTrip(trip, id); created {
					w.WriteHeader(http.StatusAccepted) //202
					fmt.Println(trip)
				} else {
					w.WriteHeader(http.StatusConflict)
				}
			}
		}
	case http.MethodPut: //update trip status
		if body, err := io.ReadAll(r.Body); err == nil {
			var trip Trip
			if err := json.Unmarshal(body, &trip); err == nil {
				json.NewDecoder(r.Body).Decode(&trip)
				if status := updateTrip(trip); status {
					w.WriteHeader(http.StatusAccepted) //202
					tripJSON, _ := json.Marshal(trip)
					w.Write(tripJSON)
				} else {
					w.WriteHeader(http.StatusConflict)
				}
			} else {
				fmt.Println(err)
			}
		}
	case http.MethodDelete: //cancel trip
		if status := deleteTrip(id); status { //param = userid
			w.WriteHeader(http.StatusAccepted) //202
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	}

}

// publish a trip
func createTrip(trip Trip, userid string) bool {
	fmt.Println("in createTrip func")
	openDB()
	defer db.Close()
	trip.ID = generateTripID()
	_, err := db.Exec("insert into Trip (TripID, CarOwnerID, PickUpAddress, AlternativeAddress, StartTime, Destination, PaxLimit) values(?,?,?,?,?,?,?)", trip.ID, userid, trip.PickUpAddress, trip.AlternativeAddress, trip.StartTime, trip.Destination, trip.PaxLimit)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// update trip status
func updateTrip(trip Trip) bool {
	fmt.Println("update trip query")
	openDB()
	defer db.Close() //will run at the end of the block of the code
	rowStatus := db.QueryRow("select Status from Trip where TripID=?", trip.ID)
	var status string
	rowStatus.Scan(&status)
	if status == "Cancelled" {
		fmt.Println("Status Already Cancelled")
		return false
	}
	_, err := db.Exec("update Trip set Status=? where TripID=?", trip.Status, trip.ID)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// cancel a trip (by a carowner id)
func deleteTrip(id string) bool {
	fmt.Println("in delete trip function")
	openDB()
	defer db.Close() //will run at the end of the block of the code
	_, err := db.Exec("delete from Trip where CarOwnerID = ?", id)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// REST API - /api/v1/listOfTrip/{query}
func tripPassenger(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	switch r.Method {
	case http.MethodDelete:
		fmt.Println("in")
		query := "SELECT * FROM Trip WHERE CarOwnerID = \"" + id + "\";"
		trips, _ := getTrips(query)
		for _, trip := range trips {
			if err := deleteTripPassenger(trip.ID); err != nil { //param = userid
				w.WriteHeader(http.StatusConflict) //202
				break
			}
			w.WriteHeader(http.StatusAccepted) //202
		}
	}
}

// cancel a trip (by a carowner id)
func deleteTripPassenger(id string) error {
	fmt.Println("in delete trip passenger function")
	openDB()
	defer db.Close() //will run at the end of the block of the code
	_, err := db.Exec("delete from TripPassenger where TripID = ?", id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// REST API - /api/v1/listOfTrip/{query}
func listOfTrip(w http.ResponseWriter, r *http.Request) {
	query := mux.Vars(r)["query"]

	switch r.Method {
	case http.MethodGet: //get trip by ownerID
		if tripDetails, err := getTrips(query); err == nil {
			w.WriteHeader(http.StatusAccepted) //202
			if tripJSON, err := json.Marshal(tripDetails); err == nil {
				fmt.Println(tripDetails)
				w.Write(tripJSON)
			} else {
				fmt.Println(err)
			}
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	}
}

// get the trip details (of a carowner)
func getTrips(query string) ([]Trip, error) {
	fmt.Println("in getTrips func")
	fmt.Println(query)
	openDB()
	defer db.Close()
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var trips []Trip
	var startTimeStr string
	for rows.Next() {
		var trip Trip
		var status sql.NullString
		if err := rows.Scan(&trip.ID, &trip.Driver.UserAcc.ID, &trip.PickUpAddress, &trip.AlternativeAddress, &startTimeStr, &trip.Destination, &trip.PaxLimit, &status); err != nil {
			fmt.Println(err)
			return nil, err
		}
		if status.Valid {
			trip.Status = status.String
		} else {
			trip.Status = "-"
		}
		trip.StartTime, _ = time.Parse("2006-01-02 15:04:05", startTimeStr)
		trips = append(trips, trip)
	}
	return trips, nil
}

// REST API - /api/v1/alltrip/{searchType}/{search}
func alltrip(w http.ResponseWriter, r *http.Request) {
	searchType := mux.Vars(r)["searchType"]
	searchInput := mux.Vars(r)["search"]
	fmt.Println(r.Method)
	switch r.Method {
	case http.MethodGet: //get all trips for user
		if searchType == "listall" {
			searchType = "PickUpAddress"
			searchInput = ""
		}
		if tripDetails, err := search(searchType, searchInput); err == nil {
			w.WriteHeader(http.StatusAccepted) //202
			if tripJSON, err := json.Marshal(tripDetails); err == nil {
				fmt.Println(tripDetails)
				w.Write(tripJSON)
			} else {
				fmt.Println(err)
			}
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	}

}

func search(searchType string, searchInput string) ([]Trip, error) {
	fmt.Println("search query")
	openDB()
	defer db.Close()
	query := "SELECT * FROM Trip WHERE PaxLimit > (SELECT COUNT(*) FROM TripPassenger WHERE TripPassenger.TripID = Trip.TripID) AND NOW() < StartTime AND " + searchType + " LIKE ? AND (Status IS NULL OR Status = \"Start\" OR Status != \"Cancelled\")"
	rows, err := db.Query(query, "%"+searchInput+"%")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var trips []Trip
	var startTimeStr string
	for rows.Next() {
		var trip Trip
		var status sql.NullString
		if err := rows.Scan(&trip.ID, &trip.Driver.UserAcc.ID, &trip.PickUpAddress, &trip.AlternativeAddress, &startTimeStr, &trip.Destination, &trip.PaxLimit, &status); err != nil {
			fmt.Println(err)
			return nil, err
		}
		if status.Valid {
			trip.Status = status.String
		} else {
			trip.Status = "-"
		}
		trip.StartTime, _ = time.Parse("2006-01-02 15:04:05", startTimeStr)
		trips = append(trips, trip)
	}
	return trips, nil
}
