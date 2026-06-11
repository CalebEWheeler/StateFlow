package operations

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var CreateUser = huma.Operation{
	Description:   "Create a new platform user",
	Method:        http.MethodPost,
	OperationID:   "create-user",
	Path:          "/user",
	Summary:       "Create user",
	Tags:          []string{"Create User"},
	DefaultStatus: 200,
}
