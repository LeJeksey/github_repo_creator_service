package restclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var enabledMocks = false
var mocks = make(map[string]*Mock)

type Mock struct {
	Url        string
	HttpMethod string
	Response   *http.Response
	Error      error
}

func StartMock() {
	enabledMocks = true
}

func StopMock() {
	enabledMocks = false
	mocks = make(map[string]*Mock)
}

func AddMock(mock *Mock) {
	mocks[mock.Url] = mock
}

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	if enabledMocks {
		mock := mocks[url]
		if mock == nil {
			return nil, errors.New("no mock found for given url")
		}
		return mock.Response, mock.Error
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	request.Header = headers

	client := http.Client{}
	return client.Do(request)

}
