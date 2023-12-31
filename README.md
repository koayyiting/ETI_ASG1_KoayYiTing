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
Platform Used                  - Did not have time to test Web service

### Assumptions:
System only allows 1 registed phone number and email <br> <br>
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

#### Description
Client Application (Console): This is the user interface where passengers and car owners interact with the platform. It acts as a client that utilizes the API Gateway for communication. <br>
Microservices: I have designed three main microservices: Passenger, Trip, and Carowner. Each microservice exposes a set of endpoints (GET, POST, PUT, DELETE) corresponding to specific functionalities within its domain. These endpoints act as the "API Gateway" in this context. <br>
Database (version1): This stores all the application data, including passenger information, car owner details, trip data, and passenger-trip relationships using MySQL. The database tables (Passenger, CarOwner, Trip, and TripPassenger) represent the entities and their relationships. <br>

Instructions to setup
=====================
Microservices Setup: (can be found under [/microservices/](https://github.com/koayyiting/ETI_ASG1_KoayYiTing/tree/main/microservices) folder) <br>
Run the following .exe files
- carowner.exe
- passenger.exe
- trip.exe

User Interface Setup: (can be found under [/microservices/](https://github.com/koayyiting/ETI_ASG1_KoayYiTing/tree/main/microservices) folder) <br>
Run the following .exe file
- console.exe <br>
From this you can start using the Application

Database Setup Queries
======================
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
This DB Setup code can be found in [/db_setup.sql](https://github.com/koayyiting/ETI_ASG1_KoayYiTing/blob/main/db_setup.sql)

### Examples to test: (Deletion of Accounts)
#### Test if Passenger can delete Account (after 1 year)
Passenger Details: Phone No. = 91010101, Email = koay@gmail.com
```
INSERT INTO Passenger (ID, PhoneNumber, Email, FirstName, LastName, registerDT)
VALUES("550e8400-e29b-41d4-a716-446655440000","91010101", "koay@gmail.com", "Yi Ting", "Koay",DATE_SUB(CURDATE(), INTERVAL 2 YEAR));
```
#### Enrolled Passenger can delete Account (after 1 year)
Passenger Details: Phone No. = 80808080, Email = derick@gmail.com
```
INSERT INTO Passenger (ID, PhoneNumber, Email, FirstName, LastName, registerDT)
VALUES("b81f9ee8-05ed-44fd-9af1-e72bed98467c","80808080", "derick@gmail.com", "Derick", "Lee",DATE_SUB(CURDATE(), INTERVAL 2 YEAR));
insert into TripPassenger (TripID, PassengerID)
values("2f192d4a-d722-4520-b66c-6a89796f297e","b81f9ee8-05ed-44fd-9af1-e72bed98467c");
```
#### Test if Car Owner Account can delete (have published trip, user have enrolled) 
CarOwner Details: Phone No. = 88880000, Email = janis@gmail.com <br>
Passenger Details: Phone No. = 88998899, Email = test@gmail.com <br>
```
INSERT INTO Passenger (ID, PhoneNumber, Email, FirstName, LastName, registerDT)
VALUES("ed2e4494-2837-424b-9dc4-0669bdf8b331","88880000", "janis@gmail.com", "Janis", "Lim",DATE_SUB(CURDATE(), INTERVAL 2 YEAR));
INSERT INTO CarOwner (UserID, LicenseNo, PlateNo)
VALUES ('ed2e4494-2837-424b-9dc4-0669bdf8b331', 'ABC123456', 'XYZ789');
INSERT INTO Trip (TripID, CarOwnerID, PickUpAddress, AlternativeAddress, StartTime, Destination, PaxLimit)
VALUES ('123e4567-e89b-12d3-a456-426614174001','ed2e4494-2837-424b-9dc4-0669bdf8b331','Chua Chu Kang Avenue 5','Chua Chu Kang Avenue 4','2023-12-17 20:00:00','Jurong East MRT', 4);
INSERT INTO Passenger (ID, PhoneNumber, Email, FirstName, LastName, registerDT)
VALUES("492e4ed4-442b-3728-c9d4-8bf39066bd31","88998899", "test@gmail.com", "Shaniah", "Santiago",DATE_SUB(CURDATE(), INTERVAL 2 YEAR));
insert into TripPassenger (TripID, PassengerID)
values("123e4567-e89b-12d3-a456-426614174001","492e4ed4-442b-3728-c9d4-8bf39066bd31");
```
