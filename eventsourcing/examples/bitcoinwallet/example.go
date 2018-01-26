package bitcoinwallet

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/alexandervantrijffel/gonats/eventsourcing"
	"github.com/alexandervantrijffel/gonats/eventsourcing/examples/bitcoinwallet/contracts"
	"github.com/rs/xid"
)

type BitcoinWallet struct {
	eventsourcing.AggregateCommonImpl
	address   string
	ownerName string
	Balance   bitcoinwalletcontracts.Amount
}

type Amount struct {
	Amount   float32
	Currency string
}

func (b *BitcoinWallet) HandleStateChange(event interface{}) {
	if castedEvent, ok := event.(*bitcoinwalletcontracts.BitcoinWalletCreated); ok {
		fmt.Printf("Received BitcoinWalletCreated %+v", castedEvent)
		b.address = castedEvent.Address
		b.ownerName = castedEvent.NameOfOwner
	} else if castedEvent, ok := event.(*bitcoinwalletcontracts.DepositMade); ok {
		b.Balance = bitcoinwalletcontracts.Amount{
			Amount:   b.Balance.Amount + castedEvent.Amount.Amount,
			Currency: b.Balance.Currency}
	} else if castedEvent, ok := event.(*bitcoinwalletcontracts.PaymentMade); ok {
		b.Balance = bitcoinwalletcontracts.Amount{
			Amount:   b.Balance.Amount - castedEvent.Amount.Amount,
			Currency: b.Balance.Currency}
	} else {
		fmt.Errorf("BitcoinWallet: could not handle unknown event of type %s\n", reflect.TypeOf(event))
	}
}

func (b *BitcoinWallet) Create(address string, ownerName string) error {
	evt := &bitcoinwalletcontracts.BitcoinWalletCreated{Address: address, NameOfOwner: ownerName}
	return b.Repository.ApplyEvent(b, evt)
}

func isPositiveBTCAmount(a bitcoinwalletcontracts.Amount) (errors []string, success bool) {
	if strings.ToLower(a.Currency) != "btc" {
		errors = append(errors, "Currency should be BTC")
	}
	if a.Amount <= 0 {
		errors = append(errors, "Amound should be above zero")
	}
	return errors, len(errors) == 0
}

func (b *BitcoinWallet) MakeDeposit(amount bitcoinwalletcontracts.Amount) error {
	if len(b.address) == 0 {
		return fmt.Errorf("Cannot make a deposit because the BitcoinWallet is not initialized with an address.")
	}
	if len(b.ownerName) == 0 {
		return fmt.Errorf("Cannot make a deposit because the owner of the wallet is not known.")
	}
	if errors, ok := isPositiveBTCAmount(amount); !ok {
		return fmt.Errorf("Cannot make a deposit because the amount is invalid. %s", strings.Join(errors, "\n"))
	}
	evt := &bitcoinwalletcontracts.DepositMade{Amount: &amount}
	return b.Repository.ApplyEvent(b, evt)
}

func (b *BitcoinWallet) MakePayment(destinationAddress string, amount bitcoinwalletcontracts.Amount) error {
	if errors, ok := isPositiveBTCAmount(amount); !ok {
		return fmt.Errorf("Cannot make a payment because the amount is invalid. %s", strings.Join(errors, "\n"))
	}
	if b.Balance.Amount < amount.Amount {
		return fmt.Errorf(
			"Cannot make a payment of %f %s to %s because the wallet has insufficient funds. Balance %f %s", amount.Amount, amount.Currency, destinationAddress, b.Balance.Amount, b.Balance.Currency)
	}
	evt := &bitcoinwalletcontracts.PaymentMade{DestinationAddress: destinationAddress, Amount: &amount}
	return b.Repository.ApplyEvent(b, evt)
}

func CreateWallet(repo eventsourcing.Repository, address string, ownerName string) (*BitcoinWallet, error) {
	aggregate := &BitcoinWallet{
		AggregateCommonImpl: eventsourcing.AggregateCommonImpl{
			IdImpl:     xid.New().String(),
			Repository: repo,
		},
	}
	err := aggregate.Create(address, ownerName)
	return aggregate, err
}
