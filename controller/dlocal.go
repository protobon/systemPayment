package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"systempayment/database"
	"systempayment/dlocal"
	"systempayment/httputil"
	"systempayment/model"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateCard godoc
//
//	@Summary		Create Card with Dlocal
//	@Description	Creates card and saves card token in database
//	@Tags			Dlocal
//	@Accept			json
//	 @Param   example     body     dlocal.Card     true  "Card example"     example(dlocal.Card)
//	@Produce		json
//	@Success		200	{object}	dlocal.CardResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/dlocal/card [post]
// func (c *Controller) CreateCard(ctx *gin.Context) {
// 	var req *http.Request
// 	var err error
// 	var Body dlocal.CardRequestBody
// 	if err := ctx.BindJSON(&Body); err != nil {
// 		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
// 		return
// 	}
// 	body_json, err := json.Marshal(Body)
// 	if err != nil {
// 		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
// 		return
// 	}
// 	if req, err = dlocal.PostRequest(body_json, "/secure_cards"); err != nil {
// 		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
// 		return
// 	}
// 	client := http.Client{
// 		Timeout: 30 * time.Second,
// 	}

// 	res, err := client.Do(req)
// 	if err != nil {
// 		ctx.JSON(400, err)
// 	}

// 	// response body
// 	res_body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		ctx.JSON(500, err)
// 	}
// 	ctx.JSON(200, res_body)
// }

// MakePayment godoc
//
//	@Summary		Make Payment with Dlocal
//	@Description	Makes a new Payment with the Dlocal API
//	@Tags			Dlocal
//	@Accept			json
//
// @Param   order_id  query  int  true  "order_id example"  example(32)
//
//	@Produce		json
//	@Success		200	{object}	dlocal.PaymentResponseBody
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/dlocal/payment [post]
func (c *Controller) MakePayment(ctx *gin.Context) {
	var order model.Order
	var req *http.Request
	var err error
	order.ID, err = strconv.Atoi(ctx.Query("order_id"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid parameter: order_id", err)
		return
	}
	code, err := order.QGetOrder(database.DB)
	if err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "Order not found", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "Error fetching Order", err)
		}
		ctx.JSON(200, order)
	}
	if order.Finished {
		httputil.NewError400(ctx, http.StatusBadRequest, "This order is already fully paid", nil)
	}

	payer := model.Payer{ID: order.PayerID}
	code, err = payer.QGetPayer(database.DB)
	if err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "Payer not found", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "Error fetching Payer", err)
		}
		return
	}
	card := model.Card{ID: payer.CardID}
	code, err = card.QGetCard(database.DB)
	if err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "Card not found", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "Error fetching Card", err)
		}
		return
	}

	// Prepare Request Body
	DlocalCard := dlocal.SecureCard{
		Token: card.Token,
	}
	DlocalAddress := dlocal.Address{
		State:   *payer.Address.State,
		City:    *payer.Address.City,
		ZipCode: *payer.Address.ZipCode,
		Street:  *payer.Address.Street,
		Number:  *payer.Address.Number,
	}
	DlocalPayer := dlocal.Payer{
		Name:      *payer.Name,
		Email:     *payer.Email,
		BirthDate: *payer.BirthDate,
		Phone:     *payer.Phone,
		Document:  *payer.Document,
		// UserReference: payer,
		Address: DlocalAddress,
	}
	Body := dlocal.PaymentRequestBody{
		Amount:            order.Product.Amount / float64(order.TotalFees),
		Currency:          *order.Currency,
		Country:           *payer.Country,
		PaymentMethodID:   "CARD",
		PaymentMethodFlow: "DIRECT",
		Payer:             DlocalPayer,
		Card:              DlocalCard,
	}

	body_json, err := json.Marshal(Body)
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "error generating request body", err)
		return
	}
	if req, err = dlocal.PostRequest(body_json, "/payments"); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "error preparing request to dlocal", err)
		return
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		if res != nil {
			ctx.JSON(res.StatusCode, err)
		}
		ctx.JSON(408, err)
	}
	defer res.Body.Close()

	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(500, err)
	}
	ctx.JSON(res.StatusCode, string(res_body))
}
