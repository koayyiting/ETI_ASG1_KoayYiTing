package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	currentPassenger Passenger
	currentCarOwner  CarOwner
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

// Menus
func mainMenu() {
	fmt.Println("===================")
	fmt.Println("Welcome to Car Pool")
	fmt.Println("1. Log In")
	fmt.Println("2. Sign Up")
	fmt.Println("0. Exit")
	fmt.Println("===================")
	fmt.Print("Enter an Option: ")
}

func carPoolMenu(userType bool) {
	var uType string
	if userType {
		uType = "Car Owner"
	} else {
		uType = "Passenger"
	}
	welcome_message := "Hi " + currentPassenger.FirstName + "! Welcome to Car Pool Main Page [" + uType + "]"
	fmt.Println(strings.Repeat("=", len(welcome_message)))
	fmt.Println(welcome_message)
	fmt.Println("1. Update User Details")
	fmt.Println("2. Delete Account")
	if userType { //for car owners
		fmt.Println("3. Publish Trip")
		fmt.Println("4. List of Published Trip")
		fmt.Println("5. Start/Cancel Published Trip")
		fmt.Println("6. View Account Detail")
	} else { //for passenger
		fmt.Println("3. Change to Car Owner")
		fmt.Println("4. Enrol Trip")
		fmt.Println("5. View Account Detail")
		fmt.Println("6. List of Enrolled Trips")
		fmt.Println("7. Search for Trips")
		fmt.Println("8. List all Trips")
	}
	fmt.Println("0. Log Out")
	fmt.Println(strings.Repeat("=", len(welcome_message)))
	fmt.Print("Enter an Option: ")
}

func (p Passenger) String() string {
	return fmt.Sprintf("PhoneNo: %s\nEmail: %s\nFirstName: %s\nLastName: %s\nRegisterDT: %s",
		p.PhoneNo, p.Email, p.FirstName, p.LastName, p.RegisterDT)
}

func (co CarOwner) String() string {
	return fmt.Sprintf("%s\nLicenseNo: %s\nPlateNo: %s", co.UserAcc.String(), co.LicenseNo, co.PlateNo)
}

// login/signup page
func main() {
	loop := true
	for loop {
		mainMenu()

		curl_login := "http://localhost:5000/api/v1/loginpassenger/"
		curl_newUser := "http://localhost:5000/api/v1/newpassenger/"
		var choice int
		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			userExist := login(curl_login)
			if userExist {
				if exit := mainpage(); exit {
					break
				}
			}
		case 2:
			signup(curl_newUser)
		case 0:
			return
		}
	}
}

