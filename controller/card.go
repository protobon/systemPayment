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
//	@Description	Creates a new payment with a CC token, saves card returned by dlocal.
//	@Tags			Payment
//	@Accept			json
//
// @Param   payer_id  query  int  true  "payer_id example"  example(1)
// @Param   token     body     model.Token    true  "Card's token example"     example(model.Token)
//
//	@Produce		json
//	@Success		200	{object}	model.PaymentResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
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
	if code, err := payer.QGetPayer(database.DB); err != nil {
		switch code {
		case 404:
			httputil.Error400(ctx, http.StatusBadRequest, "Payer not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "An error occurred while fetching the payer", err)
		}
		return
	}

	code, result, err := dlocal.PaymentWithToken(payer, token.Token)
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

	ctx.JSON(200, result)
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
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/card/{id} [get]
func (o *Controller) GetCard(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "", err)
		return
	}

	card := model.Card{ID: id}
	code, err := card.QGetCard(database.DB)
	if err != nil {
		switch code {
		case 404:
			httputil.Error404(ctx, http.StatusNotFound, "", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "", err)
		}
		return
	}

	ctx.JSON(200, card)
}
