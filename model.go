// model.go

package main

type temperature struct {
	Temperature float64 `json:"temperature"`
	Date        string  `json:"Date"`
	Change      float64 `json:"Change"`
	Tank        string  `json:"tank"`
}
