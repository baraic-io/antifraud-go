package antifraud

import (
	"bytes"
	"encoding/json"
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
