package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCleanString(t *testing.T) {
	str := cleanString(" This Is A Test ")

	if str != "this_is_a_test" {
		t.Log("Error: expected this is a test, but got ", str)
		t.Fail()
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetTemperature(t *testing.T) {
	req, _ := http.NewRequest("GET", "/temperature", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestPostTemperature(t *testing.T) {
	data := url.Values{}
	data.Add("temperature", "68.8")
	data.Add("tank", "test")
	req, _ := http.NewRequest("POST", "/temperature", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["temperature"] != "68.8" {
		t.Errorf("Expected temperature to be '68.8'. Got '%v'", m["temperature"])
	}

	if m["tank"] != "test" {
		t.Errorf("Expected tank name to be 'test'. Got '%v'", m["tank"])
	}

}
