package antifraud

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

const (
	SyncMode = iota
	AsyncMode
)

type FinalResolution SyncResolution

type Transaction struct {
	Id string `json:"id"`

	Type        string `json:"type"` // deposit, withdraw
	Date        string `json:"date"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description,optional"`

	ClientId         string `json:"client_id,optional"`
	ClientName       string `json:"client_name,optional"`
	ClientPAN        string `json:"client_pan"`
	ClientCVV        string `json:"client_card_cvv,optional"`
	ClientCardHolder string `json:"client_card_holder,optional"`
	ClientCountry    string `json:"client_country,optional"`
	ClientCity       string `json:"client_city,optional"`
	ClientPhone      string `json:"client_phone,optional"`

	MerchantId         string `json:"merchant_id,optional"`
	MerchantDescriptor string `json:"merchant_descriptor,optional"`
	MerchantTerminalId string `json:"merchant_terminal_id"`
	MerchantCountry    string `json:"merchant_country,optional"`

	Channel         string `json:"channel,optional"` // E-com, mobile, etc
	LocationIp      string `json:"location_ip,optional"`
	LocationCountry string `json:"location_country,optional"`
}

type SyncResolution struct {
	AF_Id               string                 `json:"af_id"`
	AF_Transaction      AF_Transaction         `json:"af_transaction"`
	Id                  string                 `json:"id"`
	Error               string                 `json:"af_error,omitempty"`
	Details             map[string]interface{} `json:"af_details,omitempty"`
	AddDate             time.Time              `json:"af_add_date"`
	FinalizedDate       time.Time              `json:"af_finalized_date"`
	FinalizedAction     string                 `json:"af_finalized_action"`
	ProcessTime         int64                  `json:"af_process_time"`
	ValidatedServices   []string               `json:"af_validated_services"`
	UnvalidatedServices []string               `json:"af_unvalidated_services"`
	Retry               uint                   `json:"af_retry_attempt"`
	Fraud               bool                   `json:"af_fraud"`
	Validated           bool                   `json:"af_validated"`
	Blocked             bool                   `json:"af_blocked"`
	Alert               bool                   `json:"af_alert"`
}

type AsyncResolution struct {
	AF_Id      string `json:"af_id"`
	AF_AddDate string `json:"af_add_date"`
}

type ServiceResolution struct {
	AF_Id       string              `json:"af_id"`
	TxnId       string              `json:"txn_id"`
	Id          uuid.UUID           `json:"id"`
	Date        time.Time           `json:"date"`
	Service     string              `json:"service"`
	Error       string              `json:"error,omitempty"`
	Details     map[string]string   `json:"details,omitempty"`
	Action      map[string]struct{} `json:"action,omitempty"`
	ProcessTime int64               `json:"process_time"`
	Retry       uint                `json:"retry"`
	Fraud       bool                `json:"fraud"`
	Validated   bool                `json:"validated"`
	Blocked     bool                `json:"blocked"`
	Alert       bool                `json:"alert"`
	InWhiteList bool                `json:"in_white_list"`
}

type AF_Retry struct {
	RetryCount int `json:"retry_count"`
	RetryMax   int `json:"retry_max"`
}

type AF_Transaction struct {
	Transaction Transaction         `json:"transaction,omitempty"`
	AF_Id       string              `json:"af_id"`
	AF_AddDate  string              `json:"af_add_date"`
	AF_Retries  map[string]AF_Retry `json:"af_retries,omitempty"`
}
