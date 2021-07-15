// main.go

package main

import (
	"fmt"
	"log"
	"os"
)

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
