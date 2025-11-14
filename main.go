package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	// ganacheClient, err := ethclient.Dial("http://127.0.0.1:8545")
	// if err != nil {
	// 	log.Fatal("RPC error: ", err)
	// }

	//var privateKey *ecdsa.PrivateKey

	var privateKeyHex string
	fmt.Print("приватный ключ: ")
	if scanner.Scan() {
		line := scanner.Text()
		if line[:2] == "0x" {
			privateKeyHex = line[2:]
		} else {
			privateKeyHex = line
		}
	}

	var toAddress string
	fmt.Print("адрес получателя: ")
	if scanner.Scan() {
		toAddress = scanner.Text()
	}

	fmt.Println(privateKeyHex)
	fmt.Println(toAddress)

}
