package repositories

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang-microservices/src/api/domain/repositories"
	"golang-microservices/src/api/services"
	"golang-microservices/src/api/utils/errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateRepoInvalidJsonBody(t *testing.T) {
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/repositories", strings.NewReader(""))

	CreateRepo(ctx)

	resBody := response.Body.Bytes()
	resError, _ := errors.NewApiErrorFromBytes(resBody)

	assert.EqualValues(t, http.StatusBadRequest, response.Code)
	assert.EqualValues(t, http.StatusBadRequest, resError.Status())
	assert.EqualValues(t, "invalid json body", resError.Message())
	assert.EqualValues(t, "", resError.Error())
}

type reposServiceMock struct{}

func (r *reposServiceMock) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	return createRepoFunc(input)
}

var createRepoFunc func(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)

func TestCreateRepoErrorFromGithub(t *testing.T) {
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/repositories",
		strings.NewReader(`{"name": "", "description": "test description"}`),
	)

	var actualCreateRepoInput repositories.CreateRepoRequest
	createRepoFunc = func(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
		actualCreateRepoInput = input
		return nil, errors.NewBadRequestApiError("invalid repository name")
	}
	services.RepositoryService = &reposServiceMock{}

	CreateRepo(ctx)

	assert.EqualValues(
		t,
		repositories.CreateRepoRequest{Name: "", Description: "test description"},
		actualCreateRepoInput,
	)

	resBody := response.Body.Bytes()
	resError, err := errors.NewApiErrorFromBytes(resBody)
	if err != nil {
		log.Println("error while creating error from testing response", err)
	}

	assert.EqualValues(t, http.StatusBadRequest, response.Code)
	assert.EqualValues(t, http.StatusBadRequest, resError.Status())
	assert.EqualValues(t, "invalid repository name", resError.Message())
	assert.EqualValues(t, "", resError.Error())
}

func TestCreateRepoNoError(t *testing.T) {
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/repositories",
		strings.NewReader(`{"name": "test_repo", "description": "test description"}`),
	)

	var actualCreateRepoInput repositories.CreateRepoRequest
	expectedResponse := &repositories.CreateRepoResponse{Id: 123, Name: "test_repo", Owner: "SomeOwnerOfGithubToken"}
	createRepoFunc = func(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
		actualCreateRepoInput = input
		return expectedResponse, nil
	}
	services.RepositoryService = &reposServiceMock{}

	CreateRepo(ctx)

	assert.EqualValues(
		t,
		repositories.CreateRepoRequest{Name: "test_repo", Description: "test description"},
		actualCreateRepoInput,
	)

	var res repositories.CreateRepoResponse
	if err := json.Unmarshal(response.Body.Bytes(), &res); err != nil {
		log.Println("error while unmarshalling testing response", err)
	}

	assert.EqualValues(t, http.StatusCreated, response.Code)
	assert.EqualValues(t, *expectedResponse, res)
}
