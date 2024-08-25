package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(w, http.StatusOK, payload)

	/*
		 	out, _ := json.MarshalIndent(payload, "", "\t")
			fmt.Println(out)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			w.Write(out)
	*/

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	log.Println("hit broker-HandleSubmission 1")
	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	log.Println("hit broker-HandleSubmission 2")
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJson(w, errors.New("unknown action"))
	}
	log.Println("hit broker-HandleSubmission 3")
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	//decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJson(w, http.StatusAccepted, payload)
}
