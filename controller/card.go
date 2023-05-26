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

// SaveCard godoc
//
//	@Summary		Saves a new Card
//	@Description	Creates a new payment of 1USD with a CC token, saves card returned by dlocal.
//	@Tags			Card
//	@Accept			json
//
// @Param   payer_id  query  int  true  "payer_id example"  example(1)
// @Param   token     body     model.Token    true  "Card's token example"     example(model.Token)
//
//	@Produce		json
//	@Success		200	{object}	model.PaymentResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/card/save-card [post]
func (c *Controller) SaveCard(ctx *gin.Context) {
	var token model.Token
	if err := ctx.BindJSON(&token); err != nil || token.Token == "" {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid request payload", err)
		return
	}
	payer_id, err := strconv.Atoi(ctx.Query("payer_id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid parameter: payer_id", err)
		return
	}
	var payer = model.Payer{ID: payer_id}
	if code, err := payer.QGetPayer(database.DB); code != 200 {
		httputil.Error400(ctx, http.StatusBadRequest, "Payer not found", err)
		return
	}

	code, response, err := dlocal.PaymentWithToken(payer, token.Token)
	if err != nil {
		ctx.JSON(code, err)
	}

	if code != 200 {
		ctx.JSON(code, response)
	}

	var card = model.Card{PayerID: payer.ID}
	code, err = card.SaveCardFromResponse(database.DB, response)
	if code != 200 {
		httputil.Error400(ctx, http.StatusBadRequest, "Card validation failed", err)
		return
	}

	ctx.JSON(200, response)
}

// GetCard godoc
//
//	@Summary		Select Card
//	@Description	Get one Card from ID
//	@Tags			Card
//	@Accept			json
//
// @Param   int  query  int  true  "example: 1"  "Card ID"
//
//	@Produce		json
//	@Success		200	{object}	model.CardResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/card/{id} [get]
func (o *Controller) GetCard(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid card ID", err)
		return
	}

	card := model.Card{ID: id}
	code, _ := card.QGetCard(database.DB)
	if code != 200 {
		httputil.Error400(ctx, http.StatusBadRequest, "Card not found", err)
		return
	}
	ctx.JSON(200, card)
}
