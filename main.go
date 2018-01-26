package main

import (
	"fmt"

	"github.com/alexandervantrijffel/gonats/eventsourcing/examples/bitcoinwallet/contracts"

	"github.com/alexandervantrijffel/gonats/eventsourcing/examples/bitcoinwallet"
)

func main() {
	wallet, err := bitcoinwallet.CreateWallet("aajb21", "Alexander")
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
}
