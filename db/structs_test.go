package db

import (
	"sync"
	"testing"

	"gorm.io/gorm/schema"
)

func TestBountySummaryIsNotPersisted(t *testing.T) {
	models := []struct {
		name  string
		value interface{}
	}{
		{name: "Bounty", value: &Bounty{}},
		{name: "NewBounty", value: &NewBounty{}},
	}

	for _, model := range models {
		t.Run(model.name, func(t *testing.T) {
			parsed, err := schema.Parse(model.value, &sync.Map{}, schema.NamingStrategy{})
			if err != nil {
				t.Fatalf("parse schema: %v", err)
			}

			if field := parsed.LookUpField("OneSentenceSummary"); field == nil || field.DBName != "" {
				t.Fatalf("OneSentenceSummary should not be mapped to a database column, got %#v", field)
			}
			if field := parsed.LookUpField("one_sentence_summary"); field != nil {
				t.Fatalf("one_sentence_summary should not be mapped to database field %q", field.Name)
			}
		})
	}
}
