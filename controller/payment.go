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
// @Param	auto  query	 bool  false  "auto example"  example(true)
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
	if code, err := order.GetOrderForPayment(database.DB); err != nil {
		switch code {
		case 404:
			httputil.Error404(ctx, http.StatusNotFound, "Order not found or already finished", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "An error occurred while fetching the order", err)
		}
		return
	}
	auto, _ := strconv.ParseBool(ctx.Query("auto"))
	if !order.Auto && auto {
		order.Auto = auto
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
		ctx.JSON(code, err)
	}
	if code != 200 {
		ctx.JSON(code, response)
	}

	code, err = order.PaymentSuccessful(database.DB)
	if err != nil {
		switch code {
		case 400:
			httputil.Error400(ctx, http.StatusBadRequest, "Order validation failed", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Could not update Order", err)
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
