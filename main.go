package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

var repository *Repository

func main() {
	db, err := sql.Open("postgres", "user=username dbname=your_database_name sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTables(db); err != nil {
		log.Fatal(err)
	}

	repository = NewRepository(db)

	router := chi.NewRouter()

	router.Get("/api/address/search", SearchAddressHandler)

	log.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}

	addresses, err := repository.GetAddressesFromSearchHistory("your_search_query", 0.7)
	if err != nil {
		return
	}

	fmt.Println(addresses)
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS search_history (
			id SERIAL PRIMARY KEY,
			search_query VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS address (
			id SERIAL PRIMARY KEY,
			address_text VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS history_search_address (
			id SERIAL PRIMARY KEY,
			search_history_id INT,
			address_id INT,
			FOREIGN KEY (search_history_id) REFERENCES search_history(id),
			FOREIGN KEY (address_id) REFERENCES address(id)
		);
	`)
	return nil
}
