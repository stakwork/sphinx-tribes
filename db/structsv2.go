package db

const (
	InvoicePaid    = "paid"
	InvoiceExpired = "expired"
	InvoicePending = "pending"
)

const (
	PaymentComplete = "COMPLETE"
	PaymentFailed   = "FAILED"
	PaymentPending  = "PENDING"
	PaymentNotFound = "NOTPAID"
)

type V2InvoiceResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	AmtMsat   string `json:"amt_msat"`
}

type V2InvoiceBody struct {
	PaymentHash string `json:"payment_hash"`
	Bolt11      string `json:"bolt11"`
}

type V2SendOnionRes struct {
	Status      string `json:"status"` // "COMPLETE", "PENDING", or "FAILED"
	Tag         string `json:"tag"`
	Preimage    string `json:"preimage"`
	PaymentHash string `json:"payment_hash"`
	Message     string `json:"message,omitempty"`
}

type V2PayInvoiceBody struct {
	Bolt11 string `json:"bolt11"`
}

type V2CreateInvoiceBody struct {
	AmtMsat uint `json:"amt_msat"`
}

type V2CreateInvoiceResponse struct {
	Bolt11      string `json:"bolt11"`
	PaymentHash string `json:"payment_hash"`
}

type V2PayInvoiceResponse struct {
	Tag         string `json:"tag"`
	Msat        string `json:"msat"`
	Timestamp   string `json:"timestamp"`
	PaymentHash string `json:"payment_hash"`
}

type V2TagRes struct {
	Tag    string `json:"tag"`
	Ts     uint64 `json:"ts"`
	Status string `json:"status"` // "COMPLETE", "PENDING", or "FAILED"
	Error  string `json:"error"`
}

type FeatureStories struct {
	UserStory string `json:"userStory"`
	Rationale string `json:"rationale"`
	Order     uint   `json:"order"`
}

type FeatureOutput struct {
	FeatureContext string           `json:"featureContext"`
	FeatureUuid    string           `json:"featureUuid"`
	Stories        []FeatureStories `json:"stories"`
}

type FeatureStoriesReponse struct {
	Output FeatureOutput `json:"output"`
}

type InviteReponse struct {
	Invite string `json:"invite"`
}

type InviteBody struct {
	Number     uint   `json:"number"`
	Pubkey     string `json:"pubkey"`
	RouteHint  string `json:"route_hint"`
	SatsAmount uint64 `json:"sats_amount"`
}
