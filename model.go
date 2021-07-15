// model.go

package main

import "errors"

type temperature struct {
	Temperature float64
	Date        string
	Change      float64
	Tank        string
}

func (t *temperature) getTemperature() error {
	return errors.New("Not implemented")
}
