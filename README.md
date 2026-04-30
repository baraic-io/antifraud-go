# Antifraud Go SDK

This is the Go SDK for the Antifraud Service.

## Installation

```bash
go get github.com/baraic-io/antifraud-go
```

## Methods

### NewClient

Creates a new Antifraud Client.

```go
func NewClient(c ClientConfig) (Client, error)
```

### StoreTransaction

Stores a transaction in the antifraud system.

```go
func (c Client) StoreTransaction(af_transaction AF_Transaction) error
```

### ValidateTransactionByAML

Validates a transaction using the AML service.

```go
func (c Client) ValidateTransactionByAML(af_transaction AF_Transaction) (ServiceResolution, error)
```

### ValidateTransactionByFC

Validates a transaction using the FC service.

```go
func (c Client) ValidateTransactionByFC(af_transaction AF_Transaction) (ServiceResolution, error)
```

### ValidateTransactionByML

Validates a transaction using the ML service.

```go
func (c Client) ValidateTransactionByML(af_transaction AF_Transaction) (ServiceResolution, error)
```

### StoreServiceResolution

Stores the resolution from a service check (AML, FC, LST).

```go
func (c Client) StoreServiceResolution(resolution ServiceResolution) error
```

### AddTransactionServiceCheck

Adds a completed service check resolution to the transaction aggregation process.

```go
func (c Client) AddTransactionServiceCheck(resolution ServiceResolution) error
```

### FinalizeTransaction

Finalizes the transaction validation process and retrieves the final resolution.

```go
func (c Client) FinalizeTransaction(af_transaction AF_Transaction) (FinalResolution, error)
```

### StoreFinalResolution

Stores the final resolution of the transaction.

```go
func (c Client) StoreFinalResolution(resolution FinalResolution) error
```

## Usage

### Transaction Validation Flow

The following example demonstrates a full transaction validation flow, including storing the transaction, validating it against AML and FC services, and finalizing the result.

```go
package main

import (
	"fmt"
	"os"
	"time"

	af "github.com/baraic-io/antifraud-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	apiKey := os.Getenv("API_KEY")
	host := os.Getenv("API_HOST")

	client, err := af.NewClient(af.ClientConfig{Host: host, APIKey: apiKey})
	if err != nil {
		panic(err)
	}

	now := time.Now()

	sourceTransaction := map[string]interface{}{
		"ip_addr": "192.168.0.1",
		"req_id":  uuid.NewString(),
		"finoper": map[string]interface{}{
			"product_id":         "C2C2Out",
			"oper_id":            uuid.NewString(),
			"reason":             "Test transfer",
			"knp_code":           "111",
			"oper_type":          float64(1),
			"oper_date_time":     now.Format(time.RFC3339Nano),
			"creation_date_time": now.Format(time.RFC3339Nano),
			"person": []interface{}{
				map[string]interface{}{
					"iin":             "123456789012",
					"full_name":       "Test Sender",
					"person_type":     float64(2),
					"mobile_number":   "+77007007070",
					"client_reg_date": now.Format(time.RFC3339Nano),
					"is_client":       true,
				},
				map[string]interface{}{
					"iin":             "987654321098",
					"full_name":       "Test Recipient",
					"person_type":     float64(2),
					"mobile_number":   "+77007007071",
					"client_reg_date": now.Format(time.RFC3339Nano),
					"is_client":       true,
				},
			},
			"finoper_dc": []interface{}{
				map[string]interface{}{
					"amount_kzt":     float64(100000.0),
					"bank_bic":       "AAAABBBB",
					"bank_name":      "Test Bank",
					"card_number":    "111111******1111",
					"card_exp_date":  "12/26",
					"card_open_date": "01/24",
					"account_number": "KZ000000",
				},
				map[string]interface{}{
					"amount_kzt":     float64(100000.0),
					"bank_bic":       "CCCCDDDD",
					"bank_name":      "Another Bank",
					"card_number":    "222222******2222",
					"card_exp_date":  "11/25",
					"card_open_date": "01/23",
					"account_number": "KZ111111",
				},
			},
		},
	}

	transaction, err := client.ToAFTransaction(sourceTransaction)
	if err != nil {
		panic(err)
	}

	af_transaction := af.AF_Transaction{
		AF_Id:       uuid.NewString(),
		AF_AddDate:  now.Format(time.RFC3339Nano),
		Transaction: transaction,
	}

	/* Step 1: Store transaction */
	if err := client.StoreTransaction(af_transaction); err != nil {
		panic("[STEP 1] failed: " + err.Error())
	}

	fmt.Println("[STEP 1] success")

	/* Step 2: Transaction validation by AML service */
	amlresult, err := client.ValidateTransactionByAML(af_transaction)
	if err != nil {
		panic("[STEP 2] failed: " + err.Error())
	}

	fmt.Println("[STEP 2] success")
	fmt.Printf("Validation transaction by AML service: %+v\n", amlresult)

	/* Step 3: Store AML service resolution */
	if err := client.StoreServiceResolution(amlresult); err != nil {
		panic("[STEP 3] failed: " + err.Error())
	}

	fmt.Println("[STEP 3] success")

	/* Step 4: Add AML service resolution to finalize process */
	if err := client.AddTransactionServiceCheck(amlresult); err != nil {
		panic("[STEP 4] failed: " + err.Error())
	}

	fmt.Println("[STEP 4] success")

	/* Step 5: Transaction validation by FC service */
	fcresult, err := client.ValidateTransactionByFC(af_transaction)
	if err != nil {
		panic("[STEP 5] failed: " + err.Error())
	}

	fmt.Println("[STEP 5] success")
	fmt.Printf("Validation transaction by FC service: %+v\n", fcresult)

	/* Step 6: Store FC service resolution */
	if err := client.StoreServiceResolution(fcresult); err != nil {
		panic("[STEP 6] failed: " + err.Error())
	}

	fmt.Println("[STEP 6] success")

	/* Step 7: Add FC service resolution to finalize process */
	if err := client.AddTransactionServiceCheck(fcresult); err != nil {
		panic("[STEP 7] failed: " + err.Error())
	}

	fmt.Println("[STEP 7] success")

	/* Step 8: Transaction validation by ML service */
	mlresult, err := client.ValidateTransactionByML(af_transaction)
	if err != nil {
		panic("[STEP 8] failed: " + err.Error())
	}

	fmt.Println("[STEP 8] success")
	fmt.Printf("Validation transaction by ML service: %+v\n", mlresult)

	/* Step 9: Store ML service resolution */
	if err := client.StoreServiceResolution(mlresult); err != nil {
		panic("[STEP 9] failed: " + err.Error())
	}

	fmt.Println("[STEP 9] success")

	/* Step 10: Add ML service resolution to finalize process */
	if err := client.AddTransactionServiceCheck(mlresult); err != nil {
		panic("[STEP 10] failed: " + err.Error())
	}

	fmt.Println("[STEP 10] success")

	/* Step 11: Finalize transaction validation */
	finalResolution, err := client.FinalizeTransaction(af_transaction)
	if err != nil {
		panic("[STEP 11] failed: " + err.Error())
	}

	fmt.Println("[STEP 11] success")
	fmt.Printf("Final resolution: %+v\n", finalResolution)

	/* Step 12: Store final resolution */
	if err := client.StoreFinalResolution(finalResolution); err != nil {
		panic("[STEP 12] failed: " + err.Error())
	}

	fmt.Println("[STEP 12] success")
}
```
