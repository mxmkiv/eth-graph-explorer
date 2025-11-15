package blockchain

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func BlockChainConnect() (*ethclient.Client, string) {

	ganacheClient, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal("Ошибка подключения к сети:", err)
	}

	msg := "успешное подключение"

	return ganacheClient, msg
}
