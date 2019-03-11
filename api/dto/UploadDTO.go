package dto

type CreateCardRequest struct {
	TokenId string `json:"token_id"`

	StripeCustomerId string `json:"stripe_customer_id"`

	ReturnURL string `json:"return_url"`

	Email string `json:"email"`

	IsThreeDSecure bool `form:"is_three_d_secure" json:"is_three_d_secure"`
}
