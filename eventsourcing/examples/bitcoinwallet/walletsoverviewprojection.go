package bitcoinwallet

import (
	"fmt"
	"reflect"

	"github.com/alexandervantrijffel/gonats/eventsourcing"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/go-nats-streaming/pb"
)

func StartWalletsOverviewProjection(repo *eventsourcing.StreamRepository) {
	receivedMsg := make(chan struct{}, 1)
	subscription, err := repo.Connection.Subscribe("BitcoinWallet.>", func(m *stan.Msg) {
		eventEnvelope, event, err := eventsourcing.Deserialize(m.Data)
		_ = eventEnvelope
		if err == nil {
			fmt.Printf("WalletOverviewProjection: Received event of type " + reflect.TypeOf(event).Name())
		} else {
			fmt.Println("Failed to deserialize event of type " + reflect.TypeOf(event).Name())
		}
		receivedMsg <- struct{}{}
	}, stan.StartAt(pb.StartPosition_First))
	if err != nil {
		fmt.Printf("Failed to subscribe " + err.Error())
		return
	}
	defer subscription.Close()
	for {
		select {
		case <-receivedMsg:
		}
	}
}
