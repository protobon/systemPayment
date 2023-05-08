package dlocal

// Tokenized card for one use
type SecureCard struct {
	Token *string `json:"token" validate:"nonzero"`
	Save  bool    `json:"save"`
}

type Card struct {
	CardId *string `json:"card_id" validate:"nonzero"`
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
