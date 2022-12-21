package handler

import (
	"encoding/json"

	schemas "echo-postgres-sample/app/api/schemas"
	"log"
	"net/http"
)

//-- UTILS --

func handleErr(w http.ResponseWriter, err error) {
	res := &schemas.CommentResponse{
		Success: false,
		Error:   err.Error(),
		Comment: nil,
	}
	err = json.NewEncoder(w).Encode(res)
	//if there's an error with encoding handle it
	if err != nil {
		log.Printf("error sending response %v\n", err)
	}
	//return a bad request and exist the function
	w.WriteHeader(http.StatusBadRequest)
}

func handleDBFromContextErr(w http.ResponseWriter) {
	res := &schemas.CommentResponse{
		Success: false,
		Error:   "could not get the DB from context",
		Comment: nil,
	}
	err := json.NewEncoder(w).Encode(res)
	//if there's an error with encoding handle it
	if err != nil {
		log.Printf("error sending response %v\n", err)
	}
	//return a bad request and exist the function
	w.WriteHeader(http.StatusBadRequest)
}
