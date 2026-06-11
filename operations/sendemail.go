package operations

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var SendEmail = huma.Operation{
	Description:   "Send an email to the created user",
	Method:        http.MethodPost,
	OperationID:   "send-email",
	Path:          "/email",
	Summary:       "Send Email",
	Tags:          []string{"Send Email"},
	DefaultStatus: 200,
}
