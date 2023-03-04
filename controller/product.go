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
//	 @Param   example     body     model.ProductRequest     true  "Product example"     example(model.ProductRequest)
//		@in body
//		@Produce		json
//		@Success		200	{object}	model.ProductResponse
//		@Failure		400	{object}	httputil.HTTPError400
//		@Failure		404	{object}	httputil.HTTPError404
//		@Failure		500	{object}	httputil.HTTPError500
//		@Router			/product/new [post]
func (c *Controller) NewProduct(ctx *gin.Context) {
	var product model.Product
	if err := ctx.BindJSON(&product); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if code, err := product.QCreateProduct(database.DB); err != nil {
		switch code {
		case 400:
			httputil.NewError400(ctx, http.StatusBadRequest, "Body validation failed", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "Could not create Product", err)
		}
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
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid parameter: start", err)
		return
	}
	count, err := strconv.Atoi(ctx.Query("count"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid parameter: count", err)
		return
	}

	if count > 30 || count < 1 {
		count = 30
	}
	if start < 0 {
		start = 0
	}
	var product = model.Product{}
	products, code, err := product.QGetProducts(database.DB, start, count)
	if err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "Query returned 0 records", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "Error fetching Products", err)
		}
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
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/product/{id} [get]
func (c *Controller) GetProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "", err)
		return
	}

	product := model.Product{ID: id}
	if code, err := product.QGetProduct(database.DB); err != nil {
		switch code {
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "Product not found", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "Error fetching Product", err)
		}
		return
	}

	ctx.JSON(200, model.ProductResponse(product))
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
// @Param   example     body     model.ProductRequest     true  "Product example"     example(model.ProductRequest)
//
//	@Produce		json
//	@Success		200	{object}	model.ProductResponse
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/product/update/{id} [put]
func (c *Controller) UpdateProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid parameter: id", err)
		return
	}

	product := model.Product{ID: id}
	if err := ctx.BindJSON(&product); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if code, err := product.QUpdateProduct(database.DB); err != nil {
		switch code {
		case 400:
			httputil.NewError400(ctx, http.StatusBadRequest, "Body validation failed", err)
		case 404:
			httputil.NewError404(ctx, http.StatusNotFound, "Product not found", err)
		default:
			httputil.NewError500(ctx, http.StatusInternalServerError, "Error updating Product", err)
		}
		return
	}

	ctx.JSON(200, product)
}
