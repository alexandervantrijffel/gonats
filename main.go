package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alexandervantrijffel/gonats/eventsourcing"
	"github.com/alexandervantrijffel/gonats/eventsourcing/examples/bitcoinwallet/contracts"

	"github.com/alexandervantrijffel/gonats/eventsourcing/examples/bitcoinwallet"
)

func main() {
	repo, err := eventsourcing.NewRepository("gonatseventsourcing_cluster", "test_client2")
	if err != nil {
		fmt.Printf("Failed to connect to NATS %+v", err)
		return
	}
	go bitcoinwallet.StartWalletsOverviewProjection(repo)
	wallet, err := bitcoinwallet.CreateWallet(repo, "aajb21", "Alexander")
	if err != nil {
		fmt.Println("Failed to create a new bitcoin wallet.", err.Error())
		return
	}
	wallet.MakeDeposit(bitcoinwalletcontracts.Amount{Amount: 0.0027, Currency: "BTC"})
	wallet.MakeDeposit(bitcoinwalletcontracts.Amount{Amount: 0.0013, Currency: "BTC"})
	wallet.MakePayment("eed7ab", bitcoinwalletcontracts.Amount{Amount: 0.004, Currency: "BTC"})
	// will fail because of insufficient funds
	err = wallet.MakePayment("jjbq24", bitcoinwalletcontracts.Amount{Amount: 0.0001, Currency: "BTC"})
	fmt.Println("Result from last payment", err)
	fmt.Print("Press a key to quit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}
