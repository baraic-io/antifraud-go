package antifraud

import (
	"bytes"
	"context"
	"encoding/json"
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

	var resolution SyncResolution
	if err := json.Unmarshal(response, &resolution); err != nil {
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

	var resolution AsyncResolution
	if err := json.Unmarshal(response, &resolution); err != nil {
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

	var resolution ServiceResolution
	if err := json.Unmarshal(response, &resolution); err != nil {
		return ServiceResolution{}, err
	}

	return resolution, nil
}

/* Validate Transaction by Custom Rules */
func (c Client) ValidateTransactionByRules(af_transaction AF_Transaction) (ServiceResolution, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ClientConfig.ValidationCtxDeadlineTimeout)*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(af_transaction)
	if err != nil {
		return ServiceResolution{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.Host+"/api/fcsvc/validate", bytes.NewBuffer(jsonData))
	if err != nil {
		return ServiceResolution{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := c.doRequest(req)
	if err != nil {
		return ServiceResolution{}, err
	}

	var resolution ServiceResolution
	if err := json.Unmarshal(response, &resolution); err != nil {
		return ServiceResolution{}, err
	}

	return resolution, nil
}
