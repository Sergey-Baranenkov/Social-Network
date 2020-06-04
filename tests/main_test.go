package main

import (
	"net/http"
	"testing"
	"os"
	"fmt"
)

func TestCorrectProfileInfo(t *testing.T) {
	fmt.Println("PORT", os.Getenv("PORT"))
	resp, err := http.Get("http://127.0.0.1:" + os.Getenv("PORT") + "/server/messenger/get_short_profile_info?conversationId=1")
	if err != nil {
		t.Error(err.Error())
	}
	status_code := resp.StatusCode

	if status_code != http.StatusOK{
		t.Errorf("Status is not 200, but %d", status_code)
	}
}

func TestIncorrectProfileInfo(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:" + os.Getenv("PORT") + "/server/messenger/get_short_profile_info?conversationId=undefined")
	if err != nil {
		t.Error(err.Error())
	}
	status_code := resp.StatusCode
	if status_code != 400{
		t.Errorf("Status is not 400, but %d", status_code)
	}
}

func TestCorrectAboutInfo(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:" + os.Getenv("PORT") + "/server/about_me/select_extended_user_info?userId=1")
	if err != nil {
		t.Error(err.Error())
	}
	status_code := resp.StatusCode
	if status_code != http.StatusOK{
		t.Errorf("Status is not 200, but %d", status_code)
	}
}

func TestIncorrectAboutInfo(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:" + os.Getenv("PORT") + "/server/about_me/select_extended_user_info?userId=undefined")
	if err != nil {
		t.Error(err.Error())
	}
	status_code := resp.StatusCode
	if status_code != 400{
		t.Errorf("Status is not 400, but %d", status_code)
	}
}