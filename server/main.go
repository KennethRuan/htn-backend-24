package main

// uuid
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Skills []UserSkill `json:"skills"`
}

type UserSkill struct {
	Name string `json:"name"`
	Rating int `json:"rating"`
}

type SkillFrequency struct {
	Name string `json:"name"`
	Frequency int `json:"frequency"`
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new router
	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.Use(jsonMiddleware)
	api.HandleFunc("/api/users", getUsers(db)).Methods("GET")
	api.HandleFunc("/api/users/{id}", getUser(db)).Methods("GET")
	// router.HandleFunc("/api/users/{id}", updateUser(db)).Methods("PUT")
	api.HandleFunc("/api/skills", getSkills(db)).Methods("GET")

	// Setup static file server
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/index.html")
	})

	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: router, 
	}

	fmt.Println("Starting Server")
	if err := server.ListenAndServe(); err != nil {
	fmt.Println(err)
	}
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
					w.Header().Set("Content-Type", "application/json")
			}
			next.ServeHTTP(w, r)
	})
}


func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query(`
			SELECT 
			u.id, 
			u.name, 
			u.email,
			ARRAY_AGG(
				JSON_BUILD_OBJECT(
					'name', s.name,
					'rating', us.rating
				)
			) AS skills
			FROM users u
			JOIN user_skills us ON u.id = us.user_id
			JOIN skills s ON us.skill_id = s.id
			GROUP BY u.id, u.name, u.email
		`)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Skills); err != nil {
				log.Fatal(err)
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(users);
	}
}

func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]


		row := db.QueryRow(`
			SELECT
			u.id,
			u.name,
			u.email,
			ARRAY_AGG(
				JSON_BUILD_OBJECT(
					'name', s.name,
					'rating', us.rating
				)
			) AS skills
			FROM users u
			JOIN user_skills us ON u.id = us.user_id
			JOIN skills s ON us.skill_id = s.id
			WHERE u.id = $1
		`, id)

		var u User
		if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Skills); err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(u)
	}
}

func getSkills(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		minFreq := r.URL.Query().Get("min_frequency")
		maxFreq := r.URL.Query().Get("max_frequency")

		// Return all skills with frequency between min and max
		// Frequency has to be aggregated as it is not stored on the table
		rows, err := db.Query(`
			SELECT
				s.name, 
				COUNT(s.name) AS frequency 
				FROM skills s 
				GROUP BY s.name 
				HAVING COUNT(s.name) >= $1 AND COUNT(s.name) <= $2`,
			minFreq, maxFreq)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		skills := []SkillFrequency{}
		for rows.Next() {
			var s SkillFrequency
			if err := rows.Scan(&s.Name, &s.Frequency); err != nil {
				log.Fatal(err)
			}
			skills = append(skills, s)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(skills);
	}
}