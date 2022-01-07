package services

import (
	"github.com/stretchr/testify/assert"
	"golang-microservices/src/api/domain/clients/restclient"
	"golang-microservices/src/api/domain/repositories"
	"golang-microservices/src/api/utils/errors"
	"io"
	"net/http"
	"strings"
	"sync"
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

func TestCreateRepoConcurrent(t *testing.T) {
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

	input := repositories.CreateRepoRequest{Name: "test", Description: "test description"}
	output := make(chan *repositories.CreateRepositoresResult)

	service := &reposService{}
	go service.createRepoConcurrent(input, output)

	res := <-output

	assert.Nil(t, res.Error)

	assert.EqualValues(t, 2304923, res.Response.Id)
	assert.EqualValues(t, "testing_repo", res.Response.Name)
	assert.EqualValues(t, "LeJeksey", res.Response.Owner)
}

func TestHandleRepoResult(t *testing.T) {
	input := make(chan *repositories.CreateRepositoresResult)
	output := make(chan *repositories.CreateReposResponse)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		input <- &repositories.CreateRepositoresResult{Error: errors.NewBadRequestApiError("invalid repository name")}
	}()

	service := &reposService{}
	go service.handleRepoResults(&wg, input, output)

	wg.Wait()
	close(input)

	result := <-output
	assert.NotNil(t, result)
	assert.EqualValues(t, 1, len(result.Results))
	assert.EqualValues(t, "invalid repository name", result.Results[0].Error.Message())
}

func TestCreateReposInvalidRequests(t *testing.T) {
	requests := []repositories.CreateRepoRequest{
		{Name: "   "},
		{},
	}

	res, err := RepositoryService.CreateRepos(requests)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
	assert.EqualValues(t, 2, len(res.Results))

	assert.Nil(t, res.Results[0].Response)
	assert.EqualValues(t, http.StatusBadRequest, res.Results[0].Error.Status())
	assert.EqualValues(t, "invalid repository name", res.Results[0].Error.Message())

	assert.Nil(t, res.Results[1].Response)
	assert.EqualValues(t, http.StatusBadRequest, res.Results[1].Error.Status())
	assert.EqualValues(t, "invalid repository name", res.Results[1].Error.Message())
}

func TestCreateReposOneSuccessOneFail(t *testing.T) {
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

	requests := []repositories.CreateRepoRequest{
		{Name: "testing_repo"},
		{},
	}

	res, err := RepositoryService.CreateRepos(requests)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.EqualValues(t, http.StatusPartialContent, res.StatusCode)
	assert.EqualValues(t, 2, len(res.Results))

	countErr := 0
	countOk := 0
	for _, reqResult := range res.Results {
		if reqResult.Error != nil {
			assert.EqualValues(t, http.StatusBadRequest, reqResult.Error.Status())
			assert.EqualValues(t, "invalid repository name", reqResult.Error.Message())
			countErr++
		} else {
			assert.EqualValues(t, 2304923, reqResult.Response.Id)
			assert.EqualValues(t, "testing_repo", reqResult.Response.Name)
			assert.EqualValues(t, "LeJeksey", reqResult.Response.Owner)
			countOk++
		}
	}
	assert.EqualValues(t, 1, countErr)
	assert.EqualValues(t, 1, countOk)
}

func TestCreateReposOnlySuccess(t *testing.T) {
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

	requests := []repositories.CreateRepoRequest{
		{Name: "testing_repo"},
	}

	res, err := RepositoryService.CreateRepos(requests)

	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.EqualValues(t, http.StatusCreated, res.StatusCode)
	assert.EqualValues(t, 1, len(res.Results))

	assert.Nil(t, res.Results[0].Error)
	assert.EqualValues(t, 2304923, res.Results[0].Response.Id)
	assert.EqualValues(t, "testing_repo", res.Results[0].Response.Name)
	assert.EqualValues(t, "LeJeksey", res.Results[0].Response.Owner)
}
