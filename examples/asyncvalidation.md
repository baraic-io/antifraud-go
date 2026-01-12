```go
package examples

import (
	"fmt"
	"os"

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

	/* TODO: Add all required fields */
	result, err := client.ValidateTransactionAsync(af.Transaction{
		Id:               uuid.New().String(),
		SourceUserId:     uuid.New().String(),
		SourceIdentifier: "000000000001",
		SourceFullname:   "John Smith",
		TargetUserId:     uuid.New().String(),
		TargetIdentifier: "000000000002",
		TargetFullname:   "Mike Hart",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("AF_ID: %s", result.AF_Id)
}
```
