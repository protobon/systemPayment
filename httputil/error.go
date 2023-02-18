package httputil

import "github.com/gin-gonic/gin"

// NewError400 example
func NewError400(ctx *gin.Context, status int, err error) {
	er := HTTPError400{
		Code:  status,
		Error: err.Error(),
	}
	ctx.JSON(status, er)
}

// NewError404 example
func NewError404(ctx *gin.Context, status int, err error) {
	er := HTTPError404{
		Code:  status,
		Error: err.Error(),
	}
	ctx.JSON(status, er)
}

// NewError500 example
func NewError500(ctx *gin.Context, status int, err error) {
	er := HTTPError500{
		Code:  status,
		Error: err.Error(),
	}
	ctx.JSON(status, er)
}

// HTTPError400 example Bad Request
type HTTPError400 struct {
	Code  int    `json:"code" example:"400"`
	Error string `json:"error" example:"Invalid request body or query parameters"`
}

// HTTPError404 example Not Found
type HTTPError404 struct {
	Code  int    `json:"code" example:"404"`
	Error string `json:"error" example:"Page not found"`
}

// HTTPError500 example Server Error
type HTTPError500 struct {
	Code  int    `json:"code" example:"500"`
	Error string `json:"error" example:"Internal Server Error"`
}
