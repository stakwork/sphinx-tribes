package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/stakwork/sphinx-tribes/db"
)

func GetAllBounties(w http.ResponseWriter, r *http.Request) {
	bounties := db.DB.GetAllBounties(r)
	var bountyResponse []db.BountyResponse = generateBountyResponse(bounties)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyResponse)
}

func GetBountyById(w http.ResponseWriter, r *http.Request) {
	bountyId := chi.URLParam(r, "bountyId")
	if bountyId == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bounties, err := db.DB.GetBountyById(bountyId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error", err)
	} else {
		var bountyResponse []db.BountyResponse = generateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func GetBountyCount(w http.ResponseWriter, r *http.Request) {
	personKey := chi.URLParam(r, "personKey")
	tabType := chi.URLParam(r, "tabType")

	if personKey == "" || tabType == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bountyCount := db.DB.GetBountiesCounty(personKey, tabType)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bountyCount)
}

func GetPersonCreatedWanteds(w http.ResponseWriter, r *http.Request) {
	pubkey := chi.URLParam(r, "pubkey")
	if pubkey == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bounties, err := db.DB.GetCreatedBounties(pubkey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error", err)
	} else {
		var bountyResponse []db.BountyResponse = generateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func GetPersonAssignedWanteds(w http.ResponseWriter, r *http.Request) {
	pubkey := chi.URLParam(r, "pubkey")
	if pubkey == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	bounties, err := db.DB.GetAssignedBounties(pubkey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Error", err)
	} else {
		var bountyResponse []db.BountyResponse = generateBountyResponse(bounties)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bountyResponse)
	}
}

func CreateOrEditBounty(w http.ResponseWriter, r *http.Request) {
	bounty := db.Bounty{}
	body, err := ioutil.ReadAll(r.Body)

	r.Body.Close()
	err = json.Unmarshal(body, &bounty)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	now := time.Now()

	//Check if bounty exists
	bounty.Updated = &now
	if bounty.Created == 0 {
		bounty.Created = time.Now().Unix()
	}

	if bounty.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Type is a required field")
		return
	}
	if bounty.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Title is a required field")
		return
	}
	if bounty.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Description is a required field")
		return
	}
	if bounty.Tribe == "" {
		bounty.Tribe = "None"
	}

	b, err := db.DB.CreateOrEditBounty(bounty)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)
}

func DeleteBounty(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	//pubKeyFromAuth, _ := ctx.Value(auth.ContextKey).(string)
	created := chi.URLParam(r, "created")
	pubkey := chi.URLParam(r, "pubkey")

	if pubkey == "" {
		fmt.Println("no pubkey from route")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if created == "" {
		fmt.Println("no created timestamp from route")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	b, _ := db.DB.DeleteBounty(pubkey, created)
	json.NewEncoder(w).Encode(b)
}

func generateBountyResponse(bounties []db.BountyData) []db.BountyResponse {
	var bountyResponse []db.BountyResponse

	for i := 0; i < len(bounties); i++ {
		bounty := bounties[i]
		b := db.BountyResponse{
			Bounty: db.Bounty{
				ID:                      bounty.BountyId,
				OwnerID:                 bounty.OwnerID,
				Paid:                    bounty.Paid,
				Show:                    bounty.Show,
				Type:                    bounty.Type,
				Award:                   bounty.Award,
				AssignedHours:           bounty.AssignedHours,
				BountyExpires:           bounty.BountyExpires,
				CommitmentFee:           bounty.CommitmentFee,
				Price:                   bounty.Price,
				Title:                   bounty.Title,
				Tribe:                   bounty.Tribe,
				Created:                 bounty.BountyCreated,
				Assignee:                bounty.Assignee,
				TicketUrl:               bounty.TicketUrl,
				Description:             bounty.BountyDescription,
				WantedType:              bounty.WantedType,
				Deliverables:            bounty.Deliverables,
				GithubDescription:       bounty.GithubDescription,
				OneSentenceSummary:      bounty.OneSentenceSummary,
				EstimatedSessionLength:  bounty.EstimatedSessionLength,
				EstimatedCompletionDate: bounty.EstimatedCompletionDate,
				Organization:            bounty.Organization,
				Updated:                 bounty.BountyUpdated,
				CodingLanguages:         bounty.CodingLanguages,
			},
			Assignee: db.Person{
				ID:               bounty.AssigneeId,
				Uuid:             bounty.Uuid,
				OwnerPubKey:      bounty.OwnerPubKey,
				OwnerAlias:       bounty.AssigneeAlias,
				UniqueName:       bounty.UniqueName,
				Description:      bounty.AssigneeDescription,
				Tags:             bounty.Tags,
				Img:              bounty.Img,
				Created:          bounty.AssigneeCreated,
				Updated:          bounty.AssigneeUpdated,
				LastLogin:        bounty.LastLogin,
				OwnerRouteHint:   bounty.OwnerRouteHint,
				OwnerContactKey:  bounty.OwnerContactKey,
				PriceToMeet:      bounty.PriceToMeet,
				TwitterConfirmed: bounty.TwitterConfirmed,
			},
			Owner: db.Person{
				ID:               bounty.BountyOwnerId,
				Uuid:             bounty.OwnerUuid,
				OwnerPubKey:      bounty.OwnerKey,
				OwnerAlias:       bounty.OwnerAlias,
				UniqueName:       bounty.OwnerUniqueName,
				Description:      bounty.OwnerDescription,
				Tags:             bounty.OwnerTags,
				Img:              bounty.OwnerImg,
				Created:          bounty.OwnerCreated,
				Updated:          bounty.OwnerUpdated,
				LastLogin:        bounty.OwnerLastLogin,
				OwnerRouteHint:   bounty.OwnerRouteHint,
				OwnerContactKey:  bounty.OwnerContactKey,
				PriceToMeet:      bounty.OwnerPriceToMeet,
				TwitterConfirmed: bounty.OwnerTwitterConfirmed,
			},
			Organization: db.OrganizationSHort{
				Name: bounty.OrganizationName,
				Uuid: bounty.OrganizationUuid,
				Img:  bounty.OrganizationImg,
			},
		}
		bountyResponse = append(bountyResponse, b)
	}

	return bountyResponse
}
