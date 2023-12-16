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
	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))
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

// get the trip details (of a carowner)
func getTripPassengers(query string) ([]Trip, error) {
	fmt.Println("in getTripPassengers func")
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
		if err := rows.Scan(&trip.ID, &trip.Driver.UserAcc.ID, &trip.PickUpAddress, &trip.AlternativeAddress, &startTimeStr, &trip.Destination, &trip.PaxLimit); err != nil {
			fmt.Println(err)
			return nil, err
		}
		trip.StartTime, _ = time.Parse("2006-01-02 15:04:05", startTimeStr)
		trips = append(trips, trip)
	}
	return trips, nil
}

// cancel a trip
func deleteTrip(id string) bool {
	fmt.Println("in delete trip function")
	openDB()
	defer db.Close() //will run at the end of the block of the code
	_, err := db.Exec("delete from Trip where TripID = ?", id)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// REST API
func trip(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["id"]
	switch r.Method {
	case http.MethodPost: //create trip
		if body, err := io.ReadAll(r.Body); err == nil {
			var trip Trip
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &trip); err == nil {
				if created := createTrip(trip, param); created {
					w.WriteHeader(http.StatusAccepted) //202
					fmt.Println(trip)
				} else {
					w.WriteHeader(http.StatusConflict)
				}
			}
		}
	case http.MethodGet: //get trip by ownerID
		if tripDetails, err := getTripPassengers(param); err == nil {
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
	case http.MethodDelete: //cancel trip
		if status := deleteTrip(param); status {
			w.WriteHeader(http.StatusAccepted) //202
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	}

}
