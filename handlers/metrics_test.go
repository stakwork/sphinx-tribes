package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	mocks "github.com/stakwork/sphinx-tribes/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
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
