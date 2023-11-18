package main

import (
	"fmt"
	"net/http"
)

func SearchAddressHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("query")

	if exists, err := repository.CheckAddressInDatabase(searchQuery); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if exists {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Address found in the database: %s", searchQuery)))
	} else {
		addressID, err := repository.SaveAddress(searchQuery)
		if err != nil {
			// Обработка ошибки
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := repository.SaveSearchHistory(searchQuery); err != nil {
			// Обработка ошибки
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		searchHistoryID := 1
		if err := repository.LinkAddressToSearchHistory(searchHistoryID, addressID); err != nil {
			// Обработка ошибки
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Address found from Dadata.ru: %s", searchQuery)))
	}
}
