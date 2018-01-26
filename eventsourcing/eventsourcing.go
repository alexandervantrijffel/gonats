package eventsourcing

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/alexandervantrijffel/gonats/eventsourcing/contracts"
	proto "github.com/golang/protobuf/proto"
)

type AggregateCommon interface {
	ID() string
	HandleStateChange(interface{})
}

type AggregateCommonImpl struct {
	IdImpl     string
	Repository Repository
}

func (a *AggregateCommonImpl) ID() string {
	return a.IdImpl
}

func serialize(event interface{}) (*eventsourcingcontracts.EventEnvelope, error) {
	eventBytes, err := proto.Marshal(event.(proto.Message))
	if err != nil {
		return nil, err
	}
	return &eventsourcingcontracts.EventEnvelope{
		TypeName:    getTypeName(event),
		MessageData: eventBytes,
	}, nil
}
func getTypeName(obj interface{}) string {
	return strings.Replace(fmt.Sprintf("%T", obj), "*", "", -1)
}
func deserialize(eventEnvelopeBytes []byte) (eventEnvelope *eventsourcingcontracts.EventEnvelope, event interface{}, err error) {
	envelope := &eventsourcingcontracts.EventEnvelope{}
	proto.Unmarshal(eventEnvelopeBytes, envelope)
	eventEnvelope = envelope

	t := proto.MessageType(eventEnvelope.TypeName)
	e := reflect.New(t.Elem())
	e2, isProtoMessage := e.Interface().(proto.Message)
	if !isProtoMessage {
		err = fmt.Errorf(
			"Event with EventEnvelope.TypeName %s and type %T can not be casted to type proto.Message",
			eventEnvelope.TypeName, e)
		return
	}
	proto.Unmarshal(eventEnvelope.MessageData, e2)
	event = e2
	return
}

func GetStreamName(aggregate AggregateCommon) string {
	splittedTypeNameWithoutAggregate := strings.Split(getTypeName(aggregate), ".")
	typeNameWithoutPackage := splittedTypeNameWithoutAggregate[len(splittedTypeNameWithoutAggregate)-1]
	return fmt.Sprintf("%s|%s", typeNameWithoutPackage, aggregate.ID())
}

func checkErr(err error, desc string) bool {
	if err != nil {
		log.Println(desc+" FAILED, reason ", err)
	}
	return err != nil
}
