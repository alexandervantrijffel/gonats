package eventsourcing

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
)

type Repository interface {
	LoadAggregate(AggregateCommon) (interface{}, error)
	ApplyEvent(AggregateCommon, proto.Message) error
}

var MainRepository StreamRepository

func NewRepository(stanClusterID, clientID string, options ...stan.Option) (Repository, error) {
	var err error
	if (StreamRepository{}) == MainRepository {
		conn, err := stan.Connect(stanClusterID, clientID, options...)
		if checkErr(err, "Failed to connect to NATS Streaming") {
			return nil, err
		}

		// todo close conn at shutdown
		MainRepository, err = StreamRepository{Connection: conn}, nil
	}
	return &MainRepository, err
}

type StreamRepository struct {
	Connection stan.Conn
}

func (sr *StreamRepository) LoadAggregate(aggregate AggregateCommon) (interface{}, error) {
	sr.ReadAllStreamItems(aggregate, func(event interface{}) {
		fmt.Printf("LoadAggregate: Received event %+v\n", event)
		aggregate.HandleStateChange(event)
	})
	return nil, nil
}

func (sr *StreamRepository) ApplyEvent(aggregate AggregateCommon, event proto.Message) (err error) {
	if eventEnvelope, err := serialize(event); err == nil {
		if marshalledEnvelope, err := proto.Marshal(eventEnvelope); err == nil {
			streamName := GetStreamName(aggregate)
			return sr.Connection.Publish(streamName, marshalledEnvelope)
		}
	}
	return err
}

func (sr *StreamRepository) ReadAllStreamItems(aggregate AggregateCommon, eventHandler func(interface{})) error {
	receivedMsg := make(chan struct{}, 1)
	subStreaming, err := sr.Connection.Subscribe(GetStreamName(aggregate), func(m *stan.Msg) {
		fmt.Printf("ReadAllStreamItems: Received event with sequence %d for aggregate %T id %s\n",
			m.Sequence, aggregate, aggregate.ID())
		if len(m.Data) == 0 {
			fmt.Println("ReadAllStreamItems: Skipping processing of message because the length is 0")
		} else {
			eventEnvelope, event, err := deserialize(m.Data)
			_ = eventEnvelope
			checkErr(err, "ReadAllStreamItems: Deserialize incoming message")
			eventHandler(event)
			receivedMsg <- struct{}{}
		}
	}, stan.DeliverAllAvailable())
	if checkErr(err, "ReadAllStreamItems: NATS streaming subscribe"); err != nil {
		return err
	}
	for {
		fmt.Println("ReadAllStreamItems: Waiting for ONE message")
		select {
		case <-receivedMsg:
			fmt.Printf("ReadAllStreamItems: received a message from the chan\n")
		case <-time.After(200 * time.Millisecond):
			fmt.Println("ReadAllStreamItems: timeout after 200ms")
		}
		pendingMessages, _, err := subStreaming.Pending()
		checkErr(err, "ReadAllStreamItems: Pending()")
		fmt.Printf("ReadAllStreamItems pending events to process: %d\n", pendingMessages)
		if pendingMessages == 0 {
			subStreaming.Close()
			fmt.Println("ReadAllStreamItems: no more items left")
			return nil
		}
	}
}
