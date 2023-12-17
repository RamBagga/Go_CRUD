package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Patient struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func main() {
	//connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS patients (id SERIAL PRIMARY KEY, name TEXT, email TEXT, status TEXT)")

	if err != nil {
		log.Fatal(err)
	}

	//create router
	router := mux.NewRouter()
	router.HandleFunc("/patients", getPatients(db)).Methods("GET")
	router.HandleFunc("/patients/{id}", getPatient(db)).Methods("GET")
	router.HandleFunc("/patients", createPatient(db)).Methods("POST")
	router.HandleFunc("/patients/{id}", updatePatient(db)).Methods("PUT")
	router.HandleFunc("/patients", deletePatient(db)).Methods("DELETE")

	//start server
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// get all users
func getPatients(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM patients")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		patients := []Patient{}
		for rows.Next() {
			var u Patient
			if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Status); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			patients = append(patients, u)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(patients)
	}
}

// get user by id
func getPatient(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u Patient
		err := db.QueryRow("SELECT * FROM patients WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email, &u.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

// create user
func createPatient(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u Patient
		json.NewDecoder(r.Body).Decode(&u)

		err := db.QueryRow("INSERT INTO patients (name, email,status) VALUES ($1, $2, $3) RETURNING id", u.Name, u.Email, u.Status).Scan(&u.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

// update user
func updatePatient(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u Patient
		json.NewDecoder(r.Body).Decode(&u)

		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("UPDATE patients SET name = $1, email = $2, status = $3 WHERE id = $4", u.Name, u.Email, u.Status, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

// delete user
func deletePatient(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, err := db.Exec("DELETE FROM patients WHERE status = 'Completed'")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("Patients with status 'Completed' deleted")
	}
}
