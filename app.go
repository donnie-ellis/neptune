// app.go

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

type App struct {
	Router       *mux.Router
	temperatures []temperature
	prometheus   PrometheusEndpoints
	key          string
}

type PrometheusEndpoints struct {
	total_requests    prometheus.Counter
	temperature_gets  prometheus.Counter
	temperature_posts prometheus.Counter
	temperature_gauge *prometheus.GaugeVec
}

func (a *App) Initialize(applicationKey string) {
	a.Router = mux.NewRouter()
	a.prometheus.total_requests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "neptune_http_requests_total",
			Help: "Number of get requests."})
	a.prometheus.temperature_posts = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "neptune_temperature_post_total",
			Help: "Number of posts to the temperature endpoint"})
	a.prometheus.temperature_gets = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "neptune_http_requests_total",
			Help: "Number of get requests."})
	a.prometheus.temperature_gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "neptune_tank_temperature",
			Help: "The temperature of the tank"},
		[]string{"tank"})
	a.key = applicationKey
}

func (a *App) Run(port string) {}

func (a *App) getTemperature(w http.ResponseWriter, r *http.Request) {
	a.prometheus.temperature_gets.Inc()
	a.prometheus.total_requests.Inc()
	if len(a.temperatures) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		json.NewEncoder(w).Encode(a.temperatures)
	}
}

func (a *App) postTemperature(w http.ResponseWriter, r *http.Request) {
	a.prometheus.temperature_posts.Inc()
	a.prometheus.total_requests.Inc()

	var temp temperature
	var existingIndex int
	var tankExists bool

	userName, err := validateToken(r, a)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}
	fmt.Printf(userName)

	// set the tank name
	if tankName := r.PostFormValue("tank"); len(tankName) != 0 {
		temp.Tank = cleanString(tankName)
	} else {
		temp.Tank = "default"
	}

	// Check if there is already an entry for that tank
	for index, entry := range a.temperatures {
		if entry.Tank == temp.Tank {
			temp = entry
			existingIndex = index
			tankExists = true
		}
	}

	if tempString := r.PostFormValue("temperature"); len(tempString) != 0 {
		newTemp, err := strconv.ParseFloat(tempString, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("Unable to convert the temperature to a float")
		}
		temp.Change = newTemp - temp.Temperature
		temp.Date = time.Now().Format("2006-01-02 15:04")
		temp.Temperature = newTemp
		if tankExists {
			a.temperatures[existingIndex] = temp
		} else {
			a.temperatures = append(a.temperatures, temp)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("The temperature wasn't provided")
	}
}

func validateToken(r *http.Request, a *App) (userName string, err error) {
	userName, password, hasAuth := r.BasicAuth()
	if !hasAuth {
		return "", errors.New("No authentication provided")
	} else if password == a.key {
		return userName, nil
	} else {
		return userName, errors.New("Username or password was incorrect")
	}

}
