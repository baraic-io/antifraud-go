```go
package main

import (
	"fmt"
	"os"
	"time"

	af "github.com/baraic-io/antifraud-go"
	"github.com/google/uuid"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	host := os.Getenv("API_HOST")

	client, err := af.NewClient(af.ClientConfig{Host: host, APIKey: apiKey})
	if err != nil {
		panic(err)
	}

	now := time.Now()

	result, err := client.ValidateTransactionByAML(af.AF_Transaction{
		Transaction: af.Transaction{
			Id:                 uuid.New().String(),
			SourceUserId:       uuid.NewString(),
			SourceIdentifier:   "000000000001",
			SourceFullname:     "John Smith",
			SourceCardNumber:   "3300000000000000",
			SourceAccount:      "KZ00100000000000000",
			TargetUserId:       uuid.NewString(),
			TargetIdentifier:   "000000000002",
			TargetFullname:     "Ken Doroty",
			TargetCardNumber:   "",
			TargetAccount:      "",
			MerchantId:         uuid.NewString(),
			MerchantTerminalId: uuid.NewString(),
			MerchantMCCCode:    "5146",
			Date:               now.Format(time.RFC3339Nano),
			Time:               now.Format("15:04:05"),
			Amount:             "50000",
			Currency:           "KZT",
			PaymentMode:        "Online",
			TransactionType:    "P2P",
			TransactionCountry: "",
			TransactionCity:    "Almaty",
			TransactionChannel: "Mobile",
			TransactionRRN:     uuid.NewString(),
			TransactionStatus:  "3DS",
			CardType:           "debit",
			DeviceId:           uuid.NewString(),
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("AML Validation Resolution: %v", result)
}
```