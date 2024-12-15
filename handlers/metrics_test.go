package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"fmt"
	"time"

	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/stretchr/testify/assert"
)

func TestBountyMetrics(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	mh := NewMetricHandler(db.TestDB)

	db.TestDB.DeleteAllBounties()

	person := db.Person{
		Uuid:         uuid.New().String(),
		OwnerPubKey:  "person1_pubkey",
		OwnerAlias:   "person1",
		UniqueName:   "person1",
		Description:  "description",
		Tags:         pq.StringArray{},
		Extras:       db.PropertyMap{},
		GithubIssues: db.PropertyMap{},
	}
	db.TestDB.CreateOrEditPerson(person)
	ctx := context.WithValue(context.Background(), auth.ContextKey, person.OwnerPubKey)
	now := time.Now()

	bounty1 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 1",
		Description:   "Bounty ID 1 Description",
		WorkspaceUuid: "workspace",
		Assignee:      "ali",
		OwnerID:       person.OwnerPubKey,
		Show:          true,
		Created:       now.AddDate(0, 0, -30).Unix(),
		Paid:          true,
		Price:         100,
	}
	db.TestDB.CreateOrEditBounty(bounty1)

	bounty2 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 2",
		Description:   "Bounty ID 2 Description",
		WorkspaceUuid: "workspace",
		Assignee:      "raza",
		OwnerID:       person.OwnerPubKey,
		Show:          true,
		Created:       now.AddDate(0, 0, -20).Unix(),
		Paid:          true,
		Price:         150,
	}
	db.TestDB.CreateOrEditBounty(bounty2)

	bounty3 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 3",
		Description:   "Bounty ID 3 Description",
		WorkspaceUuid: "workspace",
		Assignee:      "ali",
		OwnerID:       person.OwnerPubKey,
		Show:          true,
		Created:       now.AddDate(0, 0, -10).Unix(),
		Paid:          false,
	}
	db.TestDB.CreateOrEditBounty(bounty3)

	bounty4 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 4",
		Description:   "Bounty ID 4 Description",
		WorkspaceUuid: "workspace",
		Assignee:      "ali",
		OwnerID:       person.OwnerPubKey,
		Show:          true,
		Created:       now.Unix(),
		Paid:          false,
	}
	db.TestDB.CreateOrEditBounty(bounty4)

	dateRange := db.PaymentDateRange{
		StartDate: strconv.FormatInt(bounty1.Created, 10),
		EndDate:   strconv.FormatInt(bounty4.Created, 10),
	}

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.BountyMetrics)

		invalidJson := []byte(`{"key": "value"`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/bounty_stats", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should return error if public key not present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.BountyMetrics)

		invalidJson := []byte(`{"key": "value"}`)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/bounty_stats", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should fetch stats from db", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.BountyMetrics)

		body, _ := json.Marshal(dateRange)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/bounty_stats?workspace=workspace", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		totalBountiesPosted := int64(4)
		totalBountiesPaid := int64(2)
		totalBountiesAssigned := int64(2)
		bountiesPaidPercentage := uint((totalBountiesPaid * 100) / totalBountiesPosted)
		totalSatsPosted := bounty1.Price + bounty2.Price + bounty3.Price + bounty4.Price
		totalSatsPaid := bounty1.Price + bounty2.Price
		satsPaidPercentage := uint((totalSatsPaid * 100) / totalSatsPosted)
		avgPaidDays := uint(0)
		avgCompletedDays := uint(0)
		uniqueHuntersPaid := int64(2)
		newHuntersPaid := int64(2)
		newHunters := db.TestDB.GetNewHunters(dateRange)
		peopleByPeriod := db.TestDB.TotalPeopleByPeriod(dateRange)

		expectedMetricRes := db.BountyMetrics{
			BountiesPosted:         totalBountiesPosted,
			BountiesPaid:           int64(totalBountiesPaid),
			BountiesAssigned:       int64(totalBountiesAssigned),
			BountiesPaidPercentage: bountiesPaidPercentage,
			SatsPosted:             uint(totalSatsPosted),
			SatsPaid:               uint(totalSatsPaid),
			SatsPaidPercentage:     satsPaidPercentage,
			AveragePaid:            avgPaidDays,
			AverageCompleted:       avgCompletedDays,
			UniqueHuntersPaid:      uniqueHuntersPaid,
			NewHuntersPaid:         newHuntersPaid,
			NewHunters:             newHunters,
			NewHuntersByPeriod:     peopleByPeriod,
		}

		var res db.BountyMetrics
		_ = json.Unmarshal(rr.Body.Bytes(), &res)

		assert.EqualValues(t, expectedMetricRes, res)
	})
}

