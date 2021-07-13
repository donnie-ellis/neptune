package main

import "testing"

func TestCleanString(t *testing.T) {
	str := cleanString(" This Is A Test ")

	if str != "this_is_a_test" {
		t.Log("Error: expected this is a test, but got ", str)
		t.Fail()
	}
}
