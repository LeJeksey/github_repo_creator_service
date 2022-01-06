package repositories

import (
	"github.com/gin-gonic/gin"
	"golang-microservices/src/api/domain/repositories"
	"golang-microservices/src/api/services"
	"golang-microservices/src/api/utils/errors"
	"net/http"
)

func CreateRepo(ctx *gin.Context) {
	var request repositories.CreateRepoRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequestApiError("invalid json body")
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	res, err := services.RepositoryService.CreateRepo(request)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func CreateRepos(ctx *gin.Context) {
	var request []repositories.CreateRepoRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequestApiError("invalid json body")
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	res, err := services.RepositoryService.CreateRepos(request)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(res.StatusCode, res)
}
