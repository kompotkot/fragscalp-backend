/*
Streamglass data extractors from Geth.
*/
package streamglass

import (
	"fmt"
	// "log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// type Entry struct {

// }

// func GethTxpool(chOut chan<- []streamglass.Tx) {
func GethTxpool(chOut chan<- []Txpool, minQuantity float64) {
	// Set connection with Ethereum blockchain via geth
	gethClient, err := rpc.Dial("http://127.0.0.1:18375")
	if err != nil {
		panic(fmt.Sprintf("Could not connect to geth: %s", err.Error()))
	}
	defer gethClient.Close()

	currentTransactions := make(map[common.Hash]bool)

	// Structure of the map:
	// {pending | queued} -> "from" address -> nonce -> Transaction
	var result map[string]map[string]map[uint64]*TxGeth

	for {
		// fmt.Println("Checking pending transactions in node:")
		gethClient.Call(&result, "txpool_content")
		pendingTransactions := result["pending"]

		// Mark all transactions from previous iteration as false
		cacheSize := 0
		for transactionHash := range currentTransactions {
			currentTransactions[transactionHash] = false
			cacheSize++
		}
		// fmt.Printf("\tSize of pending transactions cache at the beginning: %d\n", cacheSize)

		txs := []Txpool{}

		// Iterate over fetched result
		addedTransactionsCounter := 0
		for _, transactionsByNonce := range pendingTransactions {
			for _, transaction := range transactionsByNonce {
				transactionHash := transaction.Hash
				_, transactionProcessed := currentTransactions[transactionHash]
				if !transactionProcessed {
					quantityBid := new(big.Float).Quo(new(big.Float).SetInt(transaction.Value.ToInt()), big.NewFloat(params.Ether))
					quantity, _ := quantityBid.Float64()
					if quantity > minQuantity {
						txs = append(txs, Txpool{
							Hash:     transactionHash.String(),
							Quantity: quantity,
						})
						addedTransactionsCounter++
					}

				}
				currentTransactions[transactionHash] = true
			}
		}
		chOut <- txs
		// fmt.Printf("\tPublishing %d transactions\n", addedTransactionsCounter)

		// Clean the slice out of disappeared transactions
		droppedTransactionsCounter := 0
		for transactionHash, justProcessed := range currentTransactions {
			if !justProcessed {
				delete(currentTransactions, transactionHash)
				droppedTransactionsCounter++
			}
		}
		// fmt.Printf("\tDropped transactions: %d\n", droppedTransactionsCounter)

		// fmt.Printf("Sleeping for %d seconds\n", 1)
		time.Sleep(time.Duration(1) * time.Second)

	}
}
