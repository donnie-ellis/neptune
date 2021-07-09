package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func postTemp(temp float64, host string, key string, tank string) bool {
	client := &http.Client{}
	//URL := "http://localhost:8000/temperature"
	data := url.Values{}
	data.Add("temperature", strconv.FormatFloat(temp, 'f', 1, 64))
	if tank != "" {
		data.Add("tank", tank)
	}

	req, err := http.NewRequest("POST", host+"/temperature", bytes.NewBufferString(data.Encode()))
	req.SetBasicAuth("Test", key)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	if err != nil {
		log.Println(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return false
	} else {
		fmt.Println(resp.Body)
		return true
	}
}

func main() {
	tank := os.Getenv("TridentTank")
	host, hostPresent := os.LookupEnv("NeptuneHost")
	key, keyPresent := os.LookupEnv("TridentKey")

	if !hostPresent || !keyPresent {
		log.Println("Key present", keyPresent)
		log.Println("Host present", hostPresent)
		log.Fatalln("Missing required environment variable")
	}

	postTemp(59.123, host, key, tank)

}
