package db

type V2InvoiceResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	AmtMsat   string `json:"amt_msat"`
}

type V2InvoiceBody struct {
	PaymentHash string `json:"payment_hash"`
}

type V2SendOnionRes struct {
	Status      string `json:"status"` // "COMPLETE", "PENDING", or "FAILED"
	Tag         string `json:"tag"`
	Preimage    string `json:"preimage"`
	PaymentHash string `json:"payment_hash"`
	Message     string `json:"message"`
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
