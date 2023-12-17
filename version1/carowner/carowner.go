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
	router.HandleFunc("/api/v1/carowner/{param}", carowner).Methods("GET", "DELETE", "POST", "PATCH", "PUT", "OPTIONS")
	router.HandleFunc("/api/v1/carowner/{loginType}/{loginDetail}", carownerLogin).Methods("GET", "DELETE", "POST", "PATCH", "PUT", "OPTIONS")
	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))
}

// API for car owner - /api/v1/carowner/{param}
func carowner(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["param"]
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
			if carowner, exist := checkCarOwnerID(param); exist { //param == id
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
		if status := deleteCarOwner(param); status { //param == id
			w.WriteHeader(http.StatusAccepted) //202
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	}
}

// update user to carowner MethodPost
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

// check carowner exist MethodPost
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

// update carowner details MethodPut
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

// delete acc for car owner MethodDelete
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

// REST API login carowner - /api/v1/carowner/{loginType}/{loginDetail}
func carownerLogin(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	loginType := param["loginType"]
	loginDetail := param["loginDetail"]

	switch r.Method {
	case http.MethodPost: //create new passenger
		if body, err := io.ReadAll(r.Body); err == nil {
			var user Passenger
			fmt.Println(string(body))
			if err := json.Unmarshal(body, &user); err == nil {
				if userDetailsDB, exist := userExist_Login(loginType, loginDetail); exist {
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
}

func userExist_Login(loginType string, loginDetail string) (CarOwner, bool) {
	fmt.Println("checking user exist..")
	openDB()
	var co CarOwner
	query := "SELECT CarOwner.* FROM CarOwner JOIN Passenger ON CarOwner.UserID = Passenger.ID WHERE Passenger." + loginType + " = " + loginDetail
	row := db.QueryRow(query)
	fmt.Println(query)
	err := row.Scan(&co.UserAcc.ID, &co.LicenseNo, &co.PlateNo)
	if err == sql.ErrNoRows { // if doesn't exist return false
		fmt.Println("user doesnt exist")
		return co, false
	} else if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	var RegisterDT_str string
	query2 := "SELECT * FROM Passenger WHERE ID = ?;"
	row2 := db.QueryRow(query2, co.UserAcc.ID)
	err2 := row2.Scan(&co.UserAcc.ID, &co.UserAcc.PhoneNo, &co.UserAcc.Email, &co.UserAcc.FirstName, &co.UserAcc.LastName, &RegisterDT_str)
	co.UserAcc.RegisterDT, _ = time.Parse("2006-01-02 15:04:05", RegisterDT_str)
	if err2 == sql.ErrNoRows { // if doesn't exist return false
		fmt.Println("user doesnt exist")
		return co, false
	} else if err2 != nil {
		fmt.Println(err)
	}
	defer db.Close()
	return co, true //exist return true
}