func mainpage() bool {
	loop := true
	curl_passenger := "http://localhost:5000/api/v1/passenger/"
	curl_carowner := "http://localhost:5001/api/v1/carowner/"
	curl_trip := "http://localhost:5002/api/v1/trip/"
	curl_alltrip := "http://localhost:5002/api/v1/alltrip/"
	curl_listOfTrip := "http://localhost:5002/api/v1/listOfTrip/"
	reader := bufio.NewReader(os.Stdin)
	for loop {
		userType := checkUserType(currentPassenger.ID, curl_carowner) //carowner ==true
		carPoolMenu(userType)
		var choice int
		fmt.Scanf("%d", &choice)
		reader.ReadString('\n')

		switch userType {
		case true: //car owners
			switch choice {
			case 1: //update acc details
				updateUserInfo(curl_passenger, userType)
			case 2: //delete acc
				if lessThanOneYear := checkOneYear(); !lessThanOneYear {
					if deleted := deleteAcc(curl_passenger, userType); deleted {
						return true
					}
				} else {
					fmt.Println("Users are only able to delete account after 1 year")
				}
			case 3: //publish trip
				createTrip(curl_trip)
			case 4: //list trips
				trips := getTripListCarOwner(curl_listOfTrip, false, "")
				tripList(trips, "Available Trips", false)
			case 5: //start/cancel trip
				fmt.Print("[1]start [2]cancel [0]exit: ")
				var tripChoice int
				fmt.Scanf("%d", &tripChoice)
				reader.ReadString('\n')
				if tripChoice == 1 {
					trips := getTripListCarOwner(curl_listOfTrip, true, "start")
					if len(trips) > 0 {
						tripList(trips, "Your Published Trip", true)
						updateTripStatus(curl_trip, trips, "Start")
					} else {
						fmt.Println("No trips to be started")
					}

				} else if tripChoice == 2 {
					trips := getTripListCarOwner(curl_listOfTrip, true, "cancel")
					if len(trips) > 0 {
						tripList(trips, "Your Published Trip", true)
						updateTripStatus(curl_trip, trips, "Cancelled")
					} else {
						fmt.Println("No trips to be cancelled")
					}

				} else if tripChoice == 0 {
					fmt.Println("back")
				}
			case 6: //view acc detail
				fmt.Println(currentCarOwner.String())
			case 0:
				var empty CarOwner
				currentCarOwner = empty
				return true
			}
		case false: //passengers
			switch choice {
			case 1:
				updateUserInfo(curl_passenger, userType)
			case 2:
				if lessThanOneYear := checkOneYear(); !lessThanOneYear {
					if deleted := deleteAcc(curl_passenger, userType); deleted {
						return true
					}
				} else {
					fmt.Println("Users are only able to delete account after 1 year")
				}
			case 3:
				updateToCarOwner(curl_carowner)
			case 4: // enrol trips
				trips := getTripList(curl_listOfTrip)
				if len(trips) > 0 {
					tripList(trips, "Available Trip(s)", true)
					enrolTrip(curl_passenger+"passengertrip/", trips)
				} else {
					fmt.Println("No Trips Available")
				}
			case 5:
				fmt.Println(currentPassenger.String())
			case 6: //list user enrolled trip
				listUserEnrolTrip(curl_passenger + "passengertrip/")
			case 7: //search
				search(curl_alltrip)
			case 8:
				listAllTrip(curl_alltrip)
			case 0:
				var empty Passenger
				currentPassenger = empty
				return true
			}
		}
	}
	return true
}

func listAllTrip(url string) []Trip {
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodGet, url+"listall/nil", nil); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				if body, err := io.ReadAll(res.Body); err == nil {
					var trips []Trip
					json.Unmarshal(body, &trips)
					if len(trips) != 0 {
						tripList(trips, "List of Trip(s)", false)
					} else {
						fmt.Println("No Trip(s) Available")
					}
				}
			} else {
				fmt.Println(res.StatusCode)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return nil
}

func search(url string) []Trip {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("[1]Pick-Up Address [2]Destination: ")
	var choice int
	fmt.Scanf("%d", &choice)
	reader.ReadString('\n')
	var param string
	var paramInput string
	if choice == 1 {
		fmt.Print("Pick-Up Address: ")
		pua, _ := reader.ReadString('\n')
		paramInput = strings.TrimSpace(pua)
		param = "PickUpAddress/"
	} else if choice == 2 {
		fmt.Print("Destination: ")
		d, _ := reader.ReadString('\n')
		paramInput = strings.TrimSpace(d)
		param = "Destination/"
	} else {
		return nil
	}
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodGet, url+param+paramInput, nil); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				if body, err := io.ReadAll(res.Body); err == nil {
					var trips []Trip
					json.Unmarshal(body, &trips)
					if len(trips) != 0 {
						tripList(trips, "Available Searched Trip", true)
					} else {
						fmt.Println("0 Searched")
					}
				}
			} else {
				fmt.Println(res.StatusCode)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return nil
}

func checkOneYear() bool {
	age := time.Since(currentPassenger.RegisterDT)
	return age < time.Hour*24*365
}

func listUserEnrolTrip(url string) []Trip {
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodGet, url+currentPassenger.ID, nil); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				if body, err := io.ReadAll(res.Body); err == nil {
					var trips []Trip
					json.Unmarshal(body, &trips)
					if len(trips) != 0 {
						tripList(trips, "Enrolled Trip(s)", false)
					} else {
						fmt.Println("No Trip(s) Enrolled")
					}

				}
			} else {
				fmt.Println(res.StatusCode)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return nil
}

