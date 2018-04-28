package eventsourcing

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alexandervantrijffel/gonats/eventsourcing/contracts"
	"github.com/gogo/protobuf/proto"
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
		fmt.Printf("UnitTestAggregate Received unknown event of type %s\n", reflect.TypeOf(event))
	}
}

func TestSerializeEvent(t *testing.T) {
	evt := &eventsourcingcontracts.TestMessage{
		Id: xid.New().String(),
		I:  42,
	}
	eventEnvelope, err := serialize(evt)
	assert.Nil(t, err)
	eventEnvelopeBytes, _ := proto.Marshal(eventEnvelope)
	deserializedEventEnvelope, deserializedEvent, err := Deserialize(eventEnvelopeBytes)
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
