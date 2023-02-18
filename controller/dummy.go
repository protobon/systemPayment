package controller

import (
	"net/http"
	"strconv"
	"systempayment/database"
	"systempayment/httputil"
	"systempayment/model"

	"github.com/gin-gonic/gin"
)

// NewDummy godoc
//
//	@Summary		Insert Dummy
//	@Description	save dummy in database
//	@Tags			dummy
//	@Accept			json
//	 @Param   example     body     model.Dummy     true  "Dummy example"     example(model.Dummy)
//	@Produce		json
//	@Success		200	{object}	model.Dummy
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/dummy/new [post]
func (c *Controller) NewDummy(ctx *gin.Context) {
	var dummy model.Dummy
	if err := ctx.BindJSON(&dummy); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	if code, err := dummy.QCreateDummy(database.DB); err != nil {
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

	ctx.JSON(200, dummy)
}

// Dummies godoc
//
//	@Summary		Select all Dummies
//	@Description	Select all Dummies
//	@Tags			dummy
//	@Accept			json
//
// @Param   start  query  int  true  "start example"  example(0)
// @Param   count  query  int  true  "count example"  example(10)
//
//	@Produce		json
//	@Success		200	{array}		model.Dummy
//	@Router			/dummy/dummies [get]
func (c *Controller) Dummies(ctx *gin.Context) {
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
	var dummy = model.Dummy{}
	dummies, code, err := dummy.QGetDummies(database.DB, start, count)
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

	ctx.JSON(200, dummies)
}

// GetDummy godoc
//
//	@Summary		Select Dummy
//	@Description	Get one Dummy from ID
//	@Tags			dummy
//	@Accept			json
//
// @Param   int  query  int  true  "example: 1"  "Dummy ID"
//
//	@Produce		json
//	@Success		200	{object}	model.Dummy
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/dummy/{id} [get]
func (c *Controller) GetDummy(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	dummy := model.Dummy{ID: id}
	if code, err := dummy.QGetDummy(database.DB); err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(200, dummy)
}

// UpdateDummy godoc
//
//	@Summary		Updates Dummy
//	@Description	Updates a dummy in database (id req)
//	@Tags			dummy
//	@Accept			json
//	 @Param   example     body     model.Dummy     true  "Dummy example"     example(model.Dummy)
//	@Produce		json
//	@Success		200	{object}	model.Dummy
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/dummy/update [put]
func (c *Controller) UpdateDummy(ctx *gin.Context) {
	var dummy model.Dummy
	if err := ctx.BindJSON(&dummy); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	if code, err := dummy.QUpdateDummy(database.DB); err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(200, dummy)
}
