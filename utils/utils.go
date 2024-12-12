package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func GetPaginationParams(r *http.Request) (int, int, string, string, string) {
	// there are cases when the request is not passed in
	if r == nil {
		return 0, 1, "updated", "asc", ""
	}

	keys := r.URL.Query()
	// trim spaces
	page := strings.TrimSpace(keys.Get("page"))
	// trim spaces
	limit := strings.TrimSpace(keys.Get("limit"))
	// convert to lowercase and trim spaces
	sortBy := strings.ToLower(strings.TrimSpace(keys.Get("sortBy")))
	direction := strings.ToLower(strings.TrimSpace(keys.Get("direction")))
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
	if direction == "" || direction != "asc" {
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
	// trim spaces
	key = strings.TrimSpace(key)
	term = strings.TrimSpace(term)

	arg1 := key + " LIKE ?"
	arg2 := "%" + term + "%"
	return arg1, arg2
}

func BuildKeysendBodyData(amount uint, receiver_pubkey string, route_hint string, memo string) string {
	var bodyData string
	if route_hint != "" {
		bodyData = fmt.Sprintf(`{"amount": %d, "destination_key": "%s", "route_hint": "%s", "text": "%s", "data": "%s"}`, amount, receiver_pubkey, route_hint, memo, memo)
	} else {
		bodyData = fmt.Sprintf(`{"amount": %d, "destination_key": "%s", "text": "%s", "data": "%s"}`, amount, receiver_pubkey, memo, memo)
	}

	return bodyData
}

func BuildV2KeysendBodyData(amount uint, receiver_pubkey string, route_hint string, memo string) string {
	// convert amount to msat
	amountMsat := amount * 1000

	// trim the memo
	memo = strings.TrimSpace(memo)
	// trim the route hint
	route_hint = strings.TrimSpace(route_hint)

	var bodyData string
	if route_hint != "" {
		bodyData = fmt.Sprintf(`{"amt_msat": %d, "dest": "%s", "route_hint": "%s", "data": "%s", "wait": true}`, amountMsat, receiver_pubkey, route_hint, memo)
	} else {
		bodyData = fmt.Sprintf(`{"amt_msat": %d, "dest": "%s", "route_hint": "", "data": "%s", "wait": true}`, amountMsat, receiver_pubkey, memo)
	}

	return bodyData
}

func BuildV2ConnectionCodes(amt_msat uint64, alias string, pubkey string, route_hint string) string {
	bodyData := fmt.Sprintf(`{
		"amt_msat": %d,
		"alias": "%s",
		"pubkey": "%s",
		"route_hint": "%s"
	}`, amt_msat, alias, pubkey, route_hint)
	return bodyData
}
