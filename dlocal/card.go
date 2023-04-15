package dlocal

type SecureCard struct {
	Token *string `json:"token" validate:"nonzero"`
}

// Create Card with Dlocal

// Card object
type Card struct {
	HolderName      string `json:"holder_name" example:"Jhon Doe" validate:"nonzero"`
	Number          string `json:"number" example:"4111111111111111" validate:"nonzero"`
	CVV             string `json:"cvv" example:"123" validate:"nonzero,min=3,max=3"`
	ExpirationMonth int    `json:"expiration_month" example:"3" validate:"nonzero,min=1,max=12"`
	ExpirationYear  int    `json:"expiration_year" example:"2033" validate:"nonzero"`
}

// Request
type CardRequestBody struct {
	Country *string `json:"country" validate:"nonzero,min=2,max=2"`
	Card    Card    `json:"card" validate:"nonzero"`
	Payer   Payer   `json:"payer" validate:"nonzero"`
}

// Response
type CardResponse struct {
	CardID          string `json:"card_id"`
	HolderName      string `json:"holder_name"`
	ExpirationMonth string `json:"expiration_month"`
	ExpirationYear  string `json:"expiration_year"`
	Last4           string `json:"last4"`
	Brand           string `json:"brand"`
}
