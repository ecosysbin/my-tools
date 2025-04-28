package handler

import (
	"vcluster-gateway/pkg/apis/response"
	v1 "vcluster-gateway/pkg/apis/v1"
	"vcluster-gateway/pkg/controller/framework"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var _ framework.VClusterHandlerInterface = &VClusterHandler{}

type VClusterHandler struct {
	controller framework.Interface
	// repo       pkg.VirtualServerRepo
}

func NewVClusterHandler(controller framework.Interface) *VClusterHandler {
	return &VClusterHandler{
		controller: controller,
		// repo:       repoImpl,
	}
}

const (
	EMPTY_STORAGE            int64 = 0
	DEFAULT_MIN_STORAGE      int64 = 20
	DEFAULT_MAX_STORAGE      int64 = 100
	DEFAULT_DEFAULT_STORGAGE int64 = 50
)

// ListProductProfiles godoc
//
//	@Summary		List Product Profiles
//	@Description	List Product Profiles
//	@Tags			Profiles
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse{data=v1.Profiles}
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/product-profiles [get]
func (w *VClusterHandler) ListProductProfiles(c *gin.Context) {
	log.Info("ListProductProfiles")
	// username := c.GetUesrName()
	profiles := v1.VCluster{
		Name:    "vcluster-test",
		Version: "v1.0.0",
	}
	response.Response(c, response.SuccessGCPResponse, &profiles)
}
