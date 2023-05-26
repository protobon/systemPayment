package httputil

import "github.com/gin-gonic/gin"

// Error400 example
func Error400(ctx *gin.Context, status int, message string, err error) {
	er := HTTPError400{
		Code:    status,
		Message: message,
		Error:   err.Error(),
	}
	ctx.JSON(status, er)
}

// Error408 example
func Error408(ctx *gin.Context, status int, message string, err error) {
	er := HTTPError400{
		Code:    status,
		Message: message,
		Error:   err.Error(),
	}
	ctx.JSON(status, er)
}

// Error500 example
func Error500(ctx *gin.Context, status int, message string, err error) {
	er := HTTPError500{
		Code:    status,
		Message: message,
		Error:   err.Error(),
	}
	ctx.JSON(status, er)
}

// HTTPError400 example Bad Request
type HTTPError400 struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message"`
	Error   string `json:"error" example:"Invalid request payload or query params"`
}

// HTTPError500 example Server Error
type HTTPError500 struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message"`
	Error   string `json:"error" example:"Internal Server Error"`
}
