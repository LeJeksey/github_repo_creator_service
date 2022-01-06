package services

import (
	"github.com/stretchr/testify/assert"
	"golang-microservices/src/api/domain/clients/restclient"
	"golang-microservices/src/api/domain/repositories"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCreateRepoInvalidName(t *testing.T) {
	request := repositories.CreateRepoRequest{}

	res, err := RepositoryService.CreateRepo(request)
	assert.Nil(t, res)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "invalid repository name", err.Message())
}

func TestCreateRepoErrorFromGithub(t *testing.T) {
	restclient.StartMock()
	defer restclient.StopMock()
	restclient.AddMock(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body: io.NopCloser(
				strings.NewReader(
					`{"message": "Requires authentication", "documentation_url": "https://developer.github.com/"}`,
				),
			),
		},
	})

	request := repositories.CreateRepoRequest{Name: "testing_repo"}

	res, err := RepositoryService.CreateRepo(request)
	assert.Nil(t, res)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, "Requires authentication", err.Message())
}

func TestCreateRepoNoError(t *testing.T) {
	restclient.StartMock()
	defer restclient.StopMock()

	restclient.AddMock(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body: io.NopCloser(strings.NewReader(`{
"id": 2304923,
"name": "testing_repo",
"owner": {
	"login": "LeJeksey"	
}}`)),
		},
	})

	request := repositories.CreateRepoRequest{Name: "testing_repo"}

	res, err := RepositoryService.CreateRepo(request)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.EqualValues(t, 2304923, res.Id)
	assert.EqualValues(t, request.Name, res.Name)
	assert.EqualValues(t, "LeJeksey", res.Owner)
}
