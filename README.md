# ZopSmart_Task

## Problem Statement
I have used a hospital's IPD as an example to perform CRUD operations with the following functionalities:
- Add a record to DB when a patient arrives
- See the list of patients currently in the hospital
- Update the record in DB when a patient is waiting, being treated or even change the values in DB
- Delete the record from DB when a patient leaves the hospital i.e Status="Completed"

## Technology Stack:

*Backend*: GoFr<br>
*Database*: Postgres SQL<br>
*Containerization*: Docker
## Testing and Validation:

All API functionalities have been thoroughly tested and validated using Postman. Access the collection for hands-on exploration: [https://lively-comet-136197.postman.co/workspace/New-Team-Workspace~f78202a5-d125-441b-91de-dfcc63352fb9/collection/31394778-7f21f2fd-0a7a-4d6b-a14b-5b978b005e4c?action=share&creator=31394778]
<br>
## How To Run
Your system should have Docker(To get the image of the Postgresql), and Go language installed
To get the dependencies run the following in a terminal
   - go get github.com/gorilla/mux"
	 - go getgithub.com/lib/pq"
In the terminal run
   -go compose build
   -go compose up go-app
This will get the API active on localhost port:8000
We can test the API on postman using the URL 0.0.0.0:8000\patients
### Unit Tests have also been created and tested for a coverage > 90% (All 4 CRUD functions GetPatients, CreatePatient, UpdatePatient, and DeletePatient have been Unit Tested). 

