package services

import (
	"golang-microservices/src/api/config"
	"golang-microservices/src/api/domain/github"
	"golang-microservices/src/api/domain/repositories"
	"golang-microservices/src/api/providers/github_provider"
	"golang-microservices/src/api/utils/errors"
	"net/http"
	"sync"
)

type reposService struct{}

type ReposServiceInterface interface {
	CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(input []repositories.CreateRepoRequest) (*repositories.CreateReposResponse, errors.ApiError)
}

var RepositoryService ReposServiceInterface

func init() {
	RepositoryService = &reposService{}
}

func (s *reposService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}

	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)
	if err != nil {
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	res := repositories.CreateRepoResponse{Id: response.Id, Name: response.Name, Owner: response.Owner.Login}
	return &res, nil
}

func (s *reposService) CreateRepos(requests []repositories.CreateRepoRequest) (*repositories.CreateReposResponse, errors.ApiError) {
	input := make(chan *repositories.CreateRepositoresResult)
	output := make(chan *repositories.CreateReposResponse)
	defer close(output)

	var wg sync.WaitGroup
	go s.handleRepoResults(&wg, input, output)

	for _, current := range requests {
		wg.Add(1)
		go s.createRepoConcurrent(current, input)
	}

	wg.Wait()
	close(input)

	result := <-output

	successCreations := 0
	for _, current := range result.Results {
		if current.Response != nil {
			successCreations++
		}
	}

	switch true {
	case successCreations == 0:
		result.StatusCode = result.Results[0].Error.Status()
	case successCreations == len(requests):
		result.StatusCode = http.StatusCreated
	default:
		result.StatusCode = http.StatusPartialContent
	}

	return result, nil
}

func (s *reposService) handleRepoResults(wg *sync.WaitGroup, input chan *repositories.CreateRepositoresResult, output chan *repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse

	for result := range input {
		results.Results = append(results.Results, *result)
		wg.Done()
	}

	output <- &results
}

func (s *reposService) createRepoConcurrent(input repositories.CreateRepoRequest, output chan *repositories.CreateRepositoresResult) {
	res, err := s.CreateRepo(input)

	output <- &repositories.CreateRepositoresResult{Response: res, Error: err}
}
