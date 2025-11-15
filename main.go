package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	ganacheClient, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal("Ganache connect error: ", err)
	}

	var privateKeyHex string
	fmt.Print("приватный ключ: ") //length should be 32
	if scanner.Scan() {
		line := scanner.Text()
		if line[:2] == "0x" {
			privateKeyHex = line[2:]
		} else {
			privateKeyHex = line
		}
	}

	privateKey, err := crypto.HexToECDSA(string(privateKeyHex))
	if err != nil {
		fmt.Println(err)
	}

	publicKey := privateKey.PublicKey
	fromAddress := crypto.PubkeyToAddress(publicKey)

	var toAddressHex string
	fmt.Print("адрес получателя: ")
	if scanner.Scan() {
		toAddressHex = scanner.Text()
	}
	toAddress := common.HexToAddress(toAddressHex)

	// nonce
	ctx := context.Background()
	nonce, err := ganacheClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatal("nonce get err", err)
	}

	var value *big.Int
	fmt.Print("сумма транзакции(ETH): ")
	if scanner.Scan() {
		line := scanner.Text()
		ethAmount := new(big.Float)
		ethAmount.SetString(line)

		wei := new(big.Float).Mul(ethAmount, big.NewFloat(1e18))
		value, _ = wei.Int(nil)

	}

	gasLimit := uint64(21000) //standart gas price
	gasPrice, err := ganacheClient.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal("Ошибка получения gas price:", err)
	}

	transaction := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	//transaction := types.NewTx()

	chainID, err := ganacheClient.NetworkID(ctx)
	if err != nil {
		log.Fatal("Ошибка получения chain id: ", err)
	}

	signedTx, err := types.SignTx(transaction, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("Ошибка подписи транзакции", err)
	}

	err = ganacheClient.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal("Ошибка отправки транзакции: ", err)
	}

	fmt.Printf("Транзакция отправлена: ", signedTx.Hash())

	balance := getEthBalance(ganacheClient, ctx, fromAddress)
	balance1 := getEthBalance(ganacheClient, ctx, toAddress)

	fmt.Println(balance)
	fmt.Println(balance1)

}

func getEthBalance(client *ethclient.Client, ctx context.Context, addr common.Address) *big.Float {
	balance, err := client.BalanceAt(ctx, addr, nil)
	if err != nil {
		log.Fatal("Ошибка получения баланса: ", err)
	}
	etherBalance := new(big.Float).SetInt(balance)
	etherUnit := new(big.Float).SetFloat64(1e18)
	etherBalance = new(big.Float).Quo(etherBalance, etherUnit)

	return etherBalance
}
