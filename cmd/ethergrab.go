package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/xllwhoami/ethergrab/internal/database"
	"github.com/xllwhoami/ethergrab/pkg/etherclient"
	"golang.org/x/sync/semaphore"
)

func processBlock(blockno string, client *etherclient.Client, db *database.Database, semaphore *semaphore.Weighted) {
	defer semaphore.Release(int64(1))

	block, err := client.GetBlockByNumber(blockno)

	if err != nil {
		fmt.Println("Can`t get block info")
		panic(err)
	}

	if len(block.Result.Transactions) < 1 {
		fmt.Printf("Block %d has no transactions\n", client.NumberFromHex(blockno))
	}

	for _, t := range block.Result.Transactions {
		from := t.From
		to := t.To

		db.InsertAddress(from)
		db.InsertAddress(to)

		fmt.Println("Transaction " + t.Hash + " processed")
	}
}

func main() {
	fmt.Println("Starting...")

	err := godotenv.Load()

	if err != nil {
		fmt.Println("Can`t load dotenv")
		panic(err)
	}

	endpoint := os.Getenv("RPCNODE_ENDPOINT")

	db, _ := database.NewDatabase("addresses.db")

	db.CreateTables()

	client := etherclient.NewClient(endpoint)

	latest_blockno, err := client.GetLatestBlockNumber()

	if err != nil {
		fmt.Println("Can`t get block number")
		panic(err)
	}

	context := context.Background()

	semaphore := semaphore.NewWeighted(int64(1))

	for i := 0; i < int(client.NumberFromHex(latest_blockno)); i++ {
		if err := semaphore.Acquire(context, 1); err != nil {
			fmt.Println("Semaphore can`t Acquire")
			panic(err)
		}

		go processBlock(client.NumberToHex(int64(i)), client, db, semaphore)
	}

}
