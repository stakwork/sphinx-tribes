package main

import (
	"net/http"
	"strconv"
)

func getPaginationParams(r *http.Request) (int, int, string, string) {

	// there are cases when the request is not passed in
	if r == nil {
		return 0, -1, "updated", "asc"
	}

	keys := r.URL.Query()
	page := keys.Get("page")
	limit := keys.Get("limit")
	sortBy := keys.Get("sortBy")
	direction := keys.Get("direction")

	// convert string to int
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)

	if intPage == 0 {
		intPage = 1
	}
	if intLimit == 0 {
		intLimit = -1
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
		offset = (intPage - 1) * intLimit
	}

	return offset, intLimit, sortBy, direction
}
