package operations

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var Order = huma.Operation{
	Description:   "Store order details and initiate order processing workflow",
	Method:        http.MethodPost,
	OperationID:   "order",
	Path:          "/order",
	Summary:       "Order",
	Tags:          []string{"Order"},
	DefaultStatus: 200,
}
