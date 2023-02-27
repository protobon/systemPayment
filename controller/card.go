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
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	var payer = model.Payer{ID: payer_id}
	if code, err := payer.QGetPayer(database.DB); err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "", err)
		}
	}
	var card model.Card
	if err := ctx.BindJSON(&card); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}

	card.PayerID = payer.ID
	if code, err := card.QCreateCard(database.DB); err != nil {
		switch code {
		case 400:
			httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "", err)
		}
		return
	}
	payer.CardID = card.ID
	if code, err := payer.QUpdatePayer(database.DB); err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "", err)
		}
	}

	ctx.JSON(200, card)
}

// cards godoc
//
//	@Summary		Select all cards
//	@Description	Select all cards
//	@Tags			Card
//
// @Param   payer_id  query  int  true  "count example"  example(1)
// @Param   start  query  int  true  "start example"  example(0)
// @Param   count  query  int  true  "count example"  example(10)
//
//	@Produce		json
//	@Success		200	{array}		model.CardResponse
//	@Router			/card/cards [get]
func (c *Controller) Cards(ctx *gin.Context) {
	payer_id, err := strconv.Atoi(ctx.Query("payer_id"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	start, err := strconv.Atoi(ctx.Query("start"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}
	count, err := strconv.Atoi(ctx.Query("count"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}

	if count < 1 {
		count = 10
	}
	if count > 100 {
		count = 100
	}
	if start < 0 {
		start = 0
	}
	var card = model.Card{PayerID: payer_id}
	cards, code, err := card.QGetCards(database.DB, payer_id)
	if err != nil {
		switch code {
		case 400:
			httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "", err)
		}
		return
	}

	ctx.JSON(200, cards)
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
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}

	card := model.Card{ID: id}
	code, err := card.QGetCard(database.DB)
	if err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "", err)
		}
		return
	}

	ctx.JSON(200, card)
}
