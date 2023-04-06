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
//
//	@Param payer_id  query  int  true  "start example"  example(1)
//	@Param   example     body     model.OrderRequest     true  "Order example"     example(model.OrderRequest)
//
//	@Produce		json
//	@Success		200	{object}	model.OrderResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/order/new [post]
func (o *Controller) NewOrder(ctx *gin.Context) {
	payer_id, err := strconv.Atoi(ctx.Query("payer_id"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid parameter: payer_id", err)
		return
	}
	var order model.Order
	if err := ctx.BindJSON(&order); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid request payload", err)
		return
	}
	order.PayerID = payer_id

	if code, err := order.QCreateOrder(database.DB); err != nil {
		switch code {
		case 400:
			httputil.NewError400(ctx, http.StatusBadRequest, "Body validation failed", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "Could not create Order", err)
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
// @Param   payerId  query  int  false  "payerId example"  example(1)
//
//	@Produce		json
//	@Success		200	{array}		model.OrderResponse
//	@Router			/order/orders [get]
func (o *Controller) Orders(ctx *gin.Context) {
	start, err := strconv.Atoi(ctx.Query("start"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid parameter: start", err)
		return
	}
	count, err := strconv.Atoi(ctx.Query("count"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid parameter: count", err)
		return
	}
	payer_id := 0
	payer_id, _ = strconv.Atoi(ctx.Query("payerId"))

	if count > 30 || count < 1 {
		count = 30
	}
	if start < 0 {
		start = 0
	}
	var order = model.Order{}
	orders, code, err := order.QGetOrders(database.DB, start, count, payer_id)
	if err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "Query returned 0 records", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "Error fetching Orders", err)
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
//	@Success		200	{object}	model.OrderResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/order/{id} [get]
func (o *Controller) GetOrder(ctx *gin.Context) {
	// var out model.OrderResponse
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid parameter: id", err)
		return
	}

	order := model.Order{ID: id}
	code, err := order.QGetOrder(database.DB)
	if err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "Order not found", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "Error fetching order", err)
		}
		return
	}

	ctx.JSON(200, order)
}
