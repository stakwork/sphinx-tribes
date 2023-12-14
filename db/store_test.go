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

func TestDeleteCache(t *testing.T) {
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

	Store.DeleteCache(key)
	_, errD := Store.GetCache(key)

	if errD == nil {
		t.Error("Could not delete cache item")
	}
}

func TestSetLnCache(t *testing.T) {
	var key = "TestLnKey"

	var value = LnStore{
		K1:     "887775666900000056890P23",
		Key:    "0000000000000000000000000000000000000",
		Status: false,
	}

	InitCache()
	Store.SetLnCache(key, value)
	cacheValue, err := Store.GetLnCache(key)

	if err != nil {
		t.Error("Cache error thrown")
	}

	if cacheValue != value {
		t.Error("Could not set cache item")
	}
}
