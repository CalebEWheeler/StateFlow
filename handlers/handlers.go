package handlers

type Handlers struct {
	CreateUserHandler    *CreateUserHandler
	CreateBillingHandler *CreateBillingHandler
	OrderHandler         *OrderHandler
	SendEmailHandler     *SendEmailHandler
}
