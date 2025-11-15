package main

import (
	"bufio"
	"context"
	"eth-graph-explorer/blockchain"
	"eth-graph-explorer/transaction"
	"eth-graph-explorer/ui"
	"fmt"
	"os"
	"strconv"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()
	client, msg := blockchain.BlockChainConnect()

Menu:
	for {
		fmt.Println("eth-graph-explorer")
		fmt.Println("1. Создание транзакций")
		fmt.Println("2. Трекер транзакций")
		fmt.Println("4. Выход")
		choose := getChoose(scanner)
		switch choose {
		case 1:
			fmt.Println(msg)
			transaction.CreateTransaction(client)
		case 2:
			targetAddress := ui.GetTargetAddress(scanner)
			trackResult := blockchain.TtackAddress(ctx, client, targetAddress)
			structResult := blockchain.ParseTransaction(trackResult)
			ui.ShowTransaction(structResult)
			//answer := ui.VisualisationAccept(scanner)
			// if answer {
			// 	driver, ctx := neo4j.CreateConnection()
			// 	//neo4j.Vizualization(structResult, driver, ctx)
			// }
		case 4:
			break Menu
		default:
			continue Menu
		}
	}

}

func getChoose(scanner *bufio.Scanner) int {

	var numb int
	if scanner.Scan() {
		numb, _ = strconv.Atoi(scanner.Text())
	}
	return numb
}
