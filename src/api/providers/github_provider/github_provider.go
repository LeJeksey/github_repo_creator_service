package github_provider

import (
	"encoding/json"
	"fmt"
	"golang-microservices/src/api/domain/clients/restclient"
	"golang-microservices/src/api/domain/github"
	"io/ioutil"
	"log"
	"net/http"
)

const headerAuthorization = "Authorization"
const headerAuthorizationFormat = "token %s"

const urlCreateRepo = "https://api.github.com/user/repos"

func getAuthHeader(accessToken string) string {
	return fmt.Sprintf(headerAuthorizationFormat, accessToken)
}

func CreateRepo(accessToken string, request github.CreateRepoRequest) (*github.CreateRepoResponse, *github.GithubErrorResponse) {
	headers := http.Header{}
	headers.Set(headerAuthorization, getAuthHeader(accessToken))

	response, err := restclient.Post(urlCreateRepo, request, headers)
	if err != nil {
		log.Println("error when trying to create repository in github", err)
		return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: err.Error()}
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("error when trying to read bytes from response body", err)
		return nil, &github.GithubErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "invalid response body",
		}
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		var errResponse github.GithubErrorResponse
		if err := json.Unmarshal(responseBytes, &errResponse); err != nil {
			return nil, &github.GithubErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "invalid json response body",
			}
		}

		errResponse.StatusCode = response.StatusCode
		return nil, &errResponse
	}

	var result github.CreateRepoResponse
	if err := json.Unmarshal(responseBytes, &result); err != nil {
		log.Println("error unmarshalling response body", err)
		return nil, &github.GithubErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error unmarshalling response body",
		}
	}

	return &result, nil
}
