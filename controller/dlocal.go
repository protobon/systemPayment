package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"systempayment/dlocal"
	"systempayment/httputil"
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
func (c *Controller) CreateCard(ctx *gin.Context) {
	var req *http.Request
	var err error
	var Body dlocal.CardRequestBody
	if err := ctx.BindJSON(&Body); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	body_json, err := json.Marshal(Body)
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	if req, err = dlocal.PostRequest(body_json, "/secure_cards"); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		ctx.JSON(400, err)
	}

	// response body
	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(500, err)
	}
	ctx.JSON(200, res_body)
}

// MakePayment godoc
//
//	@Summary		Make Payment with Dlocal
//	@Description	Makes a new Payment with the Dlocal API
//	@Tags			Dlocal
//	@Accept			json
//	 @Param   example     body     dlocal.PaymentRequestBody     true  "Payment with Dlocal example"     example(dlocal.PaymentRequestBody)
//	@Produce		json
//	@Success		200	{object}	dlocal.PaymentResponseBody
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/dlocal/payment [post]
func (c *Controller) MakePayment(ctx *gin.Context) {
	var req *http.Request
	var err error
	var Body dlocal.PaymentRequestBody
	if err := ctx.BindJSON(&Body); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	body_json, err := json.Marshal(Body)
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	if req, err = dlocal.PostRequest(body_json, "/payments"); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		ctx.JSON(400, err)
	}

	// response body
	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(500, err)
	}
	ctx.JSON(200, res_body)
}

// MakeSecurePayment godoc
//
//	@Summary		Make Secure Payment with Dlocal
//	@Description	Makes a new Secure Payment with the Dlocal API
//	@Tags			Dlocal
//	@Accept			json
//	 @Param   example     body     dlocal.SecurePaymentRequestBody     true  "Secure Payment with Dlocal example"     example(dlocal.SecurePaymentRequestBody)
//	@Produce		json
//	@Success		200	{object}	dlocal.SecurePaymentRequestBody
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/dlocal/secure-payment [post]
func (c *Controller) MakeSecurePayment(ctx *gin.Context) {
	var req *http.Request
	var err error
	var Body dlocal.SecurePaymentRequestBody
	if err := ctx.BindJSON(&Body); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	body_json, err := json.Marshal(Body)
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	if req, err = dlocal.PostRequest(body_json, "/secure_payments"); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		ctx.JSON(400, err)
	}

	// response body
	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(500, err)
	}
	ctx.JSON(200, res_body)
}
