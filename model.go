// model.go

package main

import "time"

type temperature struct {
	Temperature float64 `json:"temperature"`
	Date        string  `json:"Date"`
	Change      float64 `json:"Change"`
	Tank        string  `json:"tank"`
}

func (t *temperature) setTemperature(temp float64) {
	t.Change = temp - t.Temperature
	t.Date = time.Now().Format("2006-01-02 15:04")
	t.Temperature = temp
}
