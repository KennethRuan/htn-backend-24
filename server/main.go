package main

// uuid
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/cors"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Company string `json:"company"`
	Phone string `json:"phone"`
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
	db, err := sql.Open("postgres", "host=localhost user=user password=password dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new router
	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.Use(apiMiddleware)
	api.HandleFunc("/users", getUsers(db)).Methods("GET")
	api.HandleFunc("/users/{id}", getUser(db)).Methods("GET")
	api.HandleFunc("/users/{id}", updateUser(db)).Methods("PUT")
	api.HandleFunc("/skills/", getSkills(db)).Methods("GET")

	// Setup static file server
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/index.html")
	})

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})
	handler := c.Handler(router)

	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: handler, 
	}

	fmt.Println("Starting Server")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func apiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request for %s", r.URL.Path)
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
			u.company,
			u.phone,
			json_agg(
				JSON_BUILD_OBJECT(
					'name', s.name,
					'rating', us.rating
				)
			) AS skills
			FROM users u
			JOIN users_skills us ON u.id = us.user_id
			JOIN skills s ON us.skill_id = s.id
			GROUP BY u.id, u.name, u.email
		`)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			var skillsJSON []byte 
			if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Company, &u.Phone, &skillsJSON); err != nil {
				fmt.Println(string(skillsJSON))
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := json.Unmarshal(skillsJSON, &u.Skills); err != nil {
				fmt.Println(string(skillsJSON))
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
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
			u.company,
			u.phone,
			json_agg(
				JSON_BUILD_OBJECT(
					'name', s.name,
					'rating', us.rating
				)
			) AS skills
			FROM users u
			JOIN users_skills us ON u.id = us.user_id
			JOIN skills s ON us.skill_id = s.id
			WHERE u.id = $1
			GROUP BY u.id, u.name, u.email
		`, id)

		var u User
		var skillsJSON []byte 
		if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Company, &u.Phone, &skillsJSON); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(skillsJSON, &u.Skills); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

type UserUpdate struct {
	Name string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Company string `json:"company,omitempty"`
	Phone string `json:"phone,omitempty"`
}

func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]

		// Extract a partial user from the request body
		var u UserUpdate
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Generate the SQL query
		query, args, err := buildUpdateQuery(id, u)
		if err != nil {
			http.Error(w, "Failed to build update query", http.StatusInternalServerError)
			return
		}

		if _, err := db.Exec(query, args...); err != nil {
			http.Error(w, "Failed to execute update query", http.StatusInternalServerError)
			return
		}

		// Return the updated user
		row := db.QueryRow(`
			SELECT
			*
			FROM users
			WHERE id = $1
		`, id)
		var user User
		if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Company, &user.Phone); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func buildUpdateQuery(id string, u UserUpdate) (string, []interface{}, error) {
	query := "UPDATE users SET "
	args := []interface{}{}
	parts := []string{}
	if u.Name != "" {
			parts = append(parts, fmt.Sprintf("name = $%d", len(parts)+1))
			args = append(args, u.Name)
	}
	if u.Email != "" {
			parts = append(parts, fmt.Sprintf("email = $%d", len(parts)+1))
			args = append(args, u.Email)
	}
	if u.Company != "" {
			parts = append(parts, fmt.Sprintf("company = $%d", len(parts)+1))
			args = append(args, u.Company)
	}
	if u.Phone != "" {
			parts = append(parts, fmt.Sprintf("phone = $%d", len(parts)+1))
			args = append(args, u.Phone)
	}
	if len(parts) == 0 {
			return "", nil, fmt.Errorf("no fields to update")
	}
	query += strings.Join(parts, ", ") + fmt.Sprintf(" WHERE id = $%d", len(parts)+1)
	args = append(args, id)

	return query, args, nil
}

type Skill struct {
	id string
	Name string
}

func getSkills(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		minFreqStr := r.URL.Query().Get("min_frequency")
		minFreq, err := strconv.Atoi(minFreqStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		maxFreqStr := r.URL.Query().Get("max_frequency")
		maxFreq, err := strconv.Atoi(maxFreqStr)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("Hitting", minFreq, maxFreq)

		// Return all skills with frequency between min and max
		// Frequency has to be aggregated as it is not stored on the table
		rows, err := db.Query(`
			SELECT 
				s.name, 
				COUNT(us.user_id) AS frequency
				FROM 
				skills s
				JOIN 
				users_skills us ON s.id = us.skill_id
				GROUP BY 
				s.name
				HAVING 
				COUNT(us.user_id) >= $1 AND COUNT(us.user_id) <= $2
				ORDER BY 
				frequency DESC
			`, minFreq, maxFreq)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		skills := []SkillFrequency{}
		for rows.Next() {
			var s SkillFrequency
			if err := rows.Scan(&s.Name, &s.Frequency); err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Println(s)
			skills = append(skills, s)
		}
		if err := rows.Err(); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println(skills);
		json.NewEncoder(w).Encode(skills);
	}
}