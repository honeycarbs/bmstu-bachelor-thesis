package handler

import (
	"github.com/gin-gonic/gin"
	"ml/internal/service"
	"ml/pkg/e"
	"ml/pkg/logging"
	"net/http"
)

type MLHandler struct {
	logger        logging.Logger
	sampleService *service.SampleService
	mlService     *service.MLService
}

func NewMLHandler(logger logging.Logger, sampleService *service.SampleService, mlService *service.MLService) *MLHandler {
	return &MLHandler{logger: logger, sampleService: sampleService, mlService: mlService}
}

func (h *MLHandler) Register(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	h.logger.Trace("Register routes with prefix: /api/v1/")
	{
		v1.GET("/hmm/test", h.test)
		v1.POST("/hmm/train", h.train)
	}
}

func (h *MLHandler) train(ctx *gin.Context) {
	err := h.mlService.Train()
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Writer.WriteHeader(http.StatusNoContent)
}

func (h *MLHandler) test(ctx *gin.Context) {
	actual, predicted, err := h.mlService.Test()
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	h.mlService.GetHeatmap(actual, predicted)
	ctx.HTML(http.StatusOK, "heatmap.html", gin.H{
		"content": "This is an index page...",
	})
}
