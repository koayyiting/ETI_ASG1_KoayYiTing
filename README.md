Design Consideration
====================
Application Description:
This carpooling platform uses microservices to manage passengers, car owners and trips. Passengers create accounts and join published trips, while car owners define trips, manage seats (first-come-first-serve), and cancel trip 30 minutes before start time. Users can access past trips and update profiles, with car owners providing additional driver's license and car details. After 1 year of inactivity, users can delete accounts, but data is retained for audit purposes.

Microservices <br>
Passenger Service: This service focus solely on managing passenger data and actions, such as sign-up, login, profile updates, passenger trip enrollment, and delete account. <br>
Trip Service: This service handle trip creation, getting trip list, and cancelling trips. <br>
Carowner Service: This service manages carowner data and actions, such as car registration, update account, and delete account. <br>

Factors and Constraints <br>
User Experience: <br>
Simple and intuitive interface - Both passengers and car owners should be able to easily navigate the platform. <br>
Clear and concise information  - Trip details, user profiles, and important platform functionalities should be readily available and understandable. <br>
User Interface                 - Console/exe/command prompt 

Technical:
Microservice architecture      - Created 3 microservices.

Constraints:
Scalability                    - Did not have time to test Web service

Architecture Diagram
====================
![image](https://github.com/koayyiting/ETI_ASG1_KoayYiTing/assets/93900494/322616d1-19a8-40ea-b1a7-d99bf4ecdd87)

Instructions to setup
=====================
Microservices Setup:
Run the following .exe files
- carowner.exe
- version1.exe (passenger)
- trip.exe

User Interface Setup:
Run the following .exe file
- console.exe
From this you can start using the Application

Database Setup Queries: <br>
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
