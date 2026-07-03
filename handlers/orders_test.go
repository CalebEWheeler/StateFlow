package handlers_test

import (
	"context"
	"testing"

	"github.com/CalebEWheeler/StateFlow/handlers"
	"github.com/CalebEWheeler/StateFlow/shared"
	"github.com/CalebEWheeler/StateFlow/storage/postgres"
)

func TestHandleOrder(t *testing.T) {
	t.Parallel()

	input := handlers.OrderRequest{
		Body: shared.OrderRequestBody{
			Address: shared.Address{
				AdministrativeArea: "State",
				City:               "Test City",
				Country:            "Testopia",
				Street:             "123 Test Ave.",
			},
			Currency:   "USD",
			CustomerID: "customer123",
			Email:      "JohnDoe@test.com",
			Items: []shared.Item{
				{
					ID:       "1234567890",
					Quantity: 1,
					SKU:      "ABC123",
					MSRP:     29.99,
					Price:    24.99,
				},
			},
		},
	}

	ctx := context.Background()
	store, err := postgres.NewStore(ctx, "postgres://postgres:example@localhost:5432/stateflow")
	if err != nil {
		return
	}

	oh := handlers.NewOrderHandler(store)
	oh.Handle(ctx, &input)
}
