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
	router.HandleFunc("/api/v1/carowner/{id}", carowner).Methods("GET", "DELETE", "POST", "PATCH", "PUT", "OPTIONS")
	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))
}

// update user to carowner
func createCarOwner(co CarOwner) bool {
	fmt.Println("in createCarOwner func")
	openDB()
	defer db.Close()
	_, err := db.Exec("insert into CarOwner values(?,?,?)", co.UserAcc.ID, co.LicenseNo, co.PlateNo)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// check carowner exist
func checkCarOwnerID(id string) (CarOwner, bool) {
	fmt.Println("in carOwnerID func")
	openDB()
	defer db.Close()
	var co CarOwner
	query := "select * from CarOwner where UserID = \"" + id + "\""
	row := db.QueryRow(query)
	fmt.Println(query)
	err := row.Scan(&co.UserAcc.ID, &co.LicenseNo, &co.PlateNo)
	if err == sql.ErrNoRows { // if doesn't exist return false
		fmt.Println(err)
		return co, false
	}
	return co, true
}

// update carowner details
func updateCarOwner(user CarOwner) bool {
	fmt.Println("update carowner query")
	openDB()
	defer db.Close() //will run at the end of the block of the code
	_, err := db.Exec("update CarOwner set LicenseNo=?, PlateNo=? where UserID=?;", user.LicenseNo, user.PlateNo, user.UserAcc.ID)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// delete acc for car owner
func deleteCarOwner(id string) bool {
	fmt.Println("in delete carowner function")
	openDB()
	defer db.Close() //will run at the end of the block of the code
	_, err := db.Exec("delete from CarOwner where UserID = ?", id)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// API for car owner
func carowner(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["id"]
	fmt.Println("carowner func..")

	switch r.Method {
	case http.MethodPost: //update to car owner
		if param == "newCarOwner" {
			if body, err := io.ReadAll(r.Body); err == nil {
				var clientRes CarOwner
				fmt.Println(string(body))
				if err := json.Unmarshal(body, &clientRes); err == nil {
					if created := createCarOwner(clientRes); created {
						w.WriteHeader(http.StatusAccepted) //202
					} else {
						w.WriteHeader(http.StatusConflict)
					}
				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		} else { //check usertype
			if carowner, exist := checkCarOwnerID(param); exist {
				w.WriteHeader(http.StatusOK) //202
				userJSON, _ := json.Marshal(carowner)
				fmt.Println(carowner)
				w.Write(userJSON)
			} else {
				w.WriteHeader(http.StatusUnauthorized) //401
			}
		}

	case http.MethodPut: //update account information
		if body, err := io.ReadAll(r.Body); err == nil {
			fmt.Println("methodput")
			var user CarOwner
			if err := json.Unmarshal(body, &user); err == nil {
				json.NewDecoder(r.Body).Decode(&user)
				if status := updateCarOwner(user); status {
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
		if status := deleteCarOwner(param); status {
			w.WriteHeader(http.StatusAccepted) //202
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	}
}
