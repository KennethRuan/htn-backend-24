package main

// uuid
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	db, err := sql.Open("postgres", "host=db user=user password=password dbname=postgres sslmode=disable")
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