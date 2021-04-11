package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, data interface{}, err error) {
	if data != nil && err == nil {
		w.Header().Set("Content-Type", "application/json")

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println("che? ", err)
		}

		_, err = w.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func DecodeJSON(r *http.Request, w http.ResponseWriter, data interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	} else {
		return nil
	}
}
