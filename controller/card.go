package controller

import (
	"net/http"
	"strconv"
	"systempayment/database"
	"systempayment/httputil"
	"systempayment/model"

	"github.com/gin-gonic/gin"
)

// Card godoc
//
//	@Summary		Insert Card
//	@Description	Inserts a new Card
//	@Tags			Card
//	@Accept			json
//
//	@Param   payer_id  query  int  true  "count example"  example(1)
//	@Param   example     body     model.CardRequest     true  "Card example"     example(model.CardRequest)
//
//	@Produce		json
//	@Success		200	{object}	model.CardResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/card/new [post]
func (c *Controller) NewCard(ctx *gin.Context) {
	payer_id, err := strconv.Atoi(ctx.Query("payer_id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid parameter: payer_id", err)
		return
	}
	var payer = model.Payer{ID: payer_id}
	if code, err := payer.QGetPayer(database.DB); err != nil {
		switch code {
		case 404:
			httputil.Error404(ctx, http.StatusNotFound, "Payer not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Error fetching Payer", err)
		}
	}
	var card model.Card
	if err := ctx.BindJSON(&card); err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	card.PayerID = payer.ID
	if code, err := card.QCreateCard(database.DB); err != nil {
		switch code {
		case 400:
			httputil.Error400(ctx, http.StatusBadRequest, "Body validation failed", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Could not create Card", err)
		}
		return
	}

	if code, err := payer.QPrimaryCard(database.DB, card.ID); err != nil {
		switch code {
		case 404:
			httputil.Error404(ctx, http.StatusNotFound, "Payer not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Could not set new Card as primary", err)
		}
	}

	ctx.JSON(200, card)
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
