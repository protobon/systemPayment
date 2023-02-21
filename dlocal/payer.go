package dlocal

// DlocalPayer example
type Payer struct {
	Name          *string `validate:"nonzero,min=6,max=100"`
	Email         *string `validate:"nonzero,min=6,max=100"`
	BirthDate     *string `validate:"nonzero"`
	Phone         *string `validate:"nonzero"`
	Document      *string `validate:"nonzero"`
	UserReference *string `validate:"nonzero"`
	Address       Address `validate:"nonzero"`
	IP            *string `json:"ip"`
	DeviceID      *string `json:"device_id"`
}

type Address struct {
	State   *string `validate:"nonzero"`
	City    *string `validate:"nonzero"`
	ZipCode *string `validate:"nonzero"`
	Street  *string `validate:"nonzero"`
	Number  *string `validate:"nonzero"`
}
