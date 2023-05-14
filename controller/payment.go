package controller

import (
	"net/http"
	"strconv"
	"systempayment/database"
	"systempayment/dlocal"
	"systempayment/httputil"
	"systempayment/model"

	"github.com/gin-gonic/gin"
)

// NewPayment godoc
//
//	@Summary		New Payment
//	@Description	Creates a new payment with dlocal
//	@Tags			Payment
//	@Accept			json
//
// @Param   order_id  query  int  true  "order_id example"  example(1)
//
//	@Produce		json
//	@Success		200	{object}	model.PaymentResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/payment/new [post]
func (c *Controller) NewPayment(ctx *gin.Context) {
	order_id, err := strconv.Atoi(ctx.Query("order_id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid parameter: order_id", err)
		return
	}
	var order = model.Order{ID: order_id}
	if code, err := order.QGetOrder(database.DB); err != nil {
		switch code {
		case 404:
			httputil.Error400(ctx, http.StatusBadRequest, "Order not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "An error occurred while fetching the order", err)
		}
		return
	}

	var payer = model.Payer{ID: order.PayerID}
	if code, err := payer.QGetPayer(database.DB); err != nil {
		switch code {
		case 404:
			httputil.Error400(ctx, http.StatusBadRequest, "Payer not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "An error occurred while fetching the payer", err)
		}
		return
	}

	var card = model.Card{ID: payer.CardID}
	if code, err := card.QGetCard(database.DB); err != nil {
		switch code {
		case 404:
			httputil.Error400(ctx, http.StatusBadRequest, "Card not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "An error occurred while fetching the card", err)
		}
		return
	}

	code, response, err := dlocal.MakePayment(order, payer, card)
	if err != nil {
		switch code {
		case 408:
			httputil.Error408(ctx, http.StatusBadRequest, "Request to dlocal timed out", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Request to dlocal failed", err)
		}
		return
	}

	var payment = model.Payment{
		OrderID: order.ID,
		CardID:  card.ID,
	}
	code, err = payment.SavePaymentFromResponse(database.DB, response)
	if err != nil {
		switch code {
		case 400:
			httputil.Error400(ctx, http.StatusBadRequest, "Payment validation failed", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Could not save new Payment", err)
		}
		return
	}

	ctx.JSON(200, payment)
}

type Token struct {
	Token string `json:"token"`
}

// PaymentWithToken godoc
//
//	@Summary		New Payment with Card's token
//	@Description	Creates a new payment with a CC token
//	@Tags			Payment
//	@Accept			json
//
// @Param   order_id  query  int  true  "order_id example"  example(1)
//
//	@Produce		json
//	@Success		200	{object}	model.PaymentResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/payment/save-card [post]
func (c *Controller) PaymentWithToken(ctx *gin.Context) {
	var token Token
	if err := ctx.BindJSON(&token); err != nil || token.Token == "" {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid request payload", err)
		return
	}
	order_id, err := strconv.Atoi(ctx.Query("order_id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid parameter: order_id", err)
		return
	}
	var order = model.Order{ID: order_id}
	if code, err := order.QGetOrder(database.DB); err != nil {
		switch code {
		case 404:
			httputil.Error400(ctx, http.StatusBadRequest, "Order not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "An error occurred while fetching the order", err)
		}
		return
	}

	var payer = model.Payer{ID: order.PayerID}
	if code, err := payer.QGetPayer(database.DB); err != nil {
		switch code {
		case 404:
			httputil.Error400(ctx, http.StatusBadRequest, "Payer not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "An error occurred while fetching the payer", err)
		}
		return
	}

	code, result, err := dlocal.PaymentWithToken(order, payer, token.Token)
	if err != nil {
		switch code {
		case 408:
			httputil.Error408(ctx, http.StatusBadRequest, "Request to dlocal timed out", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Request to dlocal failed", err)
		}
		return
	}
	var card = model.Card{PayerID: payer.ID}
	code, err = card.SaveCardFromResponse(database.DB, result)
	if err != nil {
		switch code {
		case 400:
			httputil.Error400(ctx, http.StatusBadRequest, "Card validation failed", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Could not save new Card", err)
		}
		return
	}

	var payment = model.Payment{
		OrderID: order.ID,
		CardID:  card.ID,
	}
	code, err = payment.SavePaymentFromResponse(database.DB, result)
	if err != nil {
		switch code {
		case 400:
			httputil.Error400(ctx, http.StatusBadRequest, "Payment validation failed", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Could not save new Payment", err)
		}
		return
	}

	ctx.JSON(200, payment)
}

// GetPayments godoc
//
//	@Summary		Select all Payments
//	@Description	Select all Payments
//	@Tags			Payment
//
// @Param   start  query  int  true  "start example"  example(0)
// @Param   count  query  int  true  "count example"  example(10)
// @Param   orderId  query  int  false  "orderId example"  example(1)
//
//	@Produce		json
//	@Success		200	{array}		model.PaymentResponse
//	@Router			/payment/payments [get]
func (c *Controller) GetPayments(ctx *gin.Context) {
	start, err := strconv.Atoi(ctx.Query("start"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid parameter: start", err)
		return
	}
	count, err := strconv.Atoi(ctx.Query("count"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid parameter: count", err)
		return
	}
	order_id := 0
	order_id, _ = strconv.Atoi(ctx.Query("orderId"))

	if count > 30 || count < 1 {
		count = 30
	}
	if start < 0 {
		start = 0
	}
	var payment = model.Payment{}
	payments, code, err := payment.QGetAllPayments(database.DB, start, count, order_id)
	if err != nil {
		switch code {
		case 404:
			httputil.Error404(ctx, http.StatusNotFound, "Query returned 0 records", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Error fetching Payments", err)
		}
		return
	}

	ctx.JSON(200, payments)
}
