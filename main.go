package main

import (
	"log"
	"net/http"
)

var manager ContestManager

func handleAllContests(w http.ResponseWriter, r *http.Request) {
	contests := manager.GetAllContests()
	if len(contests) <= 0 {
		WriteJSON(w, http.StatusInternalServerError, Error{Message: "failed to get any contests"})
		return
	}

	WriteJSON(w, http.StatusOK, contests)
}

func handleSpecificPlatforms(w http.ResponseWriter, r *http.Request) {
	platforms := r.URL.Query()["p"]
	if len(platforms) <= 0 {
		WriteJSON(w, http.StatusBadRequest, Error{Message: "need at least one platform"})
		return
	}

	contests := manager.GetContestsOnPlatforms(platforms)
	if len(contests) <= 0 {
		WriteJSON(w, http.StatusInternalServerError, Error{Message: "failed to fetch contests"})
		return
	}

	WriteJSON(w, http.StatusOK, contests)
}

func main() {
	manager = NewContestManager()

	router := http.NewServeMux()

	router.HandleFunc("GET /contests/all", handleAllContests)
	router.HandleFunc("GET /contests", handleSpecificPlatforms)

	port := ":8080"
	log.Println("Listening on port"+port)
	http.ListenAndServe(port, router)
}