func TestMetricsBounties(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mh := NewMetricHandler(db.TestDB)

	now := time.Now()
	bountyOwner := db.Person{OwnerPubKey: "owner-1"}
	db.TestDB.CreateOrEditPerson(bountyOwner)

	bounty1 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 1",
		Description:   "Bounty ID 1 Description",
		WorkspaceUuid: "",
		Assignee:      "",
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
		Created:       now.AddDate(0, 0, -30).Unix(),
		Paid:          true,
	}
	db.TestDB.CreateOrEditBounty(bounty1)

	bounty2 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 2",
		Description:   "Bounty ID 2 Description",
		WorkspaceUuid: "",
		Assignee:      "",
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
		Created:       now.AddDate(0, 0, -20).Unix(),
		Paid:          true,
	}
	db.TestDB.CreateOrEditBounty(bounty2)

	bounty3 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 3",
		Description:   "Bounty ID 3 Description",
		WorkspaceUuid: "",
		Assignee:      "",
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
		Created:       now.AddDate(0, 0, -10).Unix(),
		Paid:          false,
	}
	db.TestDB.CreateOrEditBounty(bounty3)

	bounty4 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 4",
		Description:   "Bounty ID 4 Description",
		WorkspaceUuid: "",
		Assignee:      "",
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
		Created:       now.Unix(),
		Paid:          false,
	}
	db.TestDB.CreateOrEditBounty(bounty4)

	dateRange := db.PaymentDateRange{
		StartDate: strconv.FormatInt(bounty1.Created, 10),
		EndDate:   strconv.FormatInt(bounty4.Created, 10),
	}

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBounties)

		invalidJson := []byte(`{"key": "value"`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/bounties", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should return error if public key not present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBounties)

		invalidJson := []byte(`{"key": "value"`)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/bounties", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should fetch bounties from db", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBounties)

		body, _ := json.Marshal(dateRange)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/bounties", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var res []db.BountyData
		_ = json.Unmarshal(rr.Body.Bytes(), &res)

		bounties := db.TestDB.GetBountiesByDateRange(dateRange, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, len(bounties), len(res))

		for i, bounty := range bounties {
			assert.Equal(t, bounty.ID, res[i].BountyId)
			assert.Equal(t, bounty.OwnerID, res[i].OwnerID)
			assert.Equal(t, bounty.Price, res[i].Price)
			assert.Equal(t, bounty.Title, res[i].Title)
			assert.Equal(t, bounty.Description, res[i].BountyDescription)
			assert.Equal(t, bounty.Created, res[i].BountyCreated)
		}
	})

	t.Run("should fetch bounties from db for selected providers", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBounties)

		body, _ := json.Marshal(dateRange)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/bounties", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		// Provide multiple provider IDs in the request query parameters
		req.URL.RawQuery = "provider=owner-1"

		handler.ServeHTTP(rr, req)

		var res []db.BountyData
		_ = json.Unmarshal(rr.Body.Bytes(), &res)

		bounties := db.TestDB.GetBountiesByDateRange(dateRange, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, len(bounties), len(res))

		// Assert that the response contains bounties only from the selected providers
		for _, bounty := range res {
			assert.Equal(t, "owner-1", bounty.OwnerID)
		}
	})
}

