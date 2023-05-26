package controller

import (
	"net/http"
	"strconv"
	"systempayment/database"
	"systempayment/httputil"
	"systempayment/model"

	"github.com/gin-gonic/gin"
)

// NewProduct godoc
//
//		@Summary		Insert Product
//		@Description	save Product in database
//		@Tags			Product
//		@Accept			json
//	 @Param   product     body     model.ProductRequest     true  "Product example"     example(model.ProductRequest)
//		@in body
//		@Produce		json
//		@Success		200	{object}	model.ProductResponse
//		@Failure		400	{object}	httputil.HTTPError400
//		@Failure		500	{object}	httputil.HTTPError500
//		@Router			/product/new [post]
func (c *Controller) NewProduct(ctx *gin.Context) {
	var product model.Product
	if err := ctx.BindJSON(&product); err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if _, err := product.QCreateProduct(database.DB); err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Body validation failed", err)
		return
	}

	ctx.JSON(200, product)
}

// Products godoc
//
//	@Summary		Select all Products
//	@Description	Select all Products
//	@Tags			Product
//
// @Param   start  query  int  true  "start example"  example(0)
// @Param   count  query  int  true  "count example"  example(10)
//
//	@Produce		json
//	@Success		200	{array}		model.ProductResponse
//	@Router			/product/products [get]
func (c *Controller) Products(ctx *gin.Context) {
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
	var product = model.Product{}
	products, _, err := product.QGetProducts(database.DB, start, count)
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Query returned 0 records", err)
		return
	}

	ctx.JSON(200, products)
}

// GetProduct godoc
//
//	@Summary		Select Product
//	@Description	Get one Product from ID
//	@Tags			Product
//
// @Param   int  query  int  true  "example: 1"  "Product ID"
//
//	@Produce		json
//	@Success		200	{object}	model.ProductResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/product/{id} [get]
func (c *Controller) GetProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "", err)
		return
	}

	product := model.Product{ID: id}
	if _, err := product.QGetProduct(database.DB); err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Product not found", err)
		return
	}

	ctx.JSON(200, product)
}

// UpdateProduct godoc
//
//	@Summary		Updates Product
//	@Description	Updates a Product in database (id req)
//	@Tags			Product
//	@Accept			json
//
// @Param   int  query  int  true  "example: 1"  "Payer ID"
//
// @Param   product     body     model.ProductRequest     true  "Product example"     example(model.ProductRequest)
//
//	@Produce		json
//	@Success		200	{object}	model.ProductResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/product/update/{id} [put]
func (c *Controller) UpdateProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid parameter: id", err)
		return
	}

	product := model.Product{ID: id}
	if err := ctx.BindJSON(&product); err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if _, err := product.QUpdateProduct(database.DB); err != nil {
		httputil.Error400(ctx, http.StatusBadRequest, "Invalid request payload or query params", err)
		return
	}

	ctx.JSON(200, product)
}
