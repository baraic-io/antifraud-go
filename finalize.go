package antifraud

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (c Client) AddTransactionServiceCheck(resolution ServiceResolution) error {
	jsonData, err := json.Marshal(resolution)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Host+"/api/fzrsvc/transaction/add-service-check", bytes.NewBuffer(jsonData))
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

func (c Client) FinalizeTransaction(af_transaction AF_Transaction) (FinalResolution, error) {
	jsonData, err := json.Marshal(af_transaction)
	if err != nil {
		return FinalResolution{}, err
	}

	req, err := http.NewRequest("POST", c.Host+"/api/fzrsvc/transaction/finalize", bytes.NewBuffer(jsonData))
	if err != nil {
		return FinalResolution{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := c.doRequest(req)
	if err != nil {
		return FinalResolution{}, err
	}

	var resolution FinalResolution
	if err := json.Unmarshal(response, &resolution); err != nil {
		return FinalResolution{}, err
	}

	return resolution, nil
}
