package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8090/some/undefined/link", nil)
	if err != nil {
		panic(err)
	}
	client := new(http.Client)

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	response, err := client.Do(req)

	if err != nil && response != nil{
		if response.StatusCode != http.StatusFound {
			t.Errorf("Redirect doesnt work, expected status code %d. Got %d.", http.StatusOK, response.StatusCode)
		}
	}
}

type RegistrationStruct struct{
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Sex      string `json:"sex"`
	Password string `json:"password"`
}

func TestCorrectRegistration(t *testing.T){
	json_req, _ := json.Marshal(&RegistrationStruct{"testname","testsurname","test@mail.ru","М","12345"})
	resp, err := http.Post("http://127.0.0.1:8090/registration", "application/json", bytes.NewReader(json_req))
	if err != nil {
		t.Error(err.Error())
	}
	status_code := resp.StatusCode

	if status_code != http.StatusOK{
		t.Errorf("Status is not 200, but %d", status_code)
	}
}

type loginStruct struct{
	Email     string
	Password  string
}

func TestIncorrectLogin(t *testing.T) {
	json_req, _ := json.Marshal(&loginStruct{"badEmail","12345"})
	resp, err := http.Post("http://127.0.0.1:8090/login", "application/json", bytes.NewReader(json_req))
	if err != nil {
		t.Error(err.Error())
	}
	status_code := resp.StatusCode
	if status_code != 403{
		t.Errorf("Status is not 403, but %d", status_code)
	}
}

func TestCorrectLogin(t *testing.T) {
	json_req, _ := json.Marshal(&loginStruct{"test@mail.ru","12345"})
	resp, err := http.Post("http://127.0.0.1:8090/login", "application/json", bytes.NewReader(json_req))
	if err != nil {
		t.Error(err.Error())
	}
	status_code := resp.StatusCode
	if status_code != http.StatusOK{
		t.Errorf("Status is not 200, but %d", status_code)
	}
}