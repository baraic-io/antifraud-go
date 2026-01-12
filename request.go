package antifraud

import (
	"net/http"
)

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

func (e CodeError) Error() string {
	return e.Msg
}

var (
	ErrUnauthorized = NewCodeError(401, "Unauthorized")
	ErrInvalidToken = NewCodeError(401, "Invalid or expired token")
)

func (c Client) doRequest(request *http.Request) (*http.Response, error) {
	request.Header.Add("X-API-Key", c.APIKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
