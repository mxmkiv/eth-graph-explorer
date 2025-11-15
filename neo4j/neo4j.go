package neo4j

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateConnection() (neo4j.DriverWithContext, context.Context) {
	ctx := context.Background()
	dbUri := "neo4j://127.0.0.1:7687"
	dbUser := "neo4j"
	dbPassword := "testdatabase13"

	driver, _ := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))

	err := driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)

	fmt.Println("instance connect")

	return driver, ctx

}

// func Vizualization(data []blockchain.TransactionData, driver neo4j.DriverWithContext, ctx context.Context) {

// 	addressQuery := `RETURN EXIST { MATCH (:Block {address: $addr})}`

// 	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "eth"})

// 	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

// 	}

// }
