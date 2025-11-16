package neo4j

import (
	"context"
	"eth-graph-explorer/blockchain"
	"fmt"
	"log"

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

	fmt.Println("instance connect")

	return driver, ctx

}

func MergeAddress(addr1, addr2 string, session neo4j.SessionWithContext, driver neo4j.DriverWithContext) {

	/// переписать запрос

	ctx := context.Background()

	query := `
		MERGE (:Block{address: $address1})
		MERGE (:Block{address: $address2})
	`

	data := map[string]any{
		"address1": addr1,
		"address2": addr2,
	}

	result, err := neo4j.ExecuteQuery(ctx, driver, query, data, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("eth"))
	if err != nil {
		log.Fatal("не удалось создать адреса")
	}

	summary := result.Summary
	fmt.Printf("Created %v nodes in %+v.\n",
		summary.Counters().NodesCreated(),
		summary.ResultAvailableAfter())

}

func checkRelationExist(ctx context.Context, addr1, addr2 string, driver neo4j.DriverWithContext) bool {

	query := `
		RETURN EXISTS {
  			MATCH (a1:Block {address: $addr1})
  			MATCH (a2:Block {address: $addr2})
  			MATCH p=(a1)-[r:TRANSFER]->(a2) RETURN p
		} AS found
	`

	queryData := map[string]any{
		"addr1": addr1,
		"addr2": addr2,
	}

	result, err := neo4j.ExecuteQuery(ctx, driver, query, queryData,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("eth"))

	if err != nil {
		log.Fatal("не удалось выполнить запрос checkRelationExist", err)
	}

	if len(result.Records) == 0 {
		// пустой ответ
		return false
	}

	foundValue, ok := result.Records[0].Get("found")
	if !ok {
		// поле found не найдено
		return false
	}

	foundBool, ok := foundValue.(bool)
	if !ok {
		// неверный тип поля
		return false
	}

	return foundBool

}

func checkAddressExist(session neo4j.SessionWithContext, data blockchain.TransactionData, driver neo4j.DriverWithContext) bool {

	ctx := context.Background()

	readQuety :=
		`RETURN EXISTS {
    		MATCH (:Block {address: $fromAddress})
    		MATCH (:Block {address: $toAddress})
		} AS found`

	queryData := map[string]any{
		"fromAddress": data.FromAddress,
		"toAddress":   data.ToAddress,
	}

	//переписать через ExecuteQuery

	read, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		resords, err := tx.Run(ctx, readQuety, queryData)
		if err != nil {
			fmt.Println("ошибка тут")
			return false, err
		}
		if resords.Next(ctx) {
			value, found := resords.Record().Get("found")
			if !found {
				return false, fmt.Errorf("поле 'found' не найдено")
			}
			if result, ok := value.(bool); ok {
				return result, nil
			}
			return false, fmt.Errorf("ожидался bool, получен %T", value)
		}

		return false, fmt.Errorf("пустой результат")
	})
	if err != nil {
		log.Fatal("Ошибка выполнения запроса", err)
	}

	result := read.(bool)
	if result {
		return true
	} else {
		MergeAddress(data.FromAddress, data.ToAddress, session, driver)
		return true
	}

}

func makeRelation(addr1, addr2, value, tx string, driver neo4j.DriverWithContext) error {

	ctx := context.Background()

	relationQuety := `
		MATCH (a1:Block{address: $fromAddress})
		MATCH (a2:Block{address: $toAddress})
		CREATE (a1)-[:TRANSFER {val: $value, transaction: $tx}]->(a2)
	`

	queryData := map[string]any{
		"fromAddress": addr1,
		"toAddress":   addr2,
		"value":       value,
		"tx":          tx,
	}

	_, err := neo4j.ExecuteQuery(ctx, driver, relationQuety, queryData, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("eth"))
	if err != nil {
		return fmt.Errorf("ошибка создания связи")
	}

	return nil

}

func Vizualization(data *[]blockchain.TransactionData, driver neo4j.DriverWithContext) {

	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "eth"})

	for _, elem := range *data {
		check := checkAddressExist(session, elem, driver)
		if check {
			fmt.Println("exist")
		}
		// переписать возврат checkAddressExist

		exist := checkRelationExist(ctx, elem.FromAddress, elem.ToAddress, driver)
		if exist {
			continue
		} else {
			makeRelation(elem.FromAddress, elem.ToAddress, elem.TransactionValue, elem.Tx, driver)
		}
	}

	fmt.Println("Визуализация завершена")

}
