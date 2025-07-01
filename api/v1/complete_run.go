package v1

import (
	"github.com/SneaksAndData/nexus-receiver/api/v1/models"
	"github.com/SneaksAndData/nexus-receiver/app"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CompleteRun godoc
//
//	@Summary		Completes the provided Nexus algorithm run
//	@Description	Commits the run result to the checkpoint store and transitions the state to COMPLETED
//	@Tags			results
//	@Produce		json
//	@Produce		plain
//	@Produce		html
//	@Param			algorithmName	path		string	true	"Request id of the run to complete"
//	@Param			requestId	path		string	true	"Request id of the run to complete"
//	@Param			payload	body		models.AlgorithmResult	true	"Run result"
//	@Success		202	{object}    map[string]string
//	@Failure		400	{string}	string
//	@Failure		404	{string}	string
//	@Failure		401	{string}	string
//	@Router			/algorithm/v1.2/complete/{algorithmName}/requests/{requestId} [post]
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
