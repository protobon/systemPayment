package controller

import (
	"net/http"
	"strconv"
	"systempayment/database"
	"systempayment/httputil"
	"systempayment/model"

	"github.com/gin-gonic/gin"
)

// MockPayment godoc
//
//	@Summary		Mock Payment
//	@Description	Mocks a new Payment (for testing purposes)
//	@Tags			Payment
//	@Accept			json
//	 @Param   example     body     model.PaymentRequest     true  "Payment example"     example(model.PaymentRequest)
//	@Produce		json
//	@Success		200	{object}	model.PaymentResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/payment/new [post]
func (c *Controller) MockPayment(ctx *gin.Context) {
	payer_id, err := strconv.Atoi(ctx.Query("payer_id"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid parameter: payer_id", err)
		return
	}
	var p *model.Payer
	p, err = model.PreloadPayer(database.DB, payer_id)
	if err != nil {
		httputil.NewError404(ctx, http.StatusNotFound, "Payer not found", err)
		return
	}

	order_id, err := strconv.Atoi(ctx.Query("order_id"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid parameter: order_id", err)
		return
	}
	var o *model.Order
	o, err = model.PreloadOrder(database.DB, order_id, payer_id)
	if err != nil {
		httputil.NewError404(ctx, http.StatusNotFound, "Order not found", err)
		return
	}

	var payment model.Payment
	if err := ctx.BindJSON(&payment); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	payment.CardID = p.CardID
	payment.OrderID = o.ID
	if code, err := payment.QCreatePayment(database.DB); err != nil {
		switch code {
		case 400:
			httputil.NewError400(ctx, http.StatusBadRequest, "Body validation failed", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "Could not create Payment", err)
		}
		return
	}

	ctx.JSON(200, payment)
}