func enrolTrip(url string, trips []Trip) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a Trip to enrol: ")
	var choice int
	fmt.Scanf("%d", &choice)
	reader.ReadString('\n')

	if choice < 0 || choice > len(trips) {
		fmt.Println("Invalid Trip Option")
		return
	} else if choice == 0 {
		return
	}
	client := &http.Client{}
	tripid := trips[choice-1].ID
	postBody, _ := json.Marshal(currentPassenger)
	resBody := bytes.NewBuffer(postBody)
	if req, err := http.NewRequest(http.MethodPost, url+tripid, resBody); err == nil {
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				fmt.Println("Trip: ", tripid, " enrolled successfully")
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func getTripList(url string) []Trip {
	client := &http.Client{}
	query := "SELECT * FROM Trip WHERE PaxLimit > (SELECT COUNT(*) FROM TripPassenger WHERE TripPassenger.TripID = Trip.TripID) AND NOW() < StartTime AND (Status IS NULL OR Status = \"Start\" OR Status != \"Cancelled\");"
	if req, err := http.NewRequest(http.MethodGet, url+query, nil); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				if body, err := io.ReadAll(res.Body); err == nil {
					var trips []Trip
					json.Unmarshal(body, &trips)
					return trips
				}
			} else {
				fmt.Println(res.StatusCode)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return nil
}

func tripList(trips []Trip, caption string, option bool) {
	fmt.Println()
	fmt.Println(caption)
	fmt.Println("=====================")
	for index, trip := range trips {
		fmt.Printf("Trip %d:\n", index+1)
		fmt.Println("=====================")
		fmt.Printf("Pickup Address: %s\n", trip.PickUpAddress)
		fmt.Printf("Alternative Address: %s\n", trip.AlternativeAddress)
		fmt.Printf("Start Time: %s\n", trip.StartTime)
		fmt.Printf("Destination: %s\n", trip.Destination)
		fmt.Printf("Pax Limit: %d\n", trip.PaxLimit)
		fmt.Printf("Trip Status: %s\n", trip.Status)
		fmt.Println("=====================")
	}
	if option {
		fmt.Println("Choose 0 to Exit")
	}
	fmt.Println()
}

func updateTripStatus(url string, trips []Trip, startOrCancel string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a Trip: ")
	var choice int
	fmt.Scanf("%d", &choice)
	reader.ReadString('\n')

	if choice < 0 || choice > len(trips) {
		fmt.Println("Invalid Trip Option")
		return
	} else if choice == 0 {
		return
	}

	trip := trips[choice-1]
	if startOrCancel == "Start" {
		trip.Status = "Start"
	} else if startOrCancel == "Cancelled" {
		trip.Status = "Cancelled"
	}

	postBody, _ := json.Marshal(trip)
	resBody := bytes.NewBuffer(postBody)

	// fmt.Println(resBody)
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, url+trip.ID, resBody); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				fmt.Println("Trip: ", trip.ID, " updated successfully")
			} else {
				fmt.Println(res.StatusCode)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func getTripListCarOwner(url string, startOrCancel bool, actionType string) []Trip {
	client := &http.Client{}
	userid := currentCarOwner.UserAcc.ID
	query := "SELECT * FROM Trip WHERE CarOwnerID = '" + userid + "'"
	if startOrCancel {
		if actionType == "start" {
			query = query + " AND StartTime BETWEEN NOW() AND (NOW() + INTERVAL 30 MINUTE) AND (Status IS NULL OR Status = \"Start\" OR Status != \"Cancelled\")"
		} else if actionType == "cancel" {
			query = query + " AND StartTime > (NOW() + INTERVAL 31 MINUTE) AND (Status IS NULL OR Status = \"Start\" OR Status != \"Cancelled\")"
		}
	}
	if req, err := http.NewRequest(http.MethodGet, url+query, nil); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				if body, err := io.ReadAll(res.Body); err == nil {
					var trips []Trip
					json.Unmarshal(body, &trips)
					return trips
				}
			} else {
				fmt.Println(res.StatusCode)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return nil
}

func createTrip(url string) {
	// getting all the user information
	reader := bufio.NewReader(os.Stdin)
	// reader.ReadString('\n')
	var newTrip Trip
	fmt.Print("Pick Up Address: ")
	pickupAddr, _ := reader.ReadString('\n')
	newTrip.PickUpAddress = strings.TrimSpace(pickupAddr)
	fmt.Print("Alternative Address: ")
	altAddr, _ := reader.ReadString('\n')
	newTrip.AlternativeAddress = strings.TrimSpace(altAddr)
	fmt.Print("Start Time (YYYY-MM-DD HH:MM:SS): ")
	startTime, _ := reader.ReadString('\n')
	startTime_str := strings.TrimSpace(startTime)
	newTrip.StartTime, _ = time.Parse("2006-01-02 15:04:05", startTime_str)
	fmt.Print("Destination: ")
	dest, _ := reader.ReadString('\n')
	newTrip.Destination = strings.TrimSpace(dest)
	fmt.Print("Pax Limit: ")
	paxLimit, _ := reader.ReadString('\n')
	paxLimit_str := strings.TrimSpace(paxLimit)
	newTrip.PaxLimit, _ = strconv.Atoi(paxLimit_str)
	newTrip.Driver = currentCarOwner

	postBody, _ := json.Marshal(newTrip)
	resBody := bytes.NewBuffer(postBody)

	// fmt.Println(resBody)

	// rest api
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, url+currentCarOwner.UserAcc.ID, resBody); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				if body, err := io.ReadAll(res.Body); err == nil {
					var tripRes Trip
					json.Unmarshal(body, &tripRes)
					fmt.Println("Trip created successful")
				} else {
					fmt.Println("Trip failed to create")

				}
			} else {
				fmt.Println(res.StatusCode)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func checkUserType(id string, url string) bool {

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, url+id, nil); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == http.StatusOK {
				if body, err := io.ReadAll(res.Body); err == nil {
					var tempCarOwner CarOwner
					json.Unmarshal(body, &tempCarOwner)
					tempCarOwner.UserAcc = currentPassenger
					currentCarOwner = tempCarOwner
					// fmt.Println(currentCarOwner)
					return true
				}
			} else {
				return false
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return false
}

func updateToCarOwner(url string) {
	reader := bufio.NewReader(os.Stdin)
	var carOwner CarOwner
	fmt.Print("License No: ")
	licenseNo, _ := reader.ReadString('\n')
	carOwner.LicenseNo = strings.TrimSpace(licenseNo)
	fmt.Print("Plate No: ")
	plateNo, _ := reader.ReadString('\n')
	carOwner.PlateNo = strings.TrimSpace(plateNo)
	carOwner.UserAcc = currentPassenger

	postBody, _ := json.Marshal(carOwner)
	resBody := bytes.NewBuffer(postBody)
	// fmt.Println(resBody)

	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, url+"newCarOwner", resBody); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				currentCarOwner = carOwner
				fmt.Println("Updated to car owner")
			} else {
				fmt.Println("Fail to updated as car owner")
			}
		} else {
			fmt.Println(err)
			fmt.Println("client req")
		}
	} else {
		fmt.Println(err)
		fmt.Println("newreq")
	}
}

