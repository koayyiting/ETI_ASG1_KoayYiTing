package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Passenger struct {
	ID         string
	PhoneNo    string
	Email      string
	FirstName  string
	LastName   string
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

var (
	db  *sql.DB
	err error
)

func openDB() {
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/version1")

	if err != nil {
		panic(err.Error())
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/passenger/{id}", passenger).Methods("GET", "DELETE", "POST", "PATCH", "PUT", "OPTIONS")
	router.HandleFunc("/api/v1/passenger/passengertrip/{id}", passengertrip).Methods("GET", "POST", "OPTIONS")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

// API for enrolling the passengers in a trip
func passengertrip(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["id"]
	switch r.Method {
	case http.MethodPost:
		if body, err := io.ReadAll(r.Body); err == nil {
			var p Passenger
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &p); err == nil {
				enrolTrip(param, p)
				w.WriteHeader(http.StatusAccepted) //202
			}
		}
	case http.MethodGet:
		if tripDetails, err := getPassengerTrip(param); err == nil {
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

// get all the passenger's trip
func getPassengerTrip(pid string) ([]Trip, error) {
	fmt.Println("in getTripPassengers func")
	openDB()
	defer db.Close()
	rows, err := db.Query("SELECT T.* FROM Trip T JOIN TripPassenger TP ON T.TripID = TP.TripID WHERE TP.PassengerID = ?", pid)
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

// enrol the passenger into a trip
func enrolTrip(id string, p Passenger) bool {
	fmt.Println("in enrolTrip function")
	openDB()
	defer db.Close()

	// Add passengers to the trip
	_, err := db.Exec("INSERT INTO TripPassenger (TripID, PassengerID) VALUES (?, ?)", id, p.ID)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

// API for passenger actions
func passenger(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["id"]

	switch r.Method {
	case http.MethodPost: //create new passenger
		if param == "newUser" {
			if body, err := io.ReadAll(r.Body); err == nil {
				var newUser Passenger
				fmt.Println(string(body))
				if err := json.Unmarshal(body, &newUser); err == nil {
					if userID, exist := generateID(); !exist {
						createNewUser(newUser, userID)
						w.WriteHeader(http.StatusAccepted) //202
					}
				}
			}
		} else { // to retrieve user details
			if body, err := io.ReadAll(r.Body); err == nil {
				var user Passenger
				fmt.Println(string(body))
				if err := json.Unmarshal(body, &user); err == nil {
					if userDetailsDB, exist := userExist_Login(param); exist {
						fmt.Println("User exist")
						w.WriteHeader(http.StatusAccepted) //202
						userJSON, _ := json.Marshal(userDetailsDB)
						fmt.Println(userDetailsDB)
						w.Write(userJSON)
					} else {
						w.WriteHeader(http.StatusUnauthorized) //401
					}
				}
			}
		}

	case http.MethodPut: //update account information
		fmt.Println("method put:")
		if body, err := io.ReadAll(r.Body); err == nil {
			var user Passenger
			if err := json.Unmarshal(body, &user); err == nil {
				json.NewDecoder(r.Body).Decode(&user)
				if status := updatePassenger(user); status {
					w.WriteHeader(http.StatusAccepted) //202
					userJSON, _ := json.Marshal(user)
					w.Write(userJSON)
				} else {
					w.WriteHeader(http.StatusConflict)
				}

			} else {
				fmt.Println(err)
			}
		}
	case http.MethodDelete: //delete acc after 1 year
		if status := deletePassenger(param); status {
			w.WriteHeader(http.StatusAccepted) //202
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	}

}

// delete passenger
func deletePassenger(id string) bool {
	fmt.Println("in delete passenger function")
	openDB()
	defer db.Close() //will run at the end of the block of the code
	_, err := db.Exec("delete from Passenger where ID = ?", id)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// update passenger details
func updatePassenger(user Passenger) bool {
	fmt.Println("update passenger query")
	openDB()
	defer db.Close() //will run at the end of the block of the code
	_, err := db.Exec("update Passenger set FirstName=?, LastName=?, Email=?, PhoneNumber=? where ID=?;", user.FirstName, user.LastName, user.Email, user.PhoneNo, user.ID)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// check user exist for login
func userExist_Login(param string) (Passenger, bool) {
	fmt.Println("checking user exist..")
	openDB()
	var p Passenger
	fmt.Println(param)
	info := strings.Split(param, "_")
	detail := info[0]
	column := info[1]
	query := "select * from Passenger where " + column + " = \"" + detail + "\""
	row := db.QueryRow(query)
	fmt.Println(query)
	var RegisterDT_str string
	err := row.Scan(&p.ID, &p.PhoneNo, &p.Email, &p.FirstName, &p.LastName, &RegisterDT_str)
	p.RegisterDT, _ = time.Parse("2006-01-02 15:04:05", RegisterDT_str)
	if err == sql.ErrNoRows { // if doesn't exist return false
		fmt.Println("user doesnt exist")
		return p, false
	} else if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	return p, true //exist return true
}

// create new user
func createNewUser(p Passenger, id string) {
	fmt.Println("in createNewUser func")
	openDB()
	fmt.Println(p.PhoneNo)
	_, err := db.Exec("insert into Passenger values(?,?,?,?,?,?)", id, p.PhoneNo, p.Email, p.FirstName, p.LastName, p.RegisterDT)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
}

// generate id
func generateID() (string, bool) {
	fmt.Println("generating ID and checking..")
	unique := false

	openDB()
	var p Passenger
	for !unique {
		// generate unique ID
		id := uuid.New()
		idString := id.String() // Convert the UUID to a string

		// Check in SQL
		row := db.QueryRow("select * from Passenger where ID=?", idString)
		err := row.Scan(&p.ID, &p.PhoneNo, &p.Email, &p.FirstName, &p.LastName)
		if err == sql.ErrNoRows { // if [doesn't exist] return false
			unique = true
			return idString, false
		}
	}

	defer db.Close()
	return "", true //return exist [never going to happen]
}
