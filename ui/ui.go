package ui

import (
	"bufio"
	"eth-graph-explorer/blockchain"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func ShowTransaction(data []blockchain.TransactionData) {

	for _, elem := range data {
		fmt.Println("Transaction: ", elem.Tx)
		fmt.Println("From: ", elem.FromAddress)
		fmt.Println("To: ", elem.ToAddress)
		fmt.Println("Value: ", elem.TransactionValue)
	}
}

func VisualisationAccept(scanner *bufio.Scanner) bool {

	fmt.Println("Визуализировать данные?: (y)/n")
Acception:
	for {
		scanner.Scan()
		line := scanner.Text()
		switch line {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Println("Некорректный выбор")
			continue Acception
		}
	}

}

func GetTargetAddress(scanner *bufio.Scanner) common.Address {
	fmt.Print("адрес: ")

	var addr common.Address
	if scanner.Scan() {
		line := scanner.Text()
		addr = common.HexToAddress(line)
	}

	return addr

}
