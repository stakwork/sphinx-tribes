package main

import (
	lnurl "github.com/fiatjaf/go-lnurl"
)

func encodeLNURL() (string, error) {
	println("Hello")
	return lnurl.Encode("Hello")
}
