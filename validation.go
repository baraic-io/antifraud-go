package antifraud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

// Convert source transaction related on channel to AF transaction
/* channel: e-com, binance, mobile */
/* transactionType: deposit, withdraw */
func (c Client) ToAFTransaction(channel, transactionType string, transaction map[string]interface{}) (Transaction, error) {
	afTransaction := Transaction{
		Channel: channel,
		Type:    transactionType,
	}

	switch channel {
	default:
		return Transaction{}, fmt.Errorf("Unsupported channel type: %s", channel)
	case ChannelBinance:
		orderId, ok := transaction["order_id"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: order_id")
		}

		afTransaction.Id = orderId

		clientId, ok := transaction["client_id"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: client_id")
		}

		if transactionType == TransactionTypeDeposit {
			afTransaction.SenderId = clientId
			afTransaction.SenderType = ClientTypePerson
		} else if transactionType == TransactionTypeWithdraw {
			afTransaction.RecipientId = clientId
			afTransaction.RecipientType = ClientTypePerson
		} else {
			return Transaction{}, ErrNotSupported
		}

		amount, ok := transaction["amount"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: amount")
		}

		afTransaction.Amount = amount

		currency, ok := transaction["currency"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: currency")
		}

		afTransaction.Currency = strings.ToUpper(currency)
	case ChannelMobile:
		transactionId, ok := transaction["transferId"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: transferId")
		}

		afTransaction.Id = transactionId

		transactionType, ok := transaction["type"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: type")
		}

		afTransaction.Type = transactionType

		senderPAN, ok := transaction["sender_pan"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: sender_pan")
		}

		afTransaction.SenderPAN = senderPAN

		senderCardYearMonth, ok := transaction["sender_yearmonth"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: sender_yearmonth")
		}

		afTransaction.SenderCardYearMonth = senderCardYearMonth

		recipientPAN, ok := transaction["recipient_pan"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: recipient_pan")
		}

		afTransaction.RecipientPAN = recipientPAN

		message, ok := transaction["message"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: message")
		}

		afTransaction.Description = message

		amount, ok := transaction["amount"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: amount")
		}

		afTransaction.Amount = amount
	case ChannelEcom:
		/* TODO: Parse client card details */
		transferId, ok := transaction["merchant_order_id"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: merchant_order_id")
		}

		afTransaction.Id = transferId

		amount, ok := transaction["amount"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: amount")
		}

		afTransaction.Amount = amount

		currency, ok := transaction["currency"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: currency")
		}

		afTransaction.Currency = strings.ToUpper(currency)

		description, ok := transaction["description"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: description")
		}

		afTransaction.Description = description

		clientPAN, ok := transaction["pan"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: pan")
		}

		/* TODO: Client fields validation */
		client, ok := transaction["client"].(map[string]interface{})
		if !ok {
			return Transaction{}, errors.New("field not found: client")
		}

		if transactionType == TransactionTypeDeposit {
			afTransaction.SenderPAN = clientPAN
		} else if transactionType == TransactionTypeWithdraw {
			afTransaction.RecipientPAN = clientPAN
		}

		// Map client details
		id := fmt.Sprintf("%v", client["id"])
		if id == "<nil>" {
			id = ""
		}
		name := fmt.Sprintf("%v", client["name"])
		if name == "<nil>" {
			name = ""
		}
		phone := fmt.Sprintf("%v", client["phone"])
		if phone == "<nil>" {
			phone = ""
		}
		country := fmt.Sprintf("%v", client["country"])
		if country == "<nil>" {
			country = ""
		}

		if transactionType == TransactionTypeWithdraw {
			afTransaction.RecipientId = id
			afTransaction.RecipientName = name
			afTransaction.RecipientPhone = phone
			afTransaction.RecipientCountry = country
			afTransaction.RecipientType = ClientTypePerson
		} else {
			afTransaction.SenderId = id
			afTransaction.SenderName = name
			afTransaction.SenderPhone = phone
			afTransaction.SenderCountry = country
			afTransaction.SenderType = ClientTypePerson
		}

		location, ok := transaction["location"].(map[string]interface{})
		if !ok {
			return Transaction{}, errors.New("field not found: location")
		}

		afTransaction.LocationIp = fmt.Sprintf("%v", location["ip"])
		afTransaction.LocationCountry = fmt.Sprintf("%v", location["country"])
	}

	return afTransaction, nil
}

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
