package db

import (
	"testing"
)

func TestSetCache(t *testing.T) {
	var key = "TestKey"
	var value = "Trial"

	InitCache()

	Store.SetCache(key, value)

	cacheValue, err := Store.GetCache(key)

	if err != nil {
		t.Error("Cache error thrown")
	}

	if cacheValue != value {
		t.Error("Could not set cache item")
	}
}
