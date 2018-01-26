package eventsourcing

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/alexandervantrijffel/gonats/eventsourcing/contracts"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"github.com/rs/xid"

	"github.com/stretchr/testify/assert"
)

type UnitTestAggregate struct {
	AggregateCommonImpl
	RecordedTestMessageI int32
}

func (a *UnitTestAggregate) PerformTest(i int32) {
	evt := &eventsourcingcontracts.TestMessage{
		Id: xid.New().String(),
		I:  i,
	}
	a.Repository.ApplyEvent(a, evt)
}

func (a *UnitTestAggregate) HandleStateChange(event interface{}) {
	if casted, ok := event.(*eventsourcingcontracts.TestMessage); ok {
		fmt.Printf("We received a TestMessage with properties %+v\n", casted)
		a.RecordedTestMessageI = casted.I
	} else {
		fmt.Errorf("UnitTestAggregate Received unknown event of type %s\n", reflect.TypeOf(event))
	}
}

func TestGetStreamNameShouldReturnAggregateNameAndId(t *testing.T) {
	agg := &UnitTestAggregate{AggregateCommonImpl: AggregateCommonImpl{IdImpl: "my-name-is-billy"}}
	assert.Equal(t, "UnitTestAggregate|my-name-is-billy", GetStreamName(agg))
}

func TestLoadAggregate(t *testing.T) {
	agg := &UnitTestAggregate{AggregateCommonImpl: AggregateCommonImpl{IdImpl: "my-name-is-billy"}}
	repo, err := NewRepository("gonatseventsourcing_cluster",
		"test_client"+strconv.Itoa(random(0, 99999)),
		stan.NatsURL(stan.DefaultNatsURL))
	assert.Nil(t, err)
	repo.LoadAggregate(agg)
}

func TestSerializeEvent(t *testing.T) {
	evt := &eventsourcingcontracts.TestMessage{
		Id: xid.New().String(),
		I:  42,
	}
	eventEnvelope, err := serialize(evt)
	assert.Nil(t, err)
	eventEnvelopeBytes, _ := proto.Marshal(eventEnvelope)
	deserializedEventEnvelope, deserializedEvent, err := deserialize(eventEnvelopeBytes)
	assert.Nil(t, err)
	typ := reflect.TypeOf(deserializedEvent).Elem().Name()
	_ = typ
	typeOfEvent := reflect.TypeOf(deserializedEvent)
	_ = typeOfEvent
	deserializedMessage, ok := deserializedEvent.(*eventsourcingcontracts.TestMessage)
	assert.True(t, ok)
	_ = deserializedEventEnvelope
	assert.True(t, deserializedMessage.I == 42)
}

func TestApplyEventAndLoadAggregate(t *testing.T) {
	repo, err := NewRepository("gonatseventsourcing_cluster",
		"test_client"+strconv.Itoa(random(0, 99999)),
		stan.NatsURL(stan.DefaultNatsURL))
	assert.Nil(t, err)

	aggregateId := strconv.Itoa(random(0, 9999999))
	agg := &UnitTestAggregate{
		AggregateCommonImpl: AggregateCommonImpl{
			IdImpl:     aggregateId,
			Repository: repo,
		},
	}
	var i int32
	i = 51
	agg.PerformTest(i)

	var loadedAgg = &UnitTestAggregate{
		AggregateCommonImpl: AggregateCommonImpl{
			IdImpl: aggregateId,
		},
	}

	repo.LoadAggregate(loadedAgg)
	assert.Equal(t, i, loadedAgg.RecordedTestMessageI)
}
