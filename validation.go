package antifraud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
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

// Convert source transaction related on channel to AF transaction
func (c Client) ToAFTransaction(transaction map[string]interface{}) (Transaction, error) {
	/* TODO: Determine channel, type of transaction if new channels will added */
	/* Default: C2C2Out (mobile) */
	afTransaction := Transaction{Channel: ChannelMobile, Type: TransactionTypeDeposit}

	switch afTransaction.Channel {
	default:
		return Transaction{}, fmt.Errorf("Unsupported channel type: %s", afTransaction.Channel)
	case ChannelBinance:
		/* TODO: Determine transaction type */
		/* Default: deposit */
		transactionType := TransactionTypeDeposit

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
		/* Validate channel related on product id */
		finoper, ok := transaction["finoper"].(map[string]interface{})
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "finoper")
		}

		productId, ok := finoper["product_id"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "product_id")
		}

		if !slices.Contains(ChannelMobileProductIds[:], productId) {
			return Transaction{}, fmt.Errorf("product (%s) id not found in channel mobile product ids list", productId)
		}

		afTransaction.ProductId = productId

		/* Parse transaction */
		operId, ok := finoper["oper_id"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "oper_id")
		}
		afTransaction.Id = operId

		reqId, ok := transaction["req_id"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "req_id")
		}
		afTransaction.RequestId = reqId

		reason, ok := finoper["reason"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "reason")
		}
		afTransaction.Description = reason

		knpCode, ok := finoper["knp_code"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "knp_code")
		}
		afTransaction.KNPCode = knpCode

		operTypeFloat, ok := finoper["oper_type"].(float64)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "oper_type")
		}
		afTransaction.OperType = int(operTypeFloat)

		operDate, ok := finoper["oper_date_time"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "oper_date_time")
		}
		afTransaction.Date = operDate

		creationDate, ok := finoper["creation_date_time"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "creation_date_time")
		}
		afTransaction.CreationDate = creationDate

		ipAddr, ok := transaction["ip_addr"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "ip_addr")
		}
		afTransaction.LocationIp = ipAddr

		/* Parse participants (first is sender, second is recipient) */
		participantsRaw, ok := finoper["person"].([]interface{})
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "person")
		}
		var participants []map[string]interface{}
		for _, pRaw := range participantsRaw {
			if pm, ok := pRaw.(map[string]interface{}); ok {
				participants = append(participants, pm)
			}
		}

		participantsFinoperRaw, ok := finoper["finoper_dc"].([]interface{})
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "finoper_dc")
		}
		var participantsFinoper []map[string]interface{}
		for _, pRaw := range participantsFinoperRaw {
			if pm, ok := pRaw.(map[string]interface{}); ok {
				participantsFinoper = append(participantsFinoper, pm)
			}
		}

		if len(participants) < 2 {
			return Transaction{}, fmt.Errorf("not enough participants: %d", len(participants))
		}
		if len(participantsFinoper) < 1 {
			return Transaction{}, fmt.Errorf("not enough finoper_dc: %d", len(participantsFinoper))
		}

		/* Parse amount and currency from participants */
		amountFloat, ok := participantsFinoper[0]["amount_kzt"].(float64)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found: %s", "amount_kzt")
		}
		afTransaction.Amount = strconv.FormatFloat(amountFloat, 'f', 2, 64)

		/* TODO: Operate with currency code */
		afTransaction.Currency = "KZT"

		/* Parse sender */
		sender := participants[0]
		senderFinoper := participantsFinoper[0]

		senderIINBIN, ok := sender["iin"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "iin")
		}
		afTransaction.SenderIINBIN = senderIINBIN

		senderName, ok := sender["full_name"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "full_name")
		}
		afTransaction.SenderName = senderName

		senderTypeFloat, ok := sender["person_type"].(float64)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "person_type")
		}
		senderType := int(senderTypeFloat)

		switch senderType {
		default:
			return Transaction{}, fmt.Errorf("unsupported person type (sender): %d", senderType)
		case 1:
			afTransaction.SenderType = ClientTypeOrganization
		case 2:
			afTransaction.SenderType = ClientTypePerson
		}

		senderPhone, ok := sender["mobile_number"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "mobile_number")
		}
		afTransaction.SenderPhone = senderPhone

		senderRegDate, ok := sender["client_reg_date"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "client_reg_date")
		}
		afTransaction.SenderRegDate = senderRegDate

		senderBankBic, ok := senderFinoper["bank_bic"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "bank_bic")
		}
		afTransaction.SenderBankBic = senderBankBic

		senderBankName, ok := senderFinoper["bank_name"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "bank_name")
		}
		afTransaction.SenderBankName = senderBankName

		senderCardNumber, ok := senderFinoper["card_number"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "card_number")
		}
		afTransaction.SenderPAN = senderCardNumber

		senderCardExpDate, ok := senderFinoper["card_exp_date"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "card_exp_date")
		}
		afTransaction.SenderCardExpDate = senderCardExpDate

		senderCardOpenDate, ok := senderFinoper["card_open_date"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "card_open_date")
		}
		afTransaction.SenderCardOpenDate = senderCardOpenDate

		senderAccountNumber, ok := senderFinoper["account_number"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "account_number")
		}
		afTransaction.SenderContractNumber = senderAccountNumber

		senderIsClient, ok := sender["is_client"].(bool)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (sender): %s", "is_client")
		}
		afTransaction.SenderIsClient = senderIsClient

		/* Parse recipient */
		recipient := participants[1]
		recipientFinoper := participantsFinoper[1]

		recipientIINBIN, ok := recipient["iin"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "iin")
		}
		afTransaction.RecipientIINBIN = recipientIINBIN

		recipientName, ok := recipient["full_name"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "full_name")
		}
		afTransaction.RecipientName = recipientName

		recipientTypeFloat, ok := recipient["person_type"].(float64)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "person_type")
		}
		recipientType := int(recipientTypeFloat)

		switch recipientType {
		default:
			return Transaction{}, fmt.Errorf("unsupported person type (recipient): %d", recipientType)
		case 1:
			afTransaction.RecipientType = ClientTypeOrganization
		case 2:
			afTransaction.RecipientType = ClientTypePerson
		}

		recipientPhone, ok := recipient["mobile_number"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "mobile_number")
		}
		afTransaction.RecipientPhone = recipientPhone

		recipientRegDate, ok := recipient["client_reg_date"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "client_reg_date")
		}
		afTransaction.RecipientRegDate = recipientRegDate

		recipientBankBic, ok := recipientFinoper["bank_bic"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "bank_bic")
		}
		afTransaction.RecipientBankBic = recipientBankBic

		recipientBankName, ok := recipientFinoper["bank_name"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "bank_name")
		}
		afTransaction.RecipientBankName = recipientBankName

		recipientCardNumber, ok := recipientFinoper["card_number"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "card_number")
		}
		afTransaction.RecipientPAN = recipientCardNumber

		recipientCardExpDate, ok := recipientFinoper["card_exp_date"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "card_exp_date")
		}
		afTransaction.RecipientCardExpDate = recipientCardExpDate

		recipientCardOpenDate, ok := recipientFinoper["card_open_date"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "card_open_date")
		}
		afTransaction.RecipientCardOpenDate = recipientCardOpenDate

		recipientAccountNumber, ok := recipientFinoper["account_number"].(string)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "account_number")
		}
		afTransaction.RecipientContractNumber = recipientAccountNumber

		recipientIsClient, ok := recipient["is_client"].(bool)
		if !ok {
			return Transaction{}, fmt.Errorf("field not found (recipient): %s", "is_client")
		}
		afTransaction.RecipientIsClient = recipientIsClient
	case ChannelEcom:
		/* TODO: Determine transaction type */
		/* Default: deposit */
		transactionType := TransactionTypeDeposit

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

		description, _ := transaction["description"].(string)
		afTransaction.Description = description

		clientPAN, ok := transaction["pan"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: pan")
		}

		if transactionType == TransactionTypeDeposit {
			afTransaction.SenderPAN = clientPAN
		} else if transactionType == TransactionTypeWithdraw {
			afTransaction.RecipientPAN = clientPAN
		}

		/* Parsing client information */
		/* TODO: Is client information are required? */
		client, ok := transaction["client"].(map[string]interface{})
		if !ok {
			return Transaction{}, errors.New("field not found: client")
		}

		var id, name, phone, country string
		id, _ = client["id"].(string)
		name, _ = client["name"].(string)
		phone, _ = client["phone"].(string)
		country, _ = client["country"].(string)

		/* Parsing merchant */
		options, ok := transaction["options"].(map[string]interface{})
		if !ok {
			return Transaction{}, errors.New("field not found: options")
		}

		merchantTerminalId, ok := options["terminal"].(string)
		if !ok {
			return Transaction{}, errors.New("field not found: terminal")
		}

		if transactionType == TransactionTypeWithdraw {
			afTransaction.RecipientId = id
			afTransaction.RecipientName = name
			afTransaction.RecipientPhone = phone
			afTransaction.RecipientCountry = country
			afTransaction.RecipientType = ClientTypePerson

			afTransaction.SenderId = merchantTerminalId
			afTransaction.SenderType = ClientTypeOrganization
		} else if transactionType == TransactionTypeDeposit {
			afTransaction.SenderId = id
			afTransaction.SenderName = name
			afTransaction.SenderPhone = phone
			afTransaction.SenderCountry = country
			afTransaction.SenderType = ClientTypePerson

			afTransaction.RecipientId = merchantTerminalId
			afTransaction.RecipientType = ClientTypeOrganization
		}

		/* TODO: Is location are required? */
		location, ok := transaction["location"].(map[string]interface{})
		if !ok {
			return Transaction{}, errors.New("field not found: location")
		}

		afTransaction.LocationIp, _ = location["ip"].(string)
		afTransaction.LocationCountry, _ = location["country"].(string)
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
