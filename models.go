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

	SourceUserId     string `json:"source_user_id"`
	SourceIdentifier string `json:"source_identifier"`
	SourceFullname   string `json:"source_fullname"`
	SourceCardNumber string `json:"source_card_number"`
	SourceAccount    string `json:"source_account"`

	TargetUserId     string `json:"target_user_id"`
	TargetIdentifier string `json:"target_identifier"`
	TargetFullname   string `json:"target_fullname"`
	TargetCardNumber string `json:"target_card_number"`
	TargetAccount    string `json:"target_account"`

	MerchantId         string `json:"merchant_id"`
	MerchantTerminalId string `json:"merchant_terminal_id"`
	MerchantMCCCode    string `json:"merchant_mcc_code"`

	Date               string `json:"date"`
	Time               string `json:"time"`
	Amount             string `json:"amount"`
	Currency           string `json:"currency"`
	PaymentMode        string `json:"payment_mode"`
	TransactionType    string `json:"transaction_type"`
	TransactionCountry string `json:"transaction_country"`
	TransactionCity    string `json:"transaction_city"`
	TransactionChannel string `json:"transaction_channel"`
	TransactionRRN     string `json:"transaction_rrn"`
	TransactionStatus  string `json:"transaction_status,optional"`

	RegistrationDate string `json:"registration_date,optional"`
	CardType         string `json:"card_type"`

	NewRecipient string `json:"new_recipient,optional"`
	NewTerminal  string `json:"new_terminal,optional"`

	DeviceId             string `json:"device_id" aml:"device_id"`
	LastDeviceUpdateDate string `json:"last_device_update_date,optional"`
	IPConnentionType     string `json:"ip_connection_type,optional"`
	RemoteAccess         string `json:"remote_access,optional"`
	ScreenSharing        string `json:"screen_sharing,optional"`
	HardwareId           string `json:"hardware_id,optional"`
	OSID                 string `json:"os_id,optional"`
	IsTokenized          string `json:"is_tokenized,optional"`
	WebLocale            string `json:"web_locale,optional"`
	CookieEnabled        string `json:"cookie_enabled,optional"`

	LastLoginDate        string `json:"last_login_date,optional"`
	LastRegistrationDate string `json:"last_registration_date,optional"`

	LastDenyEvent   string `json:"last_deny_event_date,optional"`
	LastReviewEvent string `json:"last_review_event_date,optional"`

	LastLimitsUpdateDate    string `json:"last_limit_update_date,optional"`
	LastLoanApplicationDate string `json:"last_loan_application_date,optional"`
	LastLoanApprovalDate    string `json:"last_loan_approval_date,optional"`
	PinUpdateDate           string `json:"pin_update_date,optional"`
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
