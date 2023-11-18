package main

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CheckAddressInDatabase(address string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM address WHERE address_text = $1", address).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) SaveSearchHistory(searchQuery string) error {
	_, err := r.db.Exec("INSERT INTO search_history (search_query) VALUES ($1)", searchQuery)
	return err
}

func (r *Repository) SaveAddress(address string) (int, error) {
	var addressID int
	err := r.db.QueryRow("INSERT INTO address (address_text) VALUES ($1) RETURNING id", address).Scan(&addressID)
	return addressID, err
}

func (r *Repository) LinkAddressToSearchHistory(searchHistoryID, addressID int) error {
	_, err := r.db.Exec("INSERT INTO history_search_address (search_history_id, address_id) VALUES ($1, $2)", searchHistoryID, addressID)
	return err
}

func (r *Repository) GetAddressesFromSearchHistory(searchQuery string, similarityThreshold float64) ([]string, error) {
	rows, err := r.db.Query(`
        SELECT address.address_text
        FROM search_history
        JOIN history_search_address ON search_history.id = history_search_address.search_history_id
        JOIN address ON history_search_address.address_id = address.id
        WHERE search_history.search_query ILIKE $1
          AND similarity(search_history.search_query, $2) > $3
    `, "%"+searchQuery+"%", searchQuery, similarityThreshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []string
	for rows.Next() {
		var address string
		if err := rows.Scan(&address); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}
