package v1

import (
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/request"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"net/http"
)

// CheckRun godoc
//
//	@Summary		Verifies the run has been modified by receiver.
//	@Description	Checks if one of the completion statuses has been assigned by receiver.
//	@Tags			results
//	@Produce		json
//	@Produce		plain
//	@Produce		html
//	@Param			algorithmName	path		string	true	"Request id of the run to complete"
//	@Param			requestId	path		string	true	"Request id of the run to complete"
//	@Success		200	{object}    map[string]bool
//	@Failure		400	{string}	string
//	@Failure		404	{string}	string
//	@Failure		401	{string}	string
//	@Router			/algorithm/v1/check/{algorithmName}/requests/{requestId} [get]
func CheckRun(cqlStore *request.CqlStore, logger klog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		algorithmName := ctx.Param("algorithmName")
		requestId := ctx.Param("requestId")

		requestToCheck, err := cqlStore.ReadCheckpoint(algorithmName, requestId)
		if err != nil {
			logger.V(0).Error(err, "error when reading a checkpoint", "requestId", requestId, "algorithmName", algorithmName)
			ctx.String(http.StatusBadRequest, `unable to check status of a requested checkpoint: %s/%s`, algorithmName, requestId)
			return
		}

		if requestToCheck == nil {
			ctx.String(http.StatusNotFound, `checkpoint not found: %s/%s`, algorithmName, requestId)
			return
		}

		ctx.JSON(http.StatusOK, map[string]bool{
			"processed": requestToCheck.IsFinished(),
		})
	}
}
