package config

import (
	"log"
	"os"
)

const apiGithubAccessToken = "SECRET_API_GITHUB_ACCESS_TOKEN"

var githubAccessToken = os.Getenv(apiGithubAccessToken)

func init() {
	if githubAccessToken == "" {
		log.Println("WARNING: githubAccessToken is empty")
	}
}

func GetGithubAccessToken() string {
	return githubAccessToken
}
