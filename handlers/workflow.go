package handlers

import (
	"fmt"
	"github.com/stakwork/sphinx-tribes/db"
	"io"
	"net/http"
)

type workflowHandler struct {
	db db.Database
}

func NewWorkFlowHandler(database db.Database) *workflowHandler {
	return &workflowHandler{
		db: database,
	}
}

func (oh *workflowHandler) HandleWorkflowRequest(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Request Body:", string(body))

}
