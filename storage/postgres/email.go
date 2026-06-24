package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/CalebEWheeler/StateFlow/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Email struct {
	To      string
	Subject string
	Body    string
}

type OrderConfirmationData struct {
	Address        shared.Address
	Carrier        string
	Items          []shared.Item
	OrderID        uuid.UUID
	RecipientEmail string
	TrackingNumber string
}

type EmailStore struct {
	pool *pgxpool.Pool
}

func NewEmailStore(pool *pgxpool.Pool) *EmailStore {
	return &EmailStore{pool: pool}
}

func (es EmailStore) SendConfirmation(ctx context.Context, job *Job) error {

	var data OrderConfirmationData

	err := es.pool.QueryRow(ctx, `
		SELECT 
			o.address,
			o.email,
			o.items,
			o.id,
			s.carrier,
			s.tracking_number
		FROM orders o
		JOIN shipments s
		ON o.id = $1
		WHERE o.id = $1
	`, job.OrderID).Scan(
		&data.Address,
		&data.RecipientEmail,
		&data.Items,
		&data.OrderID,
		&data.Carrier,
		&data.TrackingNumber,
	)

	if err != nil {
		return err
	}

	email := Email{
		To:      data.RecipientEmail,
		Subject: fmt.Sprintf("Order %s has shipped", data.OrderID),
		Body:    BuildOrderConfirmationBody(data),
	}

	fmt.Printf(
		"Sending email to %s\nSubject: %s\n%s",
		email.To,
		email.Subject,
		email.Body,
	)

	return nil
}

func BuildOrderConfirmationBody(data OrderConfirmationData) string {
	var items strings.Builder
	for _, item := range data.Items {
		fmt.Fprintf(
			&items,
			"- %s (Qty: %d)\n",
			item.SKU,
			item.Quantity,
		)
	}
	return fmt.Sprintf(`
Hello,

Your order has been processed and a label has been created. 

Order ID: %s

Items: 
%s

Carrier: %s
Tracking Number: %s

Shipping Address
%s %s %s %s

Thank you for your purchase.
	`,
		data.OrderID.String(),
		items.String(),
		data.Carrier,
		data.TrackingNumber,
		data.Address.Street,
		data.Address.City,
		data.Address.AdministrativeArea,
		data.Address.Country,
	)
}
