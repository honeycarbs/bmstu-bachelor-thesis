package handler

import (
	"cluster/internal/config"
	"cluster/internal/service"
	"cluster/pkg/e"
	"cluster/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ClusterHandler struct {
	logger         logging.Logger
	sampleService  *service.SampleService
	clusterService *service.ClusterService
	frameService   *service.FrameService
}

func NewClusterHandler(logger logging.Logger, sampleService *service.SampleService, clusterService *service.ClusterService, frameService *service.FrameService) *ClusterHandler {
	return &ClusterHandler{logger: logger,
		sampleService:  sampleService,
		clusterService: clusterService,
		frameService:   frameService}
}

func (h *ClusterHandler) Register(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	h.logger.Trace("Register routes with prefix: /api/v1/")
	{
		v1.POST("clusters", h.createClusters)
	}
}

func (h *ClusterHandler) createClusters(ctx *gin.Context) {
	samples, err := h.sampleService.GetAll()
	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	sampleLength := len(samples)
	for i := 0; i < sampleLength; i++ {
		samples[i].Frames, err = h.frameService.GetAllBySample(samples[i].ID)
		if err != nil {
			h.logger.Info(err)
			e.NewErrorResponse(ctx, http.StatusBadRequest, err)
			return
		}
		h.logger.Infof("got %v frames from sample %v out of %v", len(samples[i].Frames), i, sampleLength)
	}

	frames, err := h.sampleService.CollectAllFrames(samples)
	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}
	fmt.Println(len(frames))

	clusterAmount := config.GetConfig().ClusterAmount
	maxRounds := config.GetConfig().KMeansMaxRounds

	h.logger.Infof("creating %v kmeans clusters with %v iterations", clusterAmount, maxRounds)
	clusters, err := h.clusterService.AssignClusters(frames, clusterAmount, maxRounds)
	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	h.logger.Info("Started creating cluster entities...")
	for i, cluster := range clusters {
		err := h.clusterService.CreateCluster(cluster)
		h.logger.Infof("Created cluster entity no. (%v)", i)
		if err != nil {
			h.logger.Info(err)
			e.NewErrorResponse(ctx, http.StatusBadRequest, err)
			return
		}
	}

	h.logger.Info("Started assigning clusters to frames...")
	for i, fm := range frames {
		err := h.frameService.AssignCluster(fm, clusters)
		if err != nil {
			h.logger.Info(err)
			e.NewErrorResponse(ctx, http.StatusBadRequest, err)
			return
		}
		h.logger.Infof("Assigned cluster to frame (%v) out of (%v)", i, len(frames))
	}

	h.logger.Info("All the clusters are assigned, finishing...")

	ctx.Writer.WriteHeader(http.StatusCreated)
}

//func (h *ClusterHandler) getClusters(ctx *gin.Context) {}
