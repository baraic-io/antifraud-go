package antifraud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

/* Validate Transaction in Sync mode */
func (c Client) ValidateTransactionSync(transaction Transaction) (SyncResolution, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ClientConfig.ValidationCtxDeadlineTimeout)*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(transaction)
	if err != nil {
		return SyncResolution{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.Host+"/api/gtwsvc/sync/transaction", bytes.NewBuffer(jsonData))
	if err != nil {
		return SyncResolution{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := c.doRequest(req)
	if err != nil {
		return SyncResolution{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return SyncResolution{}, fmt.Errorf("failed to read response body: %s", err.Error())
		}

		return SyncResolution{}, CodeError{Code: response.StatusCode, Msg: string(responseBody)}
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return SyncResolution{}, err
	}

	var resolution SyncResolution
	if err := json.Unmarshal(responseBody, &resolution); err != nil {
		return SyncResolution{}, err
	}

	return resolution, nil
}

/* Validate Transaction in Async mode */
func (c Client) ValidateTransactionAsync(transaction Transaction) (AsyncResolution, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ClientConfig.ValidationCtxDeadlineTimeout)*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(transaction)
	if err != nil {
		return AsyncResolution{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.Host+"/api/gtwsvc/async/transaction", bytes.NewBuffer(jsonData))
	if err != nil {
		return AsyncResolution{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := c.doRequest(req)
	if err != nil {
		return AsyncResolution{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return AsyncResolution{}, fmt.Errorf("failed to read response body: %s", err.Error())
		}

		return AsyncResolution{}, CodeError{Code: response.StatusCode, Msg: string(responseBody)}
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return AsyncResolution{}, err
	}

	var resolution AsyncResolution
	if err := json.Unmarshal(responseBody, &resolution); err != nil {
		return AsyncResolution{}, err
	}

	return resolution, nil
}

/* Validate Transaction by AML Service */
func (c Client) ValidateTransactionByAML(af_transaction AF_Transaction) (ServiceResolution, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ClientConfig.ValidationCtxDeadlineTimeout)*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(af_transaction)
	if err != nil {
		return ServiceResolution{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.Host+"/api/amlsvc/validate", bytes.NewBuffer(jsonData))
	if err != nil {
		return ServiceResolution{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := c.doRequest(req)
	if err != nil {
		return ServiceResolution{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return ServiceResolution{}, fmt.Errorf("failed to read response body: %s", err.Error())
		}

		return ServiceResolution{}, CodeError{Code: response.StatusCode, Msg: string(responseBody)}
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return ServiceResolution{}, err
	}

	var resolution ServiceResolution
	if err := json.Unmarshal(responseBody, &resolution); err != nil {
		return ServiceResolution{}, err
	}

	return resolution, nil
}

/* Validate Transaction by Custom Rules (Not implemented) */
func (c Client) ValidateTransactionByRules(af_transaction AF_Transaction) (ServiceResolution, error) {
	return ServiceResolution{}, ErrNotImplemented
}
