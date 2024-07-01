package handlers

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/structs"
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
	keys := r.URL.Query()
	workspace := keys.Get("workspace")

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

	sumAmount := db.DB.TotalPaymentsByDateRange(request, workspace)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sumAmount)
}

func WorkspacetMetrics(w http.ResponseWriter, r *http.Request) {
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

	sumAmount := db.DB.TotalWorkspacesByDateRange(request)

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

	sumAmount := db.DB.TotalWorkspacesByDateRange(request)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sumAmount)
}

func (mh *metricHandler) BountyMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	keys := r.URL.Query()
	workspace := keys.Get("workspace")

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
	if db.RedisError == nil && db.RedisClient != nil {
		redisMetrics := db.GetMap(metricsKey)
		if len(redisMetrics) != 0 {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(redisMetrics)
			return
		}
	} else {
		fmt.Println("Redis client is not initialized or there is an error with Redis")
	}

	totalBountiesPosted := mh.db.TotalBountiesPosted(request, workspace)
	totalBountiesPaid := mh.db.TotalPaidBounties(request, workspace)
	totalBountiesAssigned := mh.db.TotalAssignedBounties(request, workspace)
	bountiesPaidPercentage := mh.db.BountiesPaidPercentage(request, workspace)
	totalSatsPosted := mh.db.TotalSatsPosted(request, workspace)
	totalSatsPaid := mh.db.TotalSatsPaid(request, workspace)
	satsPaidPercentage := mh.db.SatsPaidPercentage(request, workspace)
	avgPaidDays := mh.db.AveragePaidTime(request, workspace)
	avgCompletedDays := mh.db.AverageCompletedTime(request, workspace)
	uniqueHuntersPaid := mh.db.TotalHuntersPaid(request, workspace)
	newHuntersPaid := mh.db.NewHuntersPaid(request, workspace)

	bountyMetrics := db.BountyMetrics{
		BountiesPosted:         totalBountiesPosted,
		BountiesPaid:           totalBountiesPaid,
		BountiesAssigned:       totalBountiesAssigned,
		BountiesPaidPercentage: bountiesPaidPercentage,
		SatsPosted:             totalSatsPosted,
		SatsPaid:               totalSatsPaid,
		SatsPaidPercentage:     satsPaidPercentage,
		AveragePaid:            avgPaidDays,
		AverageCompleted:       avgCompletedDays,
		UniqueHuntersPaid:      uniqueHuntersPaid,
		NewHuntersPaid:         newHuntersPaid,
	}

	if db.RedisError == nil && db.RedisClient != nil {
		metricsMap := structs.Map(bountyMetrics)
		db.SetMap(metricsKey, metricsMap)
	} else {
		fmt.Println("Redis client is not initialized or there is an error with Redis")
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

func (mh *metricHandler) MetricsBountiesProviders(w http.ResponseWriter, r *http.Request) {
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

	bountiesProviders := mh.db.GetBountiesProviders(request, r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountiesProviders)
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

func GetAdminWorkspaces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)

	if pubKeyFromAuth == "" {
		fmt.Println("no pubkey from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workspaces := db.DB.GetWorkspaces(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspaces)
}

func (mh *metricHandler) GetMetricsBountiesData(metricBounties []db.NewBounty) []db.BountyData {
	var metricBountiesData []db.BountyData
	for _, bounty := range metricBounties {
		bountyOwner := mh.db.GetPersonByPubkey(bounty.OwnerID)
		bountyAssignee := mh.db.GetPersonByPubkey(bounty.Assignee)
		workspace := mh.db.GetWorkspaceByUuid(bounty.WorkspaceUuid)

		bountyData := db.BountyData{
			NewBounty:               bounty,
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
			OrganizationName:        workspace.Name,
			OrganizationImg:         workspace.Img,
			OrganizationUuid:        workspace.Uuid,
			OrganizationDescription: workspace.Description,
			WorkspaceName:           workspace.Name,
			WorkspaceImg:            workspace.Img,
			WorkspaceUuid:           workspace.Uuid,
			WorkspaceDescription:    workspace.Description,
		}
		metricBountiesData = append(metricBountiesData, bountyData)
	}
	return metricBountiesData
}

func getMetricsBountyCsv(metricBounties []db.NewBounty) []db.MetricsBountyCsv {
	var metricBountiesCsv []db.MetricsBountyCsv
	for _, bounty := range metricBounties {
		bountyOwner := db.DB.GetPersonByPubkey(bounty.OwnerID)
		bountyAssignee := db.DB.GetPersonByPubkey(bounty.Assignee)
		workspace := db.DB.GetWorkspaceByUuid(bounty.WorkspaceUuid)

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
			Organization: workspace.Name,
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
	opts := &jsonconv.ToCsvOption{
		BaseHeaders: []string{"DatePosted", "Workspace", "BountyAmount", "Provider", "Hunter", "BountyTitle", "BountyLink", "BountyStatus", "DateAssigned", "DatePaid"},
	}
	result := jsonconv.ToCsv(metricsData, opts)
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

	err, postPresignedUrl := createPresignedUrl(path)

	if err != nil {
		fmt.Println("Presigned Error", err)
	}

	r, err := http.NewRequest(http.MethodPut, postPresignedUrl, bytes.NewReader(fileBuffer))
	if err != nil {
		fmt.Println("Posting presign s3 error:", err)
	}
	r.Header.Set("Content-Type", "multipart/form-data")
	client := &http.Client{}
	_, err = client.Do(r)

	if err != nil {
		fmt.Println("Error occured while posting presigned URL", err)
	}

	// Delete image from uploads folder
	DeleteFileFromUploadsFolder(filePath)

	err, presignedUrlGet := getPresignedUrl(path)

	return err, presignedUrlGet
}

func createPresignedUrl(path string) (error, string) {
	presignedUrl, err := config.PresignClient.PresignPutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket: aws.String(config.S3BucketName),
			Key:    aws.String(path),
		},
		s3.WithPresignExpires(time.Minute*15),
	)

	if err != nil {
		return err, ""
	}

	return nil, presignedUrl.URL
}

func getPresignedUrl(path string) (error, string) {
	presignedUrl, err := config.PresignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket:                     aws.String(config.S3BucketName),
			Key:                        aws.String(path),
			ResponseContentDisposition: aws.String("attachment"),
		},
		s3.WithPresignExpires(time.Minute*15),
	)

	if err != nil {
		return err, ""
	}

	return nil, presignedUrl.URL
}
