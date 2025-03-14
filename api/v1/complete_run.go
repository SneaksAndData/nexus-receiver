package v1

import (
	"github.com/SneaksAndData/nexus-receiver/api/v1/models"
	"github.com/SneaksAndData/nexus-receiver/app"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CompleteRun(actor *app.CompletionActor) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: log errors
		algorithmName := ctx.Param("algorithmName")
		requestId := ctx.Param("requestId")
		var result models.AlgorithmResult

		if err := ctx.ShouldBindJSON(&result); err != nil {
			ctx.String(http.StatusBadRequest, `Submitted result is invalid: %s`, err.Error())
			return
		}

		actor.Receive(&models.CompletionInput{Result: result, RequestId: requestId, AlgorithmName: algorithmName})

		ctx.JSON(http.StatusAccepted, map[string]string{})
	}
}
