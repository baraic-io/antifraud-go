package antifraud

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

/* Validate transaction by any validation service */
func (c Client) validate(servicePath string, af_transaction AF_Transaction) (ServiceResolution, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ClientConfig.ValidationCtxDeadlineTimeout)*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(af_transaction)
	if err != nil {
		return ServiceResolution{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.Host+servicePath, bytes.NewBuffer(jsonData))
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

/* Validate Transaction in Sync mode */
// func (c Client) validateTransactionSync(transaction Transaction) (SyncResolution, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ClientConfig.ValidationCtxDeadlineTimeout)*time.Second)
// 	defer cancel()

// 	jsonData, err := json.Marshal(transaction)
// 	if err != nil {
// 		return SyncResolution{}, err
// 	}

// 	req, err := http.NewRequestWithContext(ctx, "POST", c.Host+"/api/gtwsvc/sync/transaction", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return SyncResolution{}, err
// 	}

// 	req.Header.Add("Content-Type", "application/json")

// 	response, err := c.doRequest(req)
// 	if err != nil {
// 		return SyncResolution{}, err
// 	}

// 	var resolution SyncResolution
// 	if err := json.Unmarshal(response, &resolution); err != nil {
// 		return SyncResolution{}, err
// 	}

// 	return resolution, nil
// }

/* Validate Transaction in Async mode */
// func (c Client) validateTransactionAsync(transaction Transaction) (AsyncResolution, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ClientConfig.ValidationCtxDeadlineTimeout)*time.Second)
// 	defer cancel()

// 	jsonData, err := json.Marshal(transaction)
// 	if err != nil {
// 		return AsyncResolution{}, err
// 	}

// 	req, err := http.NewRequestWithContext(ctx, "POST", c.Host+"/api/gtwsvc/async/transaction", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return AsyncResolution{}, err
// 	}

// 	req.Header.Add("Content-Type", "application/json")

// 	response, err := c.doRequest(req)
// 	if err != nil {
// 		return AsyncResolution{}, err
// 	}

// 	var resolution AsyncResolution
// 	if err := json.Unmarshal(response, &resolution); err != nil {
// 		return AsyncResolution{}, err
// 	}

// 	return resolution, nil
// }

/* Validate Transaction by AML Service */
func (c Client) ValidateTransactionByAML(af_transaction AF_Transaction) (ServiceResolution, error) {
	resolution, err := c.validate("/api/amlsvc/validate", af_transaction)
	if err != nil {
		return ServiceResolution{}, err
	}

	return resolution, nil
}

/* Validate Transaction by FC service */
func (c Client) ValidateTransactionByFC(af_transaction AF_Transaction) (ServiceResolution, error) {
	resolution, err := c.validate("/api/fcsvc/validate", af_transaction)
	if err != nil {
		return ServiceResolution{}, err
	}

	return resolution, nil
}

/* Validate Transaction by ML service */
func (c Client) ValidateTransactionByML(af_transaction AF_Transaction) (ServiceResolution, error) {
	resolution, err := c.validate("/api/mlsvc/validate", af_transaction)
	if err != nil {
		return ServiceResolution{}, err
	}

	return resolution, nil
}

/* DEPRECATED: Validate Transaction by LST service */
func (c Client) validateTransactionByLST(af_transaction AF_Transaction) (ServiceResolution, error) {
	resolution, err := c.validate("/api/lstsvc/validate", af_transaction)
	if err != nil {
		return ServiceResolution{}, err
	}

	return resolution, nil
}
