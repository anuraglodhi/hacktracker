package main

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return err
	}

	return nil
}