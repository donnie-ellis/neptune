package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	totalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
	)

	temperatureGets = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "temperature_gets_total",
			Help: "Number of gets to the temperature endpoint",
		},
	)

	temperatureSets = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "temperature_post_total",
			Help: "Number of posts to the temperature endpoint",
		},
	)
	temperatureGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "temperature_of_tank",
			Help: "The temperature of a tank",
		},
		[]string{"tank"},
	)

	Temperatures []Temperature
)

func cleanString(in string) string {
	out := strings.TrimSpace(in)
	out = strings.ReplaceAll(out, " ", "_")
	out = strings.ToLower(out)

	return out
}

func setTemperature(w http.ResponseWriter, r *http.Request) {
	temperatureSets.Inc()
	fmt.Println("Endpoint: setTemperature")

	// Authenticate the request
	userName, password, hasAuth := r.BasicAuth()
	if hasAuth && validateToken(userName, password) {
		var Temp Temperature
		var existingIndex int
		var tankExists bool = false
		if tankName := r.PostFormValue("tank"); len(tankName) != 0 {
			Temp.Tank = cleanString(tankName)
		} else {
			Temp.Tank = "default"
		}
		// Check if there is already an entry for that tank
		for index, Tank := range Temperatures {
			if Tank.Tank == Temp.Tank {
				Temp = Tank
				existingIndex = index
				tankExists = true
			}
		}
		// Make sure a temperature was provided
		if temp, err := strconv.ParseFloat(r.PostFormValue("temperature"), 64); err == nil {
			Temp.Change = temp - Temp.Temperature
			Temp.Date = time.Now().Format("2006-01-02 15:04")
			Temp.Temperature = temp
			if tankExists {
				Temperatures[existingIndex] = Temp
			} else {
				Temperatures = append(Temperatures, Temp)
			}
			temperatureGauge.WithLabelValues(Temp.Tank).Set(Temp.Temperature)

			json.NewEncoder(w).Encode(Temp)
		} else {
			http.Error(w, `temperature(float64) is a required parameter`, http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func handleRequests(port string) {
	router := mux.NewRouter().StrictSlash(true)
	portNumber := fmt.Sprintf(":%v", port)
	fileServer := http.FileServer(http.Dir("./static"))
	totalRequests.Inc()
	router.Handle("/", http.StripPrefix("/", fileServer))
	router.Path("/metrics").Handler(promhttp.Handler())
	router.HandleFunc("/temperature", setTemperature).Methods("POST")
	router.HandleFunc("/temperature", returnTemperature)

	fmt.Println("Listening on ", portNumber, " . . .")
	err := http.ListenAndServe(portNumber, router)
	log.Fatal(err)
}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(temperatureSets)
	prometheus.Register(temperatureGets)
	prometheus.MustRegister(temperatureGauge)
}

func returnTemperature(w http.ResponseWriter, r *http.Request) {
	temperatureGets.Inc()
	fmt.Println("Endpoint Hit: returnTemperature")
	if len(Temperatures) == 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		json.NewEncoder(w).Encode(Temperatures)
	}
}

func validateToken(userName, password string) bool {
	key, keyPresent := os.LookupEnv("NeptuneKey")
	if keyPresent && password == key {
		return true
	}
	return false
}

func main() {
	port, portPresent := os.LookupEnv("NeptunePort")

	if !portPresent {
		port = "8000"
		fmt.Println("Port not specified using 8000")
	}

	now := time.Now()

	// For testing create a new temperature
	var Temp Temperature
	Temp.Tank = "default"
	Temp.Change = -0.5
	Temp.Date = now.Format("2006-01-02 15:04")
	Temp.Temperature = 69.5
	//Temperatures = append(Temperatures, Temp)

	handleRequests(port)
}
