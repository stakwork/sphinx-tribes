package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
)

// Tribe struct
type Tribe struct {
	UUID            string         `json:"uuid"`
	OwnerPubKey     string         `json:"ownerPubkey"`
	OwnerAlias      string         `json:"ownerAlias"`
	GroupKey        string         `json:"groupKey"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	Tags            pq.StringArray `json:"tags"`
	Img             string         `json:"img"`
	PriceToJoin     int64          `json:"priceToJoin"`
	PricePerMessage int64          `json:"pricePerMessage"`
	Created         *time.Time     `json:"created"`
	Updated         *time.Time     `json:"updated"`
}

// PropertyMap ...
type PropertyMap map[string]interface{}

// Value ...
func (p PropertyMap) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

// Scan ...
func (p *PropertyMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var i interface{}
	if err := json.Unmarshal(source, &i); err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("type assertion .(map[string]interface{}) failed")
	}

	return nil
}
