package bitcoinwallet

import (
	"fmt"
	"reflect"

	"github.com/alexandervantrijffel/gonats/eventsourcing"
	"github.com/nats-io/go-nats-streaming"
)

func StartWalletsOverviewProjection(repo *eventsourcing.StreamRepository) {
	repo.Connection.Subscribe("BitcoinWallet.b9lluamvik2kojvb8hl0", func(m *stan.Msg) {
		eventEnvelope, event, err := eventsourcing.Deserialize(m.Data)
		_ = eventEnvelope
		if err == nil {
			fmt.Printf("WalletOverviewProjection: Received event of type %s", reflect.TypeOf(event))
		} else {
			fmt.Println("Failed to deserialize event of type %s", reflect.TypeOf(event))
		}
	})

	//, stan.StartAt(pb.StartPosition_First)
}
