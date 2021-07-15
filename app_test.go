package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize("testKey")
	code := m.Run()
	os.Exit(code)
}

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
	var tank temperature
	tank.Change = 67.9
	tank.Temperature = 67.9
	tank.Tank = "test"
	a.temperatures = append(a.temperatures, tank)
	req, _ := http.NewRequest("GET", "/temperature", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	var k []temperature
	if err := json.Unmarshal(response.Body.Bytes(), &k); err != nil {
		t.Errorf("The response wasn't in a compatible format, the message is %s", err)
	}
	found := false
	for _, tank := range k {
		if tank.Tank == "test" {
			if tank.Temperature != 67.9 {
				t.Errorf("expected a temperature of 67.9 for the tank 'test' got %s", strconv.FormatFloat(tank.Temperature, 'f', 1, 64))
			}
			found = true
		}
	}
	if !found {
		t.Errorf("expected a tank named 'test' that wasn't found")
	}
}

func TestGetNoTemperatures(t *testing.T) {
	req, _ := http.NewRequest("GET", "/temperature", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &m); err != nil {
		t.Errorf("The response wasn't in a compatible format, the message is %s", err)
	}
	if m["error"] != "temperature not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'temperature not found'. Got '%s'", m["error"])
	}
}

func TestPostTemperature(t *testing.T) {
	data := url.Values{}
	data.Add("temperature", "68.8")
	data.Add("tank", "test")
	req, _ := http.NewRequest("POST", "/temperature", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("test", "testKey")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var k temperature
	if err := json.Unmarshal(response.Body.Bytes(), &k); err != nil {
		t.Errorf("The response wasn't in a compatible format, the message is %s", err)
	}

	if k.Temperature != 68.8 {
		t.Errorf("Expected temperature to be '68.8'. Got '%v'", strconv.FormatFloat(k.Temperature, 'f', 1, 64))
	}

	if k.Tank != "test" {
		t.Errorf("Expected tank name to be 'test'. Got '%v'", k.Tank)
	}
}

func TestPostWrongToken(t *testing.T) {
	data := url.Values{}
	data.Add("temperature", "68.9")
	req, _ := http.NewRequest("POST", "/temperature", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("test", a.key+"1")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	var m map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &m); err != nil {
		t.Errorf("The response wasn't in a compatible format, the message is %s", err)
	}
	if m["error"] != "username or password was incorrect" {
		t.Errorf("Expected an error of 'username or password was incorrect' and received %s", m["error"])
	}
}

func TestPostBadTemperature(t *testing.T) {
	data := url.Values{}
	data.Add("temperature", "test")
	data.Add("tank", "test")
	req, _ := http.NewRequest("POST", "/temperature", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("test", "testKey")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &m); err != nil {
		t.Errorf("The response wasn't in a compatible format, the message is %s", err)
	}
	if m["error"] != "Unable to convert the temperature to a float" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Unable to convert the temperature to a float'. Got '%s'", m["error"])
	}
}
