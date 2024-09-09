package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetPaginationParams(r *http.Request) (int, int, string, string, string) {
	// there are cases when the request is not passed in
	if r == nil {
		return 0, 1, "updated", "asc", ""
	}

	keys := r.URL.Query()
	page := keys.Get("page")
	limit := keys.Get("limit")
	sortBy := keys.Get("sortBy")
	direction := keys.Get("direction")
	search := keys.Get("search")

	// convert string to int
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)

	if intPage == 0 {
		intPage = 1
	}
	if intLimit == 0 {
		intLimit = 1
	}
	if sortBy == "" {
		sortBy = "created"
	}
	if direction == "" {
		direction = "desc"
	}

	// offset for page, start index
	offset := 0
	if intLimit > 0 && intPage > 0 {
		// this will give us an offset that includes part of the next/previous page,
		// so that all results arent replaced, a "page shifting" effect
		offset = (intPage - 1) * intLimit
	}

	return offset, intLimit, sortBy, direction, search
}

func BuildSearchQuery(key string, term string) (string, string) {
	arg1 := key + " LIKE ?"
	arg2 := "%" + term + "%"
	return arg1, arg2
}

func BuildKeysendBodyData(amount uint, receiver_pubkey string, route_hint string) string {
	var bodyData string
	memoText := "memotext added for notification"
	if route_hint != "" {
		bodyData = fmt.Sprintf(`{"amount": %d, "destination_key": "%s", "route_hint": "%s", "text": "%s"}`, amount, receiver_pubkey, route_hint, memoText)
	} else {
		bodyData = fmt.Sprintf(`{"amount": %d, "destination_key": "%s", "text": "%s"}`, amount, receiver_pubkey, memoText)
	}

	return bodyData
}

func BuildV2KeysendBodyData(amount uint, receiver_pubkey string, route_hint string) string {
	amountMsat := amount * 1000
	var bodyData string
	if route_hint != "" {
		bodyData = fmt.Sprintf(`{"amt_msat": %d, "dest": "%s", "route_hint": "%s", "wait": true}`, amountMsat, receiver_pubkey, route_hint)
	} else {
		bodyData = fmt.Sprintf(`{amt_msat": %d, "dest": "%s", "route_hint": "", "wait": true}`, amountMsat, receiver_pubkey)
	}

	return bodyData
}