func updateUserInfo(url string, userType bool) {
	reader := bufio.NewReader(os.Stdin)
	// reader.ReadString('\n')
	var updatedUser Passenger
	fmt.Print("Updated Phone Number: ")
	phoneNo, _ := reader.ReadString('\n')
	updatedUser.PhoneNo = strings.TrimSpace(phoneNo)
	fmt.Print("Updated the First Name: ")
	firstName, _ := reader.ReadString('\n')
	updatedUser.FirstName = strings.TrimSpace(firstName)
	fmt.Print("Updated Last Name: ")
	lastName, _ := reader.ReadString('\n')
	updatedUser.LastName = strings.TrimSpace(lastName)
	fmt.Print("Updated Email: ")
	email, _ := reader.ReadString('\n')
	updatedUser.Email = strings.TrimSpace(email)
	updatedUser.ID = currentPassenger.ID
	updatedUser.RegisterDT = currentPassenger.RegisterDT

	postBody, _ := json.Marshal(updatedUser)
	resBody := bytes.NewBuffer(postBody)
	// fmt.Println(resBody)

	// rest api
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPut, url+updatedUser.ID, resBody); err == nil {
		fmt.Println(url + updatedUser.ID)
		req.Header.Set("Content-Type", "application/json")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				currentPassenger = updatedUser
				currentPassenger.RegisterDT = updatedUser.RegisterDT
				fmt.Println("Update user successful")
			} else {
				fmt.Println("Fail to update user: Status Code - ", res.StatusCode)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}

	if userType {
		var carOwner CarOwner
		fmt.Print("Updated License No: ")
		licenseNo, _ := reader.ReadString('\n')
		carOwner.LicenseNo = strings.TrimSpace(licenseNo)
		fmt.Print("Updated Plate No: ")
		plateNo, _ := reader.ReadString('\n')
		carOwner.PlateNo = strings.TrimSpace(plateNo)
		carOwner.UserAcc = updatedUser
		fmt.Println(carOwner)
		url = "http://localhost:5001/api/v1/carowner/"

		postBody, _ := json.Marshal(carOwner)
		resBody := bytes.NewBuffer(postBody)
		// fmt.Println(resBody)

		client := &http.Client{}
		if req, err := http.NewRequest(http.MethodPut, url+carOwner.UserAcc.ID, resBody); err == nil {
			req.Header.Set("Content-Type", "application/json")
			if res, err := client.Do(req); err == nil {
				if res.StatusCode == 202 {
					currentCarOwner = carOwner
					fmt.Println("Update user successful")
				} else {
					fmt.Println("Fail to update user: Status Code - ", res.StatusCode)
				}
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}
}

func deleteAcc(url string, userType bool) bool {
	if userType { //true == carowner
		client1 := &http.Client{}
		if req, err := http.NewRequest(http.MethodDelete, "http://localhost:5002/api/v1/tripPassenger/"+currentCarOwner.UserAcc.ID, nil); err == nil {
			if res, err := client1.Do(req); err == nil {
				if res.StatusCode == http.StatusAccepted {
					fmt.Println("Trips of CarOwner:", currentCarOwner.UserAcc.ID, "deleted successfully")
				}
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}

		client2 := &http.Client{}
		if req, err := http.NewRequest(http.MethodDelete, "http://localhost:5002/api/v1/trip/"+currentCarOwner.UserAcc.ID, nil); err == nil {
			if res, err := client2.Do(req); err == nil {
				if res.StatusCode == http.StatusAccepted {
					fmt.Println("Trips of CarOwner:", currentCarOwner.UserAcc.ID, "deleted successfully")
				}
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}

		client3 := &http.Client{}
		if req, err := http.NewRequest(http.MethodDelete, "http://localhost:5001/api/v1/carowner/"+currentCarOwner.UserAcc.ID, nil); err == nil {
			if res, err := client3.Do(req); err == nil {
				if res.StatusCode == http.StatusAccepted {
					fmt.Println("Car Owner:", currentCarOwner.UserAcc.FirstName, "deleted successfully")
					var empty CarOwner
					currentCarOwner = empty
				}
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}

	client1 := &http.Client{}
	if req, err := http.NewRequest(http.MethodDelete, "http://localhost:5000/api/v1/passenger/passengertrip/"+currentPassenger.ID, nil); err == nil {
		if res, err := client1.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				fmt.Println("User's Enrolled Trips:", currentPassenger.ID, "deleted successfully")
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}

	client2 := &http.Client{}
	if req, err := http.NewRequest(http.MethodDelete, url+currentPassenger.ID, nil); err == nil {
		if res, err := client2.Do(req); err == nil {
			if res.StatusCode == http.StatusAccepted {
				fmt.Println("User:", currentPassenger.FirstName, "deleted successfully")
				var empty Passenger
				currentPassenger = empty
				return true
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	return false
}

func login(url string) bool {
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	var user Passenger
	fmt.Print("Login [1]Phone [2]Email: ")
	login_method, _ := reader.ReadString('\n')
	login_method = strings.TrimSpace(login_method)
	var userInfo string
	var loginType string
	if login_method == "1" {
		fmt.Print("Enter Phone Number: ")
		userInfo, _ = reader.ReadString('\n')
		userInfo = strings.TrimSpace(userInfo)
		user.PhoneNo = userInfo
		loginType = "PhoneNumber/"
	} else if login_method == "2" {
		fmt.Print("Enter Email Address: ")
		userInfo, _ = reader.ReadString('\n')
		userInfo = strings.TrimSpace(userInfo)
		user.Email = userInfo
		loginType = "Email/"
	}

	postBody, _ := json.Marshal(user)
	resBody := bytes.NewBuffer(postBody)
	// fmt.Println(resBody) // see what client side sending

	client1 := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, url+loginType+userInfo, resBody); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if res, err := client1.Do(req); err == nil {
			if res.StatusCode == 202 {
				if body, err := io.ReadAll(res.Body); err == nil {
					json.Unmarshal(body, &currentPassenger)
					// fmt.Println(currentPassenger)
					return true
				}
			} else if res.StatusCode == http.StatusUnauthorized {
				fmt.Println("User credential does not exist")
			}
		} else {
			client2 := &http.Client{}
			if req, err := http.NewRequest(http.MethodPost, "http://localhost:5001/api/v1/carowner/"+loginType+userInfo, resBody); err == nil {
				req.Header.Set("Content-Type", "application/json")
				if res, err := client2.Do(req); err == nil {
					if res.StatusCode == 202 {
						if body, err := io.ReadAll(res.Body); err == nil {
							json.Unmarshal(body, &currentCarOwner)
							fmt.Println(currentCarOwner.UserAcc)
							currentPassenger = currentCarOwner.UserAcc
							// fmt.Println(currentPassenger)
							return true
						}
					} else if res.StatusCode == http.StatusUnauthorized {
						fmt.Println("User credential does not exist")
					}
				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println(err)
	}
	return false
}

// create new user
func signup(url string) {
	// getting all the user information
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	var newUser Passenger
	fmt.Print("Enter Phone Number: ")
	phoneNo, _ := reader.ReadString('\n')
	newUser.PhoneNo = strings.TrimSpace(phoneNo)
	fmt.Print("Enter the First Name: ")
	firstName, _ := reader.ReadString('\n')
	newUser.FirstName = strings.TrimSpace(firstName)
	fmt.Print("Enter Last Name: ")
	lastName, _ := reader.ReadString('\n')
	newUser.LastName = strings.TrimSpace(lastName)
	fmt.Print("Enter Email: ")
	email, _ := reader.ReadString('\n')
	newUser.Email = strings.TrimSpace(email)
	newUser.ID = ""
	newUser.RegisterDT = time.Now()

	postBody, _ := json.Marshal(newUser)
	resBody := bytes.NewBuffer(postBody)

	// fmt.Println(resBody)

	// rest api
	client := &http.Client{}
	if req, err := http.NewRequest(http.MethodPost, url, resBody); err == nil {
		req.Header.Set("Content-Type", "application/json")
		if res, err := client.Do(req); err == nil {
			if res.StatusCode == 202 {
				if body, err := io.ReadAll(res.Body); err == nil {
					var userRes Passenger
					json.Unmarshal(body, &userRes)
					fmt.Println("New User Created!")
				}
			} else {
				fmt.Println("Fail to create user (Existing Phone Number/Email)")
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
