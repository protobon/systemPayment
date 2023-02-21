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
//	 @Param   example     body     model.Payer     true  "Payer example"     example(model.Payer)
//		@Produce		json
//		@Success		200	{object}	model.PayerResponse
//		@Failure		400	{object}	httputil.HTTPError400
//		@Failure		404	{object}	httputil.HTTPError404
//		@Failure		500	{object}	httputil.HTTPError500
//		@Router			/payer/new [post]
func (c *Controller) NewPayer(ctx *gin.Context) {
	var payer model.Payer
	if err := ctx.BindJSON(&payer); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	if code, err := payer.QCreatePayer(database.DB); err != nil {
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
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}
	count, err := strconv.Atoi(ctx.Query("count"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	var payer = model.Payer{}
	payers, code, err := payer.QGetPayers(database.DB, start, count)
	if err != nil {
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
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	payer := model.Payer{ID: id}
	if code, err := payer.QGetPayer(database.DB); err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, err)
		case 500:
			httputil.NewError500(ctx, http.StatusInternalServerError, err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, err)
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
//	 @Param   example     body     model.Payer     true  "Payer example"     example(model.Payer)
//	@Produce		json
//	@Success		200	{object}	model.PayerResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/payer/update [put]
func (c *Controller) UpdatePayer(ctx *gin.Context) {
	var payer model.Payer
	if err := ctx.BindJSON(&payer); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	if code, err := payer.QUpdatePayer(database.DB); err != nil {
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

	ctx.JSON(200, payer)
}