func TestMetricsBountiesCount(t *testing.T) {
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)

	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mh := NewMetricHandler(db.TestDB)

	now := time.Now()
	bountyOwner := db.Person{OwnerPubKey: "owner-1"}
	db.TestDB.CreateOrEditPerson(bountyOwner)

	bounty1 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 1",
		Description:   "Bounty ID 1 Description",
		WorkspaceUuid: "",
		Assignee:      "",
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
		Created:       now.AddDate(0, 0, -30).Unix(),
		Paid:          true,
	}
	db.TestDB.CreateOrEditBounty(bounty1)

	bounty2 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 2",
		Description:   "Bounty ID 2 Description",
		WorkspaceUuid: "",
		Assignee:      "",
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
		Created:       now.AddDate(0, 0, -20).Unix(),
		Paid:          true,
	}
	db.TestDB.CreateOrEditBounty(bounty2)

	bounty3 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 3",
		Description:   "Bounty ID 3 Description",
		WorkspaceUuid: "",
		Assignee:      "",
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
		Created:       now.AddDate(0, 0, -10).Unix(),
		Paid:          false,
	}
	db.TestDB.CreateOrEditBounty(bounty3)

	bounty4 := db.NewBounty{
		Type:          "coding",
		Title:         "Bounty With ID 4",
		Description:   "Bounty ID 4 Description",
		WorkspaceUuid: "",
		Assignee:      "",
		OwnerID:       bountyOwner.OwnerPubKey,
		Show:          true,
		Created:       now.Unix(),
		Paid:          false,
	}
	db.TestDB.CreateOrEditBounty(bounty4)

	dateRange := db.PaymentDateRange{
		StartDate: strconv.FormatInt(bounty1.Created, 10),
		EndDate:   strconv.FormatInt(bounty4.Created, 10),
	}

	t.Run("should return error if body is not a valid json", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBountiesCount)

		invalidJson := []byte(`{"key": "value"`)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/bounties/count", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should return error if public key not present", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBountiesCount)

		invalidJson := []byte(`{"key": "value"}`)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/bounties/count", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should fetch bounties count within specified date range", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBountiesCount)

		body, _ := json.Marshal(dateRange)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/bounties/count", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		var res int64
		_ = json.Unmarshal(rr.Body.Bytes(), &res)

		expectedCount := db.TestDB.GetBountiesByDateRangeCount(dateRange, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, expectedCount, res)
	})

	t.Run("should fetch bounties count within specified date range for selected providers", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBountiesCount)

		body, _ := json.Marshal(dateRange)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/bounties/count", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		// Provide provider IDs in the request query parameters
		req.URL.RawQuery = "provider=owner-1"

		handler.ServeHTTP(rr, req)

		var res int64
		_ = json.Unmarshal(rr.Body.Bytes(), &res)

		expectedCount := db.TestDB.GetBountiesByDateRangeCount(dateRange, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, expectedCount, res)
	})
}

func TestConvertMetricsToCSV(t *testing.T) {
	t.Run("should return for csv in correct order", func(t *testing.T) {
		now := time.Now()
		bountyLink := fmt.Sprintf("https://community.sphinx.chat/bounty/%d", 1)
		bounties := []db.MetricsBountyCsv{{
			DatePosted:   &now,
			Organization: "test-org",
			BountyAmount: 100,
			Provider:     "provider",
			Hunter:       "hunter",
			BountyTitle:  "test bounty",
			BountyLink:   bountyLink,
			BountyStatus: "paid",
			DatePaid:     &now,
			DateAssigned: &now,
		}}
		expectedHeaders := []string{"DatePosted", "Workspace", "BountyAmount", "Provider", "Hunter", "BountyTitle", "BountyLink", "BountyStatus", "DateAssigned", "DatePaid"}
		results := ConvertMetricsToCSV(bounties)

		assert.Equal(t, 2, len(results))
		assert.EqualValues(t, expectedHeaders, results[0])
	})

}

