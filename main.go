// main.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ()

func cleanString(in string) string {
	out := strings.TrimSpace(in)
	out = strings.ReplaceAll(out, " ", "_")
	out = strings.ToLower(out)

	return out
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

func vsalidateToken(userName, password string) bool {
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
	key, keyPresent := os.LookupEnv("NeptuneKey")

	if !keyPresent {
		log.Panicf("You need to specify a key")
	}
	a := App{}
	a.Initialize(key)
	a.Run(port)
}

/* func main() {
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
*/
