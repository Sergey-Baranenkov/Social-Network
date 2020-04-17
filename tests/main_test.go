package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestSecretPage(t *testing.T) {
	res, err := http.Get("http://127.0.0.1:8090/secretpage")

	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	b := string(body)
	if b != "U have no permission there" {
		t.Errorf("Secret form handled incorrectly, expecGOted: U have no permission there, got %s", b)
	}
}

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

func TestEmptyForm(t *testing.T) {

	res, err := http.PostForm("http://127.0.0.1:8090/login", nil)

	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	b:=string(body)
	if b != "Поля не заполнены" {
		t.Errorf("Empty form handled incorrectly, expected: Поля не заполнены, got %s", b)
	}
}

func TestIncorrectLogin(t *testing.T) {

	res, err := http.PostForm("http://127.0.0.1:8090/login", url.Values{"email": {"undefined@lol.en"}, "password": {"1234"}})

	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	b:= string(body)
	if b != "Incorrect email/pass combination" {
		t.Errorf("Incorrect login, expexted: Incorrect email/pass combination, got %s ", b)
	}
}

func TestCorrectLogin(t *testing.T) {
	_, err := http.PostForm("http://127.0.0.1:8090/registration", url.Values{"email": {"good@yandex.ru"}, "password": {"1234"}, "first_name": {"Vasya"}, "last_name":{"Pupkin"}})

	if err != nil {
		t.Fatal(err)
	}

	loginRes, err:= http.PostForm("http://127.0.0.1:8090/login", url.Values{"email": {"good@yandex.ru"}, "password": {"1234"}})

	if err != nil {
		t.Fatal(err)
	}

	if loginRes.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d. Got %d.", http.StatusOK, loginRes.StatusCode)
	}
}