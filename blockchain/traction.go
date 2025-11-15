package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TtackAddress(ctx context.Context, client *ethclient.Client, targetAddress common.Address) []*types.Transaction {

	header, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	last := header.Number().Uint64() // last block

	var match []*types.Transaction

	for i := uint64(0); i <= last; i++ {

		currentBlock, err := client.BlockByNumber(ctx, big.NewInt(int64(i)))
		if err != nil {
			log.Fatal("Ошибка извлечения блока по номеру: ", err)
		}

		for _, tx := range currentBlock.Transactions() {

			from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
			if err != nil {
				fmt.Print("Ошибка восстановления транзакции")
				continue
			}
			to := tx.To()

			if from == targetAddress || *to == targetAddress {
				match = append(match, tx)
			}
		}

	}

	return match

}

type TransactionData struct {
	Tx               string
	FromAddress      string
	ToAddress        string
	TransactionValue string
}

func ParseTransaction(addrs []*types.Transaction) []TransactionData {

	ParseResult := make([]TransactionData, len(addrs))

	for idx, tx := range addrs {

		signer := types.NewEIP155Signer(tx.ChainId())

		from, _ := types.Sender(signer, tx)

		to := func() string {
			if to := tx.To(); to != nil {
				return to.Hex()
			} else {
				return "[contract creation]"
			}
		}()

		etherBalance := new(big.Float).SetInt(tx.Value())
		etherUnit := new(big.Float).SetFloat64(1e18)
		etherBalance = new(big.Float).Quo(etherBalance, etherUnit)

		currentTransaction := TransactionData{
			Tx:               tx.Hash().Hex(),
			FromAddress:      from.Hex(),
			ToAddress:        to,
			TransactionValue: etherBalance.String(),
		}

		ParseResult[idx] = currentTransaction

	}

	return ParseResult

}
