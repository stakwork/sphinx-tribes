package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/ambelovsky/go-structs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stakwork/sphinx-tribes/auth"
	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/db"
	"github.com/tuan78/jsonconv"
)

type metricHandler struct {
	db db.Database
}

func NewMetricHandler(db db.Database) *metricHandler {
	return &metricHandler{db: db}
}

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

func (mh *metricHandler) BountyMetrics(w http.ResponseWriter, r *http.Request) {
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

	totalBountiesPosted := mh.db.TotalBountiesPosted(request)
	totalBountiesPaid := mh.db.TotalPaidBounties(request)
	bountiesPaidPercentage := mh.db.BountiesPaidPercentage(request)
	totalSatsPosted := mh.db.TotalSatsPosted(request)
	totalSatsPaid := mh.db.TotalSatsPaid(request)
	satsPaidPercentage := mh.db.SatsPaidPercentage(request)
	avgPaidDays := mh.db.AveragePaidTime(request)
	avgCompletedDays := mh.db.AverageCompletedTime(request)
	uniqueHuntersPaid := mh.db.TotalHuntersPaid(request)
	newHuntersPaid := mh.db.NewHuntersPaid(request)

	bountyMetrics := db.BountyMetrics{
		BountiesPosted:         totalBountiesPosted,
		BountiesPaid:           totalBountiesPaid,
		BountiesPaidPercentage: bountiesPaidPercentage,
		SatsPosted:             totalSatsPosted,
		SatsPaid:               totalSatsPaid,
		SatsPaidPercentage:     satsPaidPercentage,
		AveragePaid:            avgPaidDays,
		AverageCompleted:       avgCompletedDays,
		UniqueHuntersPaid:      uniqueHuntersPaid,
		NewHuntersPaid:         newHuntersPaid,
	}

	if db.RedisError == nil {
		metricsMap := structs.Map(bountyMetrics)
		db.SetMap(metricsKey, metricsMap)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyMetrics)
}

func (mh *metricHandler) MetricsBounties(w http.ResponseWriter, r *http.Request) {
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

	metricBounties := mh.db.GetBountiesByDateRange(request, r)
	metricBountiesData := mh.GetMetricsBountiesData(metricBounties)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metricBountiesData)
}

func (mh *metricHandler) MetricsBountiesCount(w http.ResponseWriter, r *http.Request) {
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

	MetricsBountiesCount := mh.db.GetBountiesByDateRangeCount(request, r)
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
	metricsCsv := getMetricsBountyCsv(metricBounties)
	result := ConvertMetricsToCSV(metricsCsv)
	resultLength := len(result)

	if resultLength > 0 {
		err, url := UploadMetricsCsv(result, request)

		if err != nil {
			fmt.Println("Error uploading csv ===", err)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(url)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (mh *metricHandler) GetMetricsBountiesData(metricBounties []db.Bounty) []db.BountyData {
	var metricBountiesData []db.BountyData
	for _, bounty := range metricBounties {
		bountyOwner := mh.db.GetPersonByPubkey(bounty.OwnerID)
		bountyAssignee := mh.db.GetPersonByPubkey(bounty.Assignee)
		organization := mh.db.GetOrganizationByUuid(bounty.OrgUuid)

		bountyData := db.BountyData{
			Bounty:                  bounty,
			BountyId:                bounty.ID,
			Person:                  bountyOwner,
			BountyCreated:           bounty.Created,
			BountyDescription:       bounty.Description,
			BountyUpdated:           bounty.Updated,
			AssigneeId:              bountyAssignee.ID,
			AssigneeImg:             bountyAssignee.Img,
			AssigneeAlias:           bountyAssignee.OwnerAlias,
			AssigneeDescription:     bountyAssignee.Description,
			AssigneeRouteHint:       bountyAssignee.OwnerRouteHint,
			BountyOwnerId:           bountyOwner.ID,
			OwnerUuid:               bountyOwner.Uuid,
			OwnerDescription:        bountyOwner.Description,
			OwnerUniqueName:         bountyOwner.UniqueName,
			OwnerImg:                bountyOwner.Img,
			OrganizationName:        organization.Name,
			OrganizationImg:         organization.Img,
			OrganizationUuid:        organization.Uuid,
			OrganizationDescription: organization.Description,
		}
		metricBountiesData = append(metricBountiesData, bountyData)
	}
	return metricBountiesData
}

func getMetricsBountyCsv(metricBounties []db.Bounty) []db.MetricsBountyCsv {
	var metricBountiesCsv []db.MetricsBountyCsv
	for _, bounty := range metricBounties {
		bountyOwner := db.DB.GetPersonByPubkey(bounty.OwnerID)
		bountyAssignee := db.DB.GetPersonByPubkey(bounty.Assignee)
		organization := db.DB.GetOrganizationByUuid(bounty.OrgUuid)

		bountyLink := fmt.Sprintf("https://community.sphinx.chat/bounty/%d", bounty.ID)
		bountyStatus := "Open"

		if bounty.Assignee != "" && !bounty.Paid {
			bountyStatus = "Assigned"
		} else {
			bountyStatus = "Paid"
		}

		tm := time.Unix(bounty.Created, 0)
		bountyCsv := db.MetricsBountyCsv{
			DatePosted:   &tm,
			Organization: organization.Name,
			BountyAmount: bounty.Price,
			Provider:     bountyOwner.OwnerAlias,
			Hunter:       bountyAssignee.OwnerAlias,
			BountyTitle:  bounty.Title,
			BountyLink:   bountyLink,
			BountyStatus: bountyStatus,
			DateAssigned: bounty.AssignedDate,
			DatePaid:     bounty.PaidDate,
		}
		metricBountiesCsv = append(metricBountiesCsv, bountyCsv)
	}

	return metricBountiesCsv
}

func ConvertMetricsToCSV(metricBountiesData []db.MetricsBountyCsv) [][]string {
	metricsData := db.DB.ConvertMetricsBountiesToMap(metricBountiesData)
	result := jsonconv.ToCsv(metricsData, nil)
	return result
}

func UploadMetricsCsv(data [][]string, request db.PaymentDateRange) (error, string) {
	dirName := "uploads"
	CreateUploadsDirectory(dirName)

	filePath := path.Join("./uploads", "metrics.csv")
	csvFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	w := csv.NewWriter(csvFile)
	err = w.WriteAll(data)
	if err != nil {
		log.Fatal(err)
	}

	csvFile.Close()

	upFile, _ := os.Open(filePath)
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	key := fmt.Sprintf("metrics%s-%s.csv", request.StartDate, request.EndDate)
	path := fmt.Sprintf("%s/%s", config.S3FolderName, key)
	_, err = config.S3Client.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(config.S3BucketName),
		Key:                  aws.String(path),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String("application/csv"),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	url := fmt.Sprintf("%s/%s/%s", config.S3Url, config.S3FolderName, key)

	// Delete image from uploads folder
	DeleteFileFromUploadsFolder(filePath)

	return err, url
}
