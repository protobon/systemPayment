package controller

import (
	"net/http"
	"strconv"
	"systempayment/database"
	"systempayment/httputil"
	"systempayment/model"

	"github.com/gin-gonic/gin"
)

// NewOrder godoc
//
//	@Summary		Insert Order
//	@Description	save Order in database
//	@Tags			Order
//	@Accept			json
//	 @Param   example     body     model.Order     true  "Order example"     example(model.Order)
//	@Produce		json
//	@Success		200	{object}	model.Order
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/order/new [post]
func (o *Controller) NewOrder(ctx *gin.Context) {
	var order model.Order
	if err := ctx.BindJSON(&order); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	if code, err := order.QCreateOrder(database.DB); err != nil {
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

	ctx.JSON(200, order)
}

// Orders godoc
//
//	@Summary		Select all Orders
//	@Description	Select all Orders
//	@Tags			Order
//	@Accept			json
//
// @Param   start  query  int  true  "start example"  example(0)
// @Param   count  query  int  true  "count example"  example(10)
//
//	@Produce		json
//	@Success		200	{array}		model.Order
//	@Router			/order/orders [get]
func (o *Controller) Orders(ctx *gin.Context) {
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
	var order = model.Order{}
	orders, code, err := order.QGetOrders(database.DB, start, count)
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

	ctx.JSON(200, orders)
}

// GetOrder godoc
//
//	@Summary		Select Order
//	@Description	Get one Order from ID
//	@Tags			Order
//	@Accept			json
//
// @Param   int  query  int  true  "example: 1"  "Order ID"
//
//	@Produce		json
//	@Success		200	{object}	model.Order
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/order/{id} [get]
func (o *Controller) GetOrder(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	order := model.Order{ID: id}
	if code, err := order.QGetOrder(database.DB); err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(200, order)
}

// UpdateOrder godoc
//
//	@Summary		Updates Order
//	@Description	Updates a Order in database (id req)
//	@Tags			Order
//	@Accept			json
//	 @Param   example     body     model.Order     true  "Order example"     example(model.Order)
//	@Produce		json
//	@Success		200	{object}	model.Order
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/order/update [put]
func (o *Controller) UpdateOrder(ctx *gin.Context) {
	var order model.Order
	if err := ctx.BindJSON(&order); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	if code, err := order.QUpdateOrder(database.DB); err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(200, order)
}
