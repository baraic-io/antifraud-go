package antifraud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c Client) StoreTransaction(af_transaction AF_Transaction) error {
	jsonData, err := json.Marshal(af_transaction)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Host+"/api/storagesvc/store/transaction", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) StoreServiceResolution(resolution ServiceResolution) error {
	jsonData, err := json.Marshal(resolution)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Host+"/api/storagesvc/store/service-resolution", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) StoreFinalResolution(resolution FinalResolution) error {
	jsonData, err := json.Marshal(resolution)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Host+"/api/storagesvc/store/final-resolution", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

/* Store validated transaction for ML dataset */
func (c Client) StoreValidatedTransactionToML(transaction Transaction, decision int) error {
	jsonData, err := json.Marshal(ValidatedTransaction{Transaction: transaction, Decision: decision})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Host+"/api/mlsvc/store", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	fmt.Println(req)

	req.Header.Add("Content-Type", "application/json")

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

/* Retrain ML model */
func (c Client) RetrainModel() (RetrainLog, error) {
	req, err := http.NewRequest("GET", c.Host+"/api/mlsvc/retrain", nil)
	if err != nil {
		return RetrainLog{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := c.doRequest(req)
	if err != nil {
		return RetrainLog{}, err
	}

	var log RetrainLog
	if err := json.Unmarshal(response, &log); err != nil {
		return RetrainLog{}, err
	}

	return log, nil
}
