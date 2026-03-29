package antifraud

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotSupported   = errors.New("not supported")
	ErrFieldNotFound  = errors.New("field not found")
	ErrNotImplemented = errors.New("not implemented")
)

const (
	TransactionTypeDeposit  = "deposit"
	TransactionTypeWithdraw = "withdraw"

	ClientTypePerson       = "person"
	ClientTypeOrganization = "organization"

	ChannelBinance = "binance"
	ChannelEcom    = "e-com"
	ChannelMobile  = "mobile"
)

const (
	SyncMode = iota
	AsyncMode
)

type FinalResolution SyncResolution

type Transaction struct {
	Id string `json:"id"`

	Type        string `json:"type"`             // deposit, withdraw
	Channel     string `json:"channel,optional"` // e-com, mobile, binance
	Date        string `json:"date"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description,optional"`

	LocationIp      string `json:"location_ip,optional"`
	LocationCountry string `json:"location_country,optional"`

	SenderId             string `json:"sender_id,optional"`
	SenderType           string `json:"sender_type,optional"` // person, organization, merchant
	SenderName           string `json:"sender_name,optional"`
	SenderPAN            string `json:"sender_pan,optional"`
	SenderCVV            string `json:"sender_card_cvv,optional"`
	SenderCardHolder     string `json:"sender_card_holder,optional"`
	SenderCardYearMonth  string `json:"sender_card_year_month,optional"`
	SenderContractNumber string `json:"sender_contract_number,optional"`
	SenderCountry        string `json:"sender_country,optional"`
	SenderPhone          string `json:"sender_phone,optional"`

	RecipientId             string `json:"recipient_id,optional"`
	RecipientType           string `json:"recipient_type,optional"`
	RecipientName           string `json:"recipient_name,optional"`
	RecipientPAN            string `json:"recipient_pan,optional"`
	RecipientCVV            string `json:"recipient_card_cvv,optional"`
	RecipientCardHolder     string `json:"recipient_card_holder,optional"`
	RecipientCardYearMonth  string `json:"recipient_card_year_month,optional"`
	RecipientContractNumber string `json:"recipient_contract_number,optional"`
	RecipientCountry        string `json:"recipient_country,optional"`
	RecipientPhone          string `json:"recipient_phone,optional"`
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

type ValidatedTransaction struct {
	Transaction Transaction `json:"transaction"`
	Decision    int         `json:"decision"` // 0 - not fraud, 1 - fraud
}

type RetrainLog struct {
	Timestamp   string
	DatasetSize int
	AUC         float64
	Precision   float64
	Recall      float64
	F1          float64
	TP          int
	FP          int
	FN          int
	TN          int
	FPR         float64
}
