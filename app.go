// app.go

package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	a.intializeRoutes()
}

func (a *App) Run(port string) {
	log.Println("Starting the server at localhost:", port)
	log.Fatal(http.ListenAndServe(":"+port, a.Router))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "applcation/json")
	w.WriteHeader(code)
	w.Write(response)
}

// getTemperature
// handles gets to /temperature
func (a *App) getTemperature(w http.ResponseWriter, r *http.Request) {
	a.prometheus.temperature_gets.Inc()
	a.prometheus.total_requests.Inc()
	if len(a.temperatures) == 0 {
		respondWithError(w, http.StatusNotFound, "temperature not found")
	} else {
		respondWithJSON(w, http.StatusOK, a.temperatures)
	}
}

// postTemperature
// handles posts to /temperature
func (a *App) postTemperature(w http.ResponseWriter, r *http.Request) {
	a.prometheus.temperature_posts.Inc()
	a.prometheus.total_requests.Inc()

	var temp temperature
	var existingIndex int
	var tankExists bool

	if userName, err := validateToken(r, a); err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	} else {
		// set the tank name
		if len(userName) != 0 {
			temp.Tank = cleanString(userName)
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
			if newTemp, err := strconv.ParseFloat(tempString, 64); err != nil {
				respondWithError(w, http.StatusBadRequest, "Unable to convert the temperature to a float")
			} else {
				temp.Change = newTemp - temp.Temperature
				temp.Date = time.Now().Format("2006-01-02 15:04")
				temp.Temperature = newTemp
				if tankExists {
					a.temperatures[existingIndex] = temp
				} else {
					a.temperatures = append(a.temperatures, temp)
					respondWithJSON(w, http.StatusCreated, temp)
				}
			}
		} else {
			respondWithJSON(w, http.StatusBadRequest, "The temperature wasn't provided")
		}
	}
}

// validateToken
// takes a http request and validates the basic authentication
// returns the user name
func validateToken(r *http.Request, a *App) (userName string, err error) {
	userName, password, hasAuth := r.BasicAuth()
	if !hasAuth {
		return "", errors.New("no authentication provided")
	} else if password == a.key {
		return userName, nil
	} else {
		return userName, errors.New("username or password was incorrect")
	}
}

// cleanString
// removes whitespace from a string and makes it lowercase
func cleanString(in string) string {
	out := strings.TrimSpace(in)
	out = strings.ReplaceAll(out, " ", "_")
	out = strings.ToLower(out)

	return out
}

func (a *App) intializeRoutes() {
	a.Router.HandleFunc("/temperature", a.getTemperature).Methods("GET")
	a.Router.HandleFunc("/temperature", a.postTemperature).Methods("POST")
	a.Router.Handle("/metrics", promhttp.Handler())
}
