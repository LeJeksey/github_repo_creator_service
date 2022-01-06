package github_provider

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"golang-microservices/src/api/domain/clients/restclient"
	"golang-microservices/src/api/domain/github"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestGetAuthHeader(t *testing.T) {
	header := getAuthHeader("abc123")
	assert.EqualValues(t, "token abc123", header)
}

func TestCreateRepoErrorRestclient(t *testing.T) {
	restclient.StartMock()
	defer restclient.StopMock()
	clientMock := &restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Error:      errors.New("invalid restclient response"),
	}
	restclient.AddMock(clientMock)

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, clientMock.Error.Error(), err.Message)
}

func TestCreateRepoInvalidResponseBody(t *testing.T) {
	restclient.StartMock()
	defer restclient.StopMock()

	invalidBody, _ := os.Open("-;dlfaksd;fasdf")
	clientMock := &restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       invalidBody,
		},
	}
	restclient.AddMock(clientMock)

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid response body", err.Message)
}

func TestCreateRepoInvalidJsonResponseBody(t *testing.T) {
	restclient.StartMock()
	defer restclient.StopMock()

	clientMock := &restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(strings.NewReader(`{"message": 1}`)),
		},
	}
	restclient.AddMock(clientMock)

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid json response body", err.Message)
}

func TestCreateRepoUnauthorized(t *testing.T) {
	restclient.StartMock()
	defer restclient.StopMock()

	clientMock := &restclient.Mock{
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
	}
	restclient.AddMock(clientMock)

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusUnauthorized, err.StatusCode)
	assert.EqualValues(t, "Requires authentication", err.Message)
}

func TestCreateRepoInvalidSuccessResponse(t *testing.T) {
	restclient.StartMock()
	defer restclient.StopMock()

	clientMock := &restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body: io.NopCloser(
				strings.NewReader(
					`{"id": "123"}`,
				),
			),
		},
	}
	restclient.AddMock(clientMock)

	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, response)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "error unmarshalling response body", err.Message)
}

func TestCreateRepoOk(t *testing.T) {
	restclient.StartMock()
	defer restclient.StopMock()

	clientMock := &restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body: io.NopCloser(
				strings.NewReader(
					`{
  "id": 2304923,
  "node_id": "MSKADLAKSDJ0192eladkj",
  "name": "my-repo",
  "full_name": "myrepo/my-repo",
  "private": false,
  "owner": {
    "id": 1234,
    "url": "https://github.com/myrepo/my",
    "html_url": "<a>https://github.com/myrepo/my</a>",
    "login": "Ivanov"
  },
  "permissions": {
    "is_admin": true,
    "pull": true,
    "push": false
  }
}`,
				),
			),
		},
	}
	restclient.AddMock(clientMock)

	expectedResponse := &github.CreateRepoResponse{
		Id:       2304923,
		Name:     "my-repo",
		FullName: "myrepo/my-repo",
		Owner: github.RepoOwner{
			Id:      1234,
			Login:   "Ivanov",
			Url:     "https://github.com/myrepo/my",
			HtmlUrl: "<a>https://github.com/myrepo/my</a>",
		},
		Permissions: github.RepoPermissions{
			IsAdmin: true,
			HasPull: true,
			HasPush: false,
		},
	}
	response, err := CreateRepo("", github.CreateRepoRequest{})

	assert.Nil(t, err)
	assert.NotNil(t, response)

	assert.EqualValues(t, expectedResponse, response)
}
