package operations

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var CreateBilling = huma.Operation{
	Description: "Create new user billing details",
	Method:      http.MethodPost,
	OperationID: "create-billing",
	Path:        "/create/billing",
	Summary:     "Create billing",
	Tags:        []string{"Create Billing"},
}
