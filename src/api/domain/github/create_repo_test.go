package github

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRepoRequestAsJson(t *testing.T) {
	request := CreateRepoRequest{
		Name:        "testName",
		Description: "My description",
		Homepage:    "https://github.com/",
		Private:     true,
		HasIssues:   false,
		HasProjects: false,
		HasWiki:     true,
	}

	expectedJson := "{\"name\":\"testName\",\"description\":\"My description\",\"homepage\":\"https://github.com/\",\"private\":true,\"has_issues\":false,\"has_projects\":false,\"has_wiki\":true}"
	bytes, err := json.Marshal(request)

	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	assert.EqualValues(t, expectedJson, string(bytes))
}
