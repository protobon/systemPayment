package controller

import (
	"net/http"
	"strconv"
	"systempayment/database"
	"systempayment/httputil"
	"systempayment/model"

	"github.com/gin-gonic/gin"
)

// NewPayer godoc
//
//		@Summary		Insert Payer
//		@Description	save payer in database
//		@Tags			Payer
//		@Accept			json
//	 @Param   payer     body     model.Payer     true  "Payer example"     example(model.Payer)
//		@Produce		json
//		@Success		200	{object}	model.PayerResponse
//		@Failure		400	{object}	httputil.HTTPError400
//		@Failure		404	{object}	httputil.HTTPError404
//		@Failure		500	{object}	httputil.HTTPError500
//		@Router			/payer/new [post]
func (c *Controller) NewPayer(ctx *gin.Context) {
	var payer model.Payer
	if err := ctx.BindJSON(&payer); err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if code, err := payer.QCreatePayer(database.DB); err != nil {
		switch code {
		case 400:
			httputil.Error400(ctx, http.StatusBadRequest, "Body validation failed", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Error creating Payer", err)
		}
		return
	}

	ctx.JSON(200, payer)
}

// Payers godoc
//
//	@Summary		Select all Payers
//	@Description	Select all Payers
//	@Tags			Payer
//
// @Param   start  query  int  true  "start example"  example(0)
// @Param   count  query  int  true  "count example"  example(10)
//
//	@Produce		json
//	@Success		200	{array}		model.PayerResponse
//	@Router			/payer/payers [get]
func (c *Controller) Payers(ctx *gin.Context) {
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

	if count > 30 || count < 1 {
		count = 30
	}
	if start < 0 {
		start = 0
	}
	var payer = model.Payer{}
	payers, code, err := payer.QGetPayers(database.DB, start, count)
	if err != nil {
		switch code {
		case 404:
			httputil.Error404(ctx, http.StatusNotFound, "Query returned 0 records", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Error fetching Payers", err)
		}
		return
	}

	ctx.JSON(200, payers)
}

// GetPayer godoc
//
//	@Summary		Select Payer
//	@Description	Get one Payer from ID
//	@Tags			Payer
//
// @Param   int  query  int  true  "example: 1"  "Payer ID"
//
//	@Produce		json
//	@Success		200	{object}	model.PayerResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/payer/{id} [get]
func (c *Controller) GetPayer(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "", err)
		return
	}

	payer := model.Payer{ID: id}
	// var payer_out model.PayerResponse
	var code int
	code, err = payer.QGetPayer(database.DB)
	if err != nil {
		switch code {
		case 404:
			httputil.Error404(ctx, http.StatusNotFound, "Payer not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Error fetching Payer", err)
		}
		return
	}

	ctx.JSON(200, payer)
}

// UpdatePayer godoc
//
//	@Summary		Updates Payer
//	@Description	Updates a payer in database (id req)
//	@Tags			Payer
//	@Accept			json
//
// @Param   id  query  int  true  "example: 1"  "Payer ID"
// @Param   payer     body     model.Payer     true  "Payer example"     example(model.Payer)
//
//	@Produce		json
//	@Success		200	{object}	model.PayerResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/payer/update/{id} [put]
func (c *Controller) UpdatePayer(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid parameter: id", err)
		return
	}

	payer := model.Payer{ID: id}
	if err := ctx.BindJSON(&payer); err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if code, err := payer.QUpdatePayer(database.DB); err != nil {
		switch code {
		case 400:
			httputil.Error400(ctx, http.StatusBadRequest, "Body validation failed", err)
		case 404:
			httputil.Error404(ctx, http.StatusNotFound, "Payer not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Error updating Payer", err)
		}
		return
	}

	ctx.JSON(200, payer)
}

// PrimaryCard godoc
//
//	@Summary		Sets Payers primary card
//	@Description	?payer_id=1&card_id=1
//	@Tags			Payer
//	@Accept			json
//
// @Param   payer_id  query  int  true  "example: 1"  "Payer ID"
// @Param   card_id  query  int  true  "example: 1"  "Card ID"
//
//	@Produce		json
//	@Success		200	{object}	model.PayerResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/payer/primary-card [put]
func (c *Controller) PrimaryCard(ctx *gin.Context) {
	payer_id, err := strconv.Atoi(ctx.Query("payer_id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid parameter: payer_id", err)
		return
	}
	card_id, err := strconv.Atoi(ctx.Query("card_id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid parameter: card_id", err)
		return
	}

	payer := model.Payer{ID: payer_id}
	if code, err := payer.QGetPayer(database.DB); err != nil {
		switch code {
		case 404:
			httputil.Error404(ctx, http.StatusNotFound, "Payer not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Error updating Payer", err)
		}
	}

	if code, err := payer.QPrimaryCard(database.DB, card_id); err != nil {
		switch code {
		case 400:
			httputil.Error400(ctx, http.StatusNotFound, "Invalid card ID", err)
		case 404:
			httputil.Error404(ctx, http.StatusNotFound, "Card not found", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Error updating Payer", err)
		}
		return
	}

	ctx.JSON(200, payer)
}

// PayerCards godoc
//
//	@Summary		Get all cards from payer id
//	@Description	?payer_id=1
//	@Tags			Payer
//
// @Param   payer_id  query  int  true  "payer_id example"  example(1)
//
//	@Produce		json
//	@Success		200	{array}		model.CardResponse
//	@Router			/payer/cards [get]
func (c *Controller) PayerCards(ctx *gin.Context) {
	payer_id, err := strconv.Atoi(ctx.Query("payer_id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid parameter: start", err)
		return
	}

	card := model.Card{PayerID: payer_id}
	cards, code, err := card.QGetCards(database.DB, payer_id)
	if err != nil {
		switch code {
		case 404:
			httputil.Error404(ctx, http.StatusNotFound, "Query returned 0 records", err)
		default:
			httputil.Error500(ctx, http.StatusInternalServerError, "Error fetching cards", err)
		}
		return
	}

	ctx.JSON(200, cards)
}
