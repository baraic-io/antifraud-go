```
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

	af_transaction := af.AF_Transaction{
		AF_Id:      uuid.NewString(),
		AF_AddDate: now.Format(time.RFC3339Nano),
		Transaction: af.Transaction{
			Id:                 uuid.New().String(),
			Type:               "deposit",
			Date:               now.Format(time.RFC3339Nano),
			Amount:             "100000",
			Currency:           "KZT",
			ClientId:           uuid.New().String(),
			ClientName:         "John Smith",
			ClientPAN:          "111111******1111",
			ClientCVV:          "111",
			ClientCardHolder:   "JOHN SMITH",
			ClientPhone:        "+77007007070",
			MerchantTerminalId: "00000001",
			Channel:            "E-com",
			LocationIp:         "192.168.0.1",
		},
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

	/* Step 8: Finalize transaction validation */
	finalResolution, err := client.FinalizeTransaction(af_transaction)
	if err != nil {
		panic("[STEP 8] failed: " + err.Error())
	}

	fmt.Println("[STEP 8] success")
	fmt.Printf("Final resolution: %+v\n", finalResolution)

	/* Step 9: Store final resolution */
	if err := client.StoreFinalResolution(finalResolution); err != nil {
		panic("[STEP 9] failed: " + err.Error())
	}

	fmt.Println("[STEP 9] success")
}
```