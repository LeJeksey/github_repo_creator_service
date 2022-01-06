package app

import (
	"golang-microservices/src/api/controllers/marcopolo"
	"golang-microservices/src/api/controllers/repositories"
)

func mapUrls() {
	router.GET("/marco", marcopolo.Marco)
	router.POST("/repository", repositories.CreateRepo)
	router.POST("/repositories", repositories.CreateRepos)
}
