Design Consideration
====================
## Application Description
This carpooling platform uses microservices to manage passengers, car owners and trips. These microservices are loosely cooupled and working together. Passengers create accounts and join published trips, while car owners define trips, manage seats (first-come-first-serve), and cancel trip 30 minutes before start time. Users can access past trips and update profiles, with car owners providing additional driver's license and car details. After 1 year of inactivity, users can delete accounts, but data is retained for audit purposes.

## Microservices
Passenger Service: This service focus solely on managing passenger data and actions, such as sign-up, login, profile updates, change to car owner, passenger trip enrollment, and delete account. <br>
Trip Service: This service handle getting trip list, trip creation, and cancelling trips. <br>
Carowner Service: This service manages carowner data and actions, such as update account, and delete account, login. <br>

## Factors and Constraints
### User Experience: <br>
Simple and intuitive interface - Both passengers and car owners should be able to easily navigate the platform. <br>
Clear and concise information  - Trip details, user profiles, and important platform functionalities should be readily available and understandable. <br>
User Interface                 - Console/exe/command prompt 

### Technical:
Microservice architecture      - Created 3 microservices (Passenger, Car Owner, Trip). <br>
They are designed to be loosely cooupled but working together. <br>
Example: (Login function) <br>
- Both the Passenger and CarOwner allows user login. (For Car owner to login by itself, the user needs to be upgraded to Car Owner first)
- Means: if passenger.exe not running Car Owner user can still login into the car pool platform

### Constraints:
Scalability                    - Did not have time to test Web service

### Assumptions:
System only allows 1 regiested phone number and email <br> <br>
Passenger can update account to Car Owner <br>
Passenger can enrol Trips <br>
Passenger account can be deleted (after 1 year of registration) and enrolled trips will be deleted <br> <br>
Car Owner can publish Trips <br>
Car Owner account can be deleted (after 1 year of registration) this include all the published trip by the car owner will be deleted (which also means passengers that is enrolled to the carowner's trip will be deleted) <br> <br>
Trips have status "Start", "Cancelled" or default value "null"  <br>
Trip status can be changed to "Cancelled" if (current time < start time) by 30 minutes <br>
Trip status can be changed to "Start" if within 30 minutes of current time and start time (current time < start time) <br>

### Features/Functions (UI)
#### Main Page
1. Login
2. Sign up
3. Exit Program
   
#### Passenger
1. Update User Details
2. Delete Account
3. Change to Car Owner
4. Enrol Trip
5. View Account Detail
6. List of Enrolled Trips
7. Search for Trips
8. List all Trips
0. Log Out

#### Car Owner
1. Update User Details
2. Delete Account
3. Publish Trip
4. List of Published Trip
5. Start/Cancel Published Trip
6. View Account Detail
0. Log Out

Architecture Diagram
====================
![image](https://github.com/koayyiting/ETI_ASG1_KoayYiTing/assets/93900494/5ddff562-a0c9-4ac7-8d04-82fe48b5269d)


Instructions to setup
=====================
Microservices Setup: (can be found under /microservices/ folder) <br>
Run the following .exe files
- carowner.exe
- passenger.exe
- trip.exe

User Interface Setup: <br>
Run the following .exe file
- console.exe
From this you can start using the Application

Database Setup Queries: <br>
```
CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost';

CREATE database version1;

USE version1;
CREATE TABLE Passenger (
ID varchar(36) NOT NULL PRIMARY KEY,
PhoneNumber varchar(50), 
Email varchar (320),
FirstName varchar(100), 
LastName varchar(50),
RegistedDT datetime
);

CREATE TABLE CarOwner (
UserID varchar(36),
LicenseNo varchar(9),
PlateNo varchar(10), 
FOREIGN KEY (UserID) REFERENCES Passenger(ID)
);

CREATE TABLE TripPassenger(
	TripID varchar(36),
    PassengerID varchar(36),
    PRIMARY KEY (TripID, PassengerID),
    FOREIGN KEY (TripID) REFERENCES Trip(TripID),
    FOREIGN KEY (PassengerID) REFERENCES Passenger(ID)
);

CREATE TABLE Trip (
    TripID varchar(36) NOT NULL PRIMARY KEY,
    CarOwnerID varchar(36),
    PickUpAddress varchar(255),
    AlternativeAddress varchar(255),
    StartTime datetime,
    Destination varchar(255),
    PaxLimit int,
    FOREIGN KEY (CarOwnerID) REFERENCES CarOwner(UserID)
);
```
