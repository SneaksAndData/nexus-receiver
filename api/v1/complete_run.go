package v1

import (
	"github.com/SneaksAndData/nexus-core/pkg/checkpoint/request"
	"github.com/gin-gonic/gin"
)

func CompleteRun(checkpointStore *request.CqlStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: log errors
		algorithmName := ctx.Param("algorithmName")
		requestId := ctx.Param("requestId")

		//if err := ctx.ShouldBindJSON(&payload); err != nil {
		//	ctx.String(http.StatusBadRequest, `Algorithm payload is invalid: %s`, err.Error())
		//	return
		//}
		//
		//config, cacheErr := configCache.GetConfiguration(algorithmName)
		//
		//if cacheErr != nil {
		//	ctx.String(http.StatusInternalServerError, `Internal error occurred when processing your request.`, algorithmName)
		//	return
		//}
		//
		//if config == nil {
		//	ctx.String(http.StatusBadRequest, `No valid configuration found for: %s. Please check that algorithm name is spelled correctly and try again. Contact an algorithm author if this problem persists.`, algorithmName)
		//	return
		//}
		//
		//if err := buffer.Add(requestId.String(), algorithmName, &payload, &config.Spec); err != nil {
		//	ctx.String(http.StatusBadRequest, `Request buffering failed for: %s, error: %s`, requestId.String(), err.Error())
		//	return
		//}
		//
		//ctx.JSON(http.StatusAccepted, map[string]string{
		//	"requestId": requestId.String(),
		//})
	}
}
