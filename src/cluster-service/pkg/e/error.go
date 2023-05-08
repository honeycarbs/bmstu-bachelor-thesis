package e

import (
	"github.com/gin-gonic/gin"
)

func NewErrorResponse(ctx *gin.Context, status int, err error) {
	er := ErrorResponse{
		Code:    status,
		Message: err.Error(),
	}

	ctx.AbortWithStatusJSON(status, er)
}

type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
