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
//	 @Param   example     body     model.Product     true  "Product example"     example(model.Product)
//		@in body
//		@Produce		json
//		@Success		200	{object}	model.Product
//		@Failure		400	{object}	httputil.HTTPError400
//		@Failure		404	{object}	httputil.HTTPError404
//		@Failure		500	{object}	httputil.HTTPError500
//		@Router			/product/new [post]
func (c *Controller) NewProduct(ctx *gin.Context) {
	var product model.Product
	if err := ctx.BindJSON(&product); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	if code, err := product.QCreateProduct(database.DB); err != nil {
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
//	@Success		200	{array}		model.Product
//	@Router			/product/products [get]
func (c *Controller) Products(ctx *gin.Context) {
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
	var product = model.Product{}
	products, code, err := product.QGetProducts(database.DB, start, count)
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
//	@Success		200	{object}	model.Product
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/product/{id} [get]
func (c *Controller) GetProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	product := model.Product{ID: id}
	if code, err := product.QGetProduct(database.DB); err != nil {
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

	ctx.JSON(200, product)
}

// UpdateProduct godoc
//
//	@Summary		Updates Product
//	@Description	Updates a Product in database (id req)
//	@Tags			Product
//	@Accept			json
//	 @Param   example     body     model.Product     true  "Product example"     example(model.Product)
//	@Produce		json
//	@Success		200	{object}	model.Product
//	@Failure		400	{object}	httputil.HTTPError400
//	@Failure		404	{object}	httputil.HTTPError404
//	@Failure		500	{object}	httputil.HTTPError500
//	@Router			/product/update [put]
func (c *Controller) UpdateProduct(ctx *gin.Context) {
	var product model.Product
	if err := ctx.BindJSON(&product); err != nil {
		httputil.NewError400(ctx, http.StatusBadRequest, err)
		return
	}

	if code, err := product.QUpdateProduct(database.DB); err != nil {
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

	ctx.JSON(200, product)
}
