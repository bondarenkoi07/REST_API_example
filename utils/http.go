package utils

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, data interface{}, err error) {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)

	if data != nil && err == nil {
		w.Header().Add("Content-Type", "application/json")
		err = encoder.Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
