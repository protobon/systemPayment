package dlocal

import (
	"encoding/json"
	"net/http"
	"systempayment/model"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Payment
type PaymentRequestBody struct {
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	Country           string  `json:"country"`
	PaymentMethodID   string  `json:"payment_method_id"`
	PaymentMethodFlow string  `json:"payment_method_flow"`
	Payer             Payer   `json:"payer"`
	Card              Card    `json:"card"`
	OrderID           string  `json:"order_id"`
	Description       string  `json:"description"`
}

// Payment with token
type PaymentWithTokenRequestBody struct {
	Amount            float64       `json:"amount"`
	Currency          string        `json:"currency"`
	Country           string        `json:"country"`
	PaymentMethodID   string        `json:"payment_method_id"`
	PaymentMethodFlow string        `json:"payment_method_flow"`
	Payer             Payer         `json:"payer"`
	Card              CardWithToken `json:"card"`
	OrderID           string        `json:"order_id"`
	Description       string        `json:"description"`
}

// Payment Response
type PaymentResponseBody struct {
	ID                string       `json:"id"`
	Amount            float64      `json:"amount"`
	Currency          string       `json:"currency"`
	Country           string       `json:"country"`
	PaymentMethodID   string       `json:"payment_method_id"`
	PaymentMethodType string       `json:"payment_method_type"`
	PaymentMethodFlow string       `json:"payment_method_flow"`
	Card              CardResponse `json:"card"`
	CreatedDate       string       `json:"created_date"`
	ApprovedDate      string       `json:"approved_date"`
	Status            string       `json:"status"`
	StatusCode        string       `json:"status_code"`
	StatusDetail      string       `json:"status_detail"`
	OrderID           string       `json:"order_id"`
	NotificationUrl   string       `json:"notification_url"`
}

func MakePayment(order model.Order, payer model.Payer, card model.Card) (int, map[string]interface{}, error) {
	var req *http.Request
	var err error
	var dlocalCard = Card{CardId: card.CardId}

	// Payer's address
	DlocalAddress := Address{
		State:   *payer.Address.State,
		City:    *payer.Address.City,
		ZipCode: *payer.Address.ZipCode,
		Street:  *payer.Address.Street,
		Number:  *payer.Address.Number,
	}
	// Payer's info
	DlocalPayer := Payer{
		Name:          *payer.Name,
		Email:         *payer.Email,
		BirthDate:     *payer.BirthDate,
		Phone:         *payer.Phone,
		Document:      *payer.Document,
		UserReference: payer.UserReference,
		Address:       DlocalAddress,
	}
	// Payment request body
	Body := PaymentRequestBody{
		Amount:            order.Product.Amount / float64(order.TotalFees),
		Currency:          *order.Currency,
		Country:           *payer.Country,
		PaymentMethodID:   "CARD",
		PaymentMethodFlow: "DIRECT",
		Payer:             DlocalPayer,
		Card:              dlocalCard,
		OrderID:           order.OrderId,
	}

	body_json, err := json.Marshal(Body)
	if err != nil {
		log.Error("MakePayment - ", err)
		return 501, nil, err
	}
	// prepare dlocal POST request
	if req, err = DlocalPostRequest(body_json, "/payments"); err != nil {
		return 501, nil, err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	var res_body map[string]interface{}
	_ = json.NewDecoder(res.Body).Decode(&res_body)
	if err != nil {
		log.Error("MakePayment - ", err)
		if res != nil {
			return 501, res_body, err
		}
		return 408, nil, err
	}
	defer res.Body.Close()

	return 200, res_body, nil
}

func PaymentWithToken(payer model.Payer, token string) (int, map[string]interface{}, error) {
	var req *http.Request
	var err error
	var dlocalCard = CardWithToken{Token: token, Save: true}

	// Payer's address
	DlocalAddress := Address{
		State:   *payer.Address.State,
		City:    *payer.Address.City,
		ZipCode: *payer.Address.ZipCode,
		Street:  *payer.Address.Street,
		Number:  *payer.Address.Number,
	}
	// Payer's info
	DlocalPayer := Payer{
		Name:          *payer.Name,
		Email:         *payer.Email,
		BirthDate:     *payer.BirthDate,
		Phone:         *payer.Phone,
		Document:      *payer.Document,
		UserReference: payer.UserReference,
		Address:       DlocalAddress,
	}
	// Payment request body
	Body := PaymentWithTokenRequestBody{
		Amount:            1,
		Currency:          "USD",
		Country:           *payer.Country,
		PaymentMethodID:   "CARD",
		PaymentMethodFlow: "DIRECT",
		Payer:             DlocalPayer,
		Card:              dlocalCard,
		OrderID:           uuid.New().String(),
	}

	body_json, err := json.Marshal(Body)
	if err != nil {
		log.Error("PaymentWithToken - ", err)
		return 501, nil, err
	}
	// prepare dlocal POST request
	if req, err = DlocalPostRequest(body_json, "/payments"); err != nil {
		return 501, nil, err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	var res_body map[string]interface{}
	_ = json.NewDecoder(res.Body).Decode(&res_body)
	if err != nil {
		log.Error("PaymentWithToken - ", err)
		if res != nil {
			return 501, res_body, err
		}
		return 408, nil, err
	}
	defer res.Body.Close()

	return 200, res_body, nil
}
