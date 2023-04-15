package dlocal

// DlocalPayer example
type Payer struct {
	Name      string `json:"name" validate:"nonzero,min=6,max=100"`
	Email     string `json:"email" validate:"nonzero,min=6,max=100"`
	BirthDate string `json:"birth_date" validate:"nonzero"`
	Phone     string `json:"phone" validate:"nonzero"`
	Document  string `json:"document" validate:"nonzero"`
	// UserReference *string `json:"user_reference" validate:"nonzero"`
	Address Address `json:"address" validate:"nonzero"`
	// IP            *string `json:"ip"`
	// DeviceID      *string `json:"device_id"`
}

type Address struct {
	State   string `json:"state" validate:"nonzero"`
	City    string `json:"city" validate:"nonzero"`
	ZipCode string `json:"zip_code" validate:"nonzero"`
	Street  string `json:"street" validate:"nonzero"`
	Number  string `json:"number" validate:"nonzero"`
}
