package controller

import (
	"net/http"
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
//	 @Param   example     body     model.Payment     true  "Payment example"     example(model.Payment)
//	@Produce		json
//	@Success		200	{object}	model.PaymentResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/payment/new [post]
func (c *Controller) MockPayment(ctx *gin.Context) {
	var payment model.Payment
	if err := ctx.BindJSON(&payment); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	if code, err := payment.QCreatePayment(database.DB); err != nil {
		switch code {
		case 400:
			httputil.NewError400(ctx, http.StatusBadRequest, err)
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(200, payment)
}