func TestMetricsBountiesProviders(t *testing.T) {
	ctx := context.Background()
	teardownSuite := SetupSuite(t)
	defer teardownSuite(t)
	mh := NewMetricHandler(db.TestDB)
	unauthorizedCtx := context.WithValue(context.Background(), auth.ContextKey, "")
	authorizedCtx := context.WithValue(ctx, auth.ContextKey, "valid-key")

	t.Run("should return 401 error if user is unauthorized", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBountiesProviders)

		req, err := http.NewRequestWithContext(unauthorizedCtx, http.MethodPost, "/bounties/providers", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return 406 error if wrong data is passed", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBountiesProviders)

		invalidJson := []byte(`{"start_date": "2021-01-01"`)
		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/bounties/providers", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotAcceptable, rr.Code)
	})

	t.Run("should return bounty providers and 200 status code if there is no error", func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBountiesProviders)

		db.TestDB.DeleteAllBounties()

		person1 := db.Person{
			Uuid:         uuid.New().String(),
			OwnerPubKey:  "person1_pubkey",
			OwnerAlias:   "person1",
			UniqueName:   "person1",
			Description:  "description",
			Tags:         pq.StringArray{},
			Extras:       db.PropertyMap{},
			GithubIssues: db.PropertyMap{},
		}
		person2 := db.Person{
			Uuid:         uuid.New().String(),
			OwnerPubKey:  "person2_pubkey",
			OwnerAlias:   "person2",
			UniqueName:   "person2",
			Description:  "description",
			Tags:         pq.StringArray{},
			Extras:       db.PropertyMap{},
			GithubIssues: db.PropertyMap{},
		}

		db.TestDB.CreateOrEditPerson(person1)
		db.TestDB.CreateOrEditPerson(person2)

		now := time.Now()
		thirtyDaysBefore := now.AddDate(0, 0, -30).Unix()
		twentyDaysBefore := now.AddDate(0, 0, -20).Unix()
		tenDaysBefore := now.AddDate(0, 0, -10).Unix()
		nowUnix := now.Unix()

		bounty1 := db.NewBounty{
			Type:          "coding",
			Title:         "Bounty With ID",
			Description:   "Bounty ID Description",
			WorkspaceUuid: "",
			Assignee:      "",
			Show:          true,
			OwnerID:       person2.OwnerPubKey,
			Paid:          true,
			Created:       thirtyDaysBefore,
		}
		bounty2 := db.NewBounty{
			Type:          "coding",
			Title:         "Bounty With ID",
			Description:   "Bounty ID Description",
			WorkspaceUuid: "",
			Assignee:      "",
			Show:          true,
			OwnerID:       person2.OwnerPubKey,
			Created:       twentyDaysBefore,
		}
		bounty3 := db.NewBounty{
			Type:          "coding",
			Title:         "Bounty With ID",
			Description:   "Bounty ID Description",
			WorkspaceUuid: "",
			Assignee:      "",
			Show:          true,
			OwnerID:       person1.OwnerPubKey,
			Paid:          true,
			Created:       tenDaysBefore,
		}
		bounty4 := db.NewBounty{
			Type:          "coding",
			Title:         "Bounty With ID",
			Description:   "Bounty ID Description",
			WorkspaceUuid: "",
			Assignee:      "",
			Show:          true,
			OwnerID:       person1.OwnerPubKey,
			Created:       nowUnix,
		}

		db.TestDB.CreateOrEditBounty(bounty1)
		db.TestDB.CreateOrEditBounty(bounty2)
		db.TestDB.CreateOrEditBounty(bounty3)
		db.TestDB.CreateOrEditBounty(bounty4)

		dateRange := db.PaymentDateRange{
			StartDate: strconv.FormatInt(bounty1.Created, 10),
			EndDate:   strconv.FormatInt(bounty4.Created, 10),
		}

		body, _ := json.Marshal(dateRange)

		req, err := http.NewRequestWithContext(authorizedCtx, http.MethodPost, "/bounties/providers?limit=10", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		fetchedProviders := db.TestDB.GetBountiesProviders(dateRange, req)

		handler.ServeHTTP(rr, req)

		var actualProviders []db.Person
		err = json.Unmarshal(rr.Body.Bytes(), &actualProviders)
		if err != nil {
			t.Fatal("Failed to unmarshal response:", err)
		}

		assert.Equal(t, http.StatusOK, rr.Code)
		//Assert that the API call response matches the value returned from the DB
		assert.EqualValues(t, fetchedProviders, actualProviders)
		//Assert that the Providers returned are equal to the persons created
		person1.ID = fetchedProviders[0].ID
		person2.ID = fetchedProviders[1].ID
		expectedProviders := []db.Person{person1, person2}
		assert.EqualValues(t, expectedProviders, actualProviders)
	})

}
