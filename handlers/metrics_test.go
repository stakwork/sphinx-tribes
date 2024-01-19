package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBountyMetrics(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	mh := NewMetricHandler(mockDb)

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
		db.RedisError = errors.New("redis not initialized")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.BountyMetrics)
		dateRange := db.PaymentDateRange{
			StartDate: "1111",
			EndDate:   "2222",
		}
		body, _ := json.Marshal(dateRange)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/bounty_stats", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		mockDb.On("TotalBountiesPosted", dateRange).Return(int64(1)).Once()
		mockDb.On("TotalPaidBounties", dateRange).Return(int64(1)).Once()
		mockDb.On("BountiesPaidPercentage", dateRange).Return(uint(1)).Once()
		mockDb.On("TotalSatsPosted", dateRange).Return(uint(1)).Once()
		mockDb.On("TotalSatsPaid", dateRange).Return(uint(1)).Once()
		mockDb.On("SatsPaidPercentage", dateRange).Return(uint(1)).Once()
		mockDb.On("AveragePaidTime", dateRange).Return(uint(1)).Once()
		mockDb.On("AverageCompletedTime", dateRange).Return(uint(1)).Once()
		mockDb.On("TotalHuntersPaid", dateRange).Return(int64(1)).Once()
		mockDb.On("NewHuntersPaid", dateRange).Return(int64(1)).Once()
		handler.ServeHTTP(rr, req)

		expectedMetricRes := db.BountyMetrics{
			BountiesPosted:         1,
			BountiesPaid:           1,
			BountiesPaidPercentage: 1,
			SatsPosted:             1,
			SatsPaid:               1,
			SatsPaidPercentage:     1,
			AveragePaid:            1,
			AverageCompleted:       1,
			UniqueHuntersPaid:      1,
			NewHuntersPaid:         1,
		}
		var res db.BountyMetrics
		_ = json.Unmarshal(rr.Body.Bytes(), &res)

		assert.EqualValues(t, expectedMetricRes, res)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestMetricsBounties(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	mh := NewMetricHandler(mockDb)

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

		invalidJson := []byte(`{"key": "value"}`)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/bounties", bytes.NewReader(invalidJson))
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should fetch bounties from db", func(t *testing.T) {
		db.RedisError = errors.New("redis not initialized")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBounties)
		dateRange := db.PaymentDateRange{
			StartDate: "1111",
			EndDate:   "2222",
		}
		body, _ := json.Marshal(dateRange)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/boutnies", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		bounties := []db.Bounty{
			{
				ID:          1,
				OwnerID:     "owner-1",
				Price:       100,
				Title:       "bounty 1",
				Description: "test bounty",
				Created:     1112,
			},
		}
		mockDb.On("GetBountiesByDateRange", dateRange, req).Return(bounties).Once()
		mockDb.On("GetPersonByPubkey", "owner-1").Return(db.Person{ID: 1}).Once()
		mockDb.On("GetPersonByPubkey", "").Return(db.Person{}).Once()
		mockDb.On("GetOrganizationByUuid", "").Return(db.Organization{}).Once()
		handler.ServeHTTP(rr, req)

		var res []db.BountyData
		_ = json.Unmarshal(rr.Body.Bytes(), &res)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, res[0].BountyId, uint(1))
		assert.Equal(t, res[0].OwnerID, "owner-1")
		assert.Equal(t, res[0].Price, uint(100))
		assert.Equal(t, res[0].Title, "bounty 1")
		assert.Equal(t, res[0].BountyDescription, "test bounty")
		assert.Equal(t, res[0].BountyCreated, int64(1112))
	})
}

func TestMetricsBountiesCount(t *testing.T) {
	ctx := context.WithValue(context.Background(), auth.ContextKey, "test-key")
	mockDb := mocks.NewDatabase(t)
	mh := NewMetricHandler(mockDb)

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
		db.RedisError = errors.New("redis not initialized")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(mh.MetricsBountiesCount)
		dateRange := db.PaymentDateRange{
			StartDate: "1111",
			EndDate:   "2222",
		}
		body, _ := json.Marshal(dateRange)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/boutnies/count", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		mockDb.On("GetBountiesByDateRangeCount", dateRange, req).Return(int64(100)).Once()
		handler.ServeHTTP(rr, req)

		var res int64
		_ = json.Unmarshal(rr.Body.Bytes(), &res)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, res, int64(100))
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
		expectedHeaders := []string{"DatePosted", "Organization", "BountyAmount", "Provider", "Hunter", "BountyTitle", "BountyLink", "BountyStatus", "DateAssigned", "DatePaid"}
		results := ConvertMetricsToCSV(bounties)

		assert.Equal(t, 2, len(results))
		assert.EqualValues(t, expectedHeaders, results[0])
	})

}
