package bitcoinwallet

import (
	"fmt"
	"reflect"

	"github.com/alexandervantrijffel/gonats/eventsourcing"
	"github.com/alexandervantrijffel/gonats/eventsourcing/examples/bitcoinwallet/contracts"
	"github.com/rs/xid"
)

type BitcoinWallet struct {
	eventsourcing.AggregateCommonImpl
	address   string
	ownerName string
}

func (b *BitcoinWallet) HandleStateChange(event interface{}) {
	if castedEvent, ok := event.(*bitcoinwalletcontracts.BitcoinWalletCreated); ok {
		fmt.Printf("Received BitcoinWalletCreated %+v", castedEvent)
		b.address = castedEvent.Address
		b.ownerName = castedEvent.NameOfOwner
	} else {
		fmt.Errorf("BitcoinWallet: could not handle unknown event of type %s\n", reflect.TypeOf(event))
	}
}

func (b *BitcoinWallet) Create(address string, ownerName string) error {
	evt := &bitcoinwalletcontracts.BitcoinWalletCreated{Address: address, NameOfOwner: ownerName}
	return b.Repository.ApplyEvent(b, evt)
}

func CreateWallet(address string, ownerName string) (*BitcoinWallet, error) {
	repo, err := eventsourcing.NewRepository("gonatseventsourcing_cluster", "test_client1")
	if err != nil {
		return nil, err
	}
	aggregate := &BitcoinWallet{
		AggregateCommonImpl: eventsourcing.AggregateCommonImpl{
			IdImpl:     xid.New().String(),
			Repository: repo,
		},
	}
	err = aggregate.Create(address, ownerName)
	return aggregate, err
}
