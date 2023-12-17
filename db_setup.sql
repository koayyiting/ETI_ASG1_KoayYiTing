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
    Status varchar(20),
    FOREIGN KEY (CarOwnerID) REFERENCES CarOwner(UserID)
);
