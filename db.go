package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type database struct {
	db *gorm.DB
}

// DB is the object
var DB database

func initDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		rdsHost := os.Getenv("RDS_HOSTNAME")
		rdsPort := os.Getenv("RDS_PORT")
		rdsDbName := os.Getenv("RDS_DB_NAME")
		rdsUsername := os.Getenv("RDS_USERNAME")
		rdsPassword := os.Getenv("RDS_PASSWORD")
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", rdsUsername, rdsPassword, rdsHost, rdsPort, rdsDbName)
	}
	if dbURL == "" {
		panic("DB env vars not found")
	}
	var err error
	db, err := gorm.Open("postgres", dbURL)
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	DB.db = db
	fmt.Println("db connected")
}

var updatables = []string{
	"name", "description", "tags", "img",
	"owner_alias", "price_to_join", "price_per_message",
}

// check that update owner_pub_key does in fact throw error
func (db database) createOrEditTribe(m Tribe) (Tribe, error) {
	if m.OwnerPubKey == "" {
		return Tribe{}, errors.New("no pub key")
	}
	onConflict := "ON CONFLICT (uuid) DO UPDATE SET"
	for i, u := range updatables {
		onConflict = onConflict + fmt.Sprintf(" %s=EXCLUDED.%s", u, u)
		if i < len(updatables)-1 {
			onConflict = onConflict + ","
		}
	}
	if m.Name == "" {
		m.Name = "name"
	}
	if m.Description == "" {
		m.Description = "description"
	}
	if m.Tags == nil {
		m.Tags = []string{}
	}
	if err := db.db.Set("gorm:insert_option", onConflict).Create(&m).Error; err != nil {
		fmt.Println(err)
		return Tribe{}, err
	}
	// not working?
	db.db.Exec(`UPDATE tribes SET tsv =
  	setweight(to_tsvector(name), 'A') ||
	setweight(to_tsvector(description), 'B') ||
	setweight(array_to_tsvector(tags), 'C')
	WHERE uuid = '` + m.UUID + "'")
	return m, nil
}

func (db database) getAllTribes() []Tribe {
	ms := []Tribe{}
	db.db.Find(&ms)
	return ms
}

func (db database) getTribe(uuid string) Tribe {
	m := Tribe{}
	db.db.Where("uuid = ?", uuid).Find(&m)
	return m
}

func (db database) searchTribes(s string) []Tribe {
	ms := []Tribe{}
	if s == "" {
		return ms
	}
	// set limit
	db.db.Raw(
		`SELECT uuid, owner_pub_key, name, img, description, ts_rank(tsv, q) as rank
		FROM tribes, to_tsquery('` + s + `') q
		WHERE tsv @@ q
		ORDER BY rank DESC LIMIT 100;`).Find(&ms)
	return ms
}
