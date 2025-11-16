package blockchain

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func BlockChainConnect() (*ethclient.Client, string) {

	ganacheClient, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal("ошибка подключения к сети:", err)
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

	// _, err = ganacheClient.BlockNumber(ctx)
	// if err != nil {
	// 	log.Fatal("узел не отвечает", err)
	// }

	msg := "успешное подключение"

	return ganacheClient, msg
}
