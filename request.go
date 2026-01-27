package antifraud

import (
	"fmt"
	"io"
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
	return fmt.Sprintf("request failed with status code: %d, response: %s", e.Code, e.Msg)
}

var (
	ErrUnauthorized = NewCodeError(401, "Unauthorized")
	ErrInvalidToken = NewCodeError(401, "Invalid or expired token")
)

func (c Client) doRequest(request *http.Request) ([]byte, error) {
	request.Header.Add("X-API-Key", c.APIKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, CodeError{Code: response.StatusCode, Msg: fmt.Errorf("failed to read response body: %s", err.Error()).Error()}
		}

		return nil, CodeError{Code: response.StatusCode, Msg: string(responseBody)}
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
