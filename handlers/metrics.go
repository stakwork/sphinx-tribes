package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ambelovsky/go-structs"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/tuan78/jsonconv"
)

func PaymentMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := db.PaymentDateRange{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode("Request body not accepted")
		return
	}

	sumAmount := db.DB.TotalPaymentsByDateRange(request)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sumAmount)
}

func OrganizationtMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := db.PaymentDateRange{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode("Request body not accepted")
		return
	}

	sumAmount := db.DB.TotalOrganizationsByDateRange(request)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sumAmount)
}

func PeopleMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := db.PaymentDateRange{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode("Request body not accepted")
		return
	}

	sumAmount := db.DB.TotalOrganizationsByDateRange(request)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sumAmount)
}

func BountyMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := db.PaymentDateRange{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode("Request body not accepted")
		return
	}

	metricsKey := fmt.Sprintf("metrics - %s - %s", request.StartDate, request.EndDate)
	/**
	check redis if cache id available for the date range
	or add to redis
	*/
	if db.RedisError == nil {
		redisMetrics := db.GetMap(metricsKey)
		if len(redisMetrics) != 0 {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(redisMetrics)
			return
		}
	}

	totalBountiesPosted := db.DB.TotalBountiesPosted(request)
	totalBountiesPaid := db.DB.TotalPaidBounties(request)
	bountiesPaidPercentage := db.DB.BountiesPaidPercentage(request)
	totalSatsPosted := db.DB.TotalSatsPosted(request)
	totalSatsPaid := db.DB.TotalSatsPaid(request)
	satsPaidPercentage := db.DB.SatsPaidPercentage(request)
	avgPaidDays := db.DB.AveragePaidTime(request)
	avgCompletedDays := db.DB.AverageCompletedTime(request)

	bountyMetrics := db.BountyMetrics{
		BountiesPosted:         totalBountiesPosted,
		BountiesPaid:           totalBountiesPaid,
		BountiesPaidPercentage: bountiesPaidPercentage,
		SatsPosted:             totalSatsPosted,
		SatsPaid:               totalSatsPaid,
		SatsPaidPercentage:     satsPaidPercentage,
		AveragePaid:            avgPaidDays,
		AverageCompleted:       avgCompletedDays,
	}

	if db.RedisError == nil {
		metricsMap := structs.Map(bountyMetrics)
		db.SetMap(metricsKey, metricsMap)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyMetrics)
}

func MetricsBounties(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := db.PaymentDateRange{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode("Request body not accepted")
		return
	}

	metricBounties := db.DB.GetBountiesByDateRange(request, r)
	metricBountiesData := GetMetricsBountiesData(metricBounties)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metricBountiesData)
}

func MetricsBountiesCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := db.PaymentDateRange{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode("Request body not accepted")
		return
	}

	MetricsBountiesCount := db.DB.GetBountiesByDateRangeCount(request, r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MetricsBountiesCount)
}

func MetricsCsv(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request := db.PaymentDateRange{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(body, &request)

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode("Request body not accepted")
		return
	}

	metricBounties := db.DB.GetBountiesByDateRange(request, r)
	metricBountiesData := GetMetricsBountiesData(metricBounties)
	result := ConvertMetricsToCSV(metricBountiesData)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func GetMetricsBountiesData(metricBounties []db.Bounty) []db.BountyData {
	var metricBountiesData []db.BountyData
	for _, bounty := range metricBounties {
		bountyOwner := db.DB.GetPersonByPubkey(bounty.OwnerID)
		bountyAssignee := db.DB.GetPersonByPubkey(bounty.Assignee)
		organization := db.DB.GetOrganizationByUuid(bounty.OrgUuid)

		bountyData := db.BountyData{
			Bounty:              bounty,
			BountyId:            bounty.ID,
			Person:              bountyOwner,
			BountyCreated:       bounty.Created,
			BountyDescription:   bounty.Description,
			BountyUpdated:       bounty.Updated,
			AssigneeId:          bountyAssignee.ID,
			AssigneeImg:         bountyAssignee.Img,
			AssigneeAlias:       bountyAssignee.OwnerAlias,
			AssigneeDescription: bountyAssignee.Description,
			AssigneeRouteHint:   bountyAssignee.OwnerRouteHint,
			BountyOwnerId:       bountyOwner.ID,
			OwnerUuid:           bountyOwner.Uuid,
			OwnerDescription:    bountyOwner.Description,
			OwnerUniqueName:     bountyOwner.UniqueName,
			OwnerImg:            bountyOwner.Img,
			OrganizationName:    organization.Name,
			OrganizationImg:     organization.Img,
			OrganizationUuid:    organization.Uuid,
		}
		metricBountiesData = append(metricBountiesData, bountyData)
	}
	return metricBountiesData
}

func ConvertMetricsToCSV(metricBountiesData []db.BountyData) [][]string {
	var metricsData []map[string]interface{}
	data, err := json.Marshal(metricBountiesData)
	if err != nil {
		fmt.Println("Could not convert metrics structs Array to JSON")
		return [][]string{}
	}
	err = json.Unmarshal(data, &metricsData)
	result := jsonconv.ToCsv(metricsData, nil)
	return result
}
