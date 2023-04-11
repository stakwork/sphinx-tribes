package main

import "testing"

func TestInitJwt(t *testing.T) {

	initConfig()
	initJwt()

	if TokenAuth == nil {
		t.Error("Could not init JWT")
	} else {
		t.Log("JWT inited successfully")
	}
}
