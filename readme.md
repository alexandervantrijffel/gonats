CQRS Event Sourcing with NATS Streaming

Not only is Go an excellent choice for writing system applications and network services, such as the container platform Docker, the container orchestrator Kubernetes and the network proxy Traefik, it is also very suitable to build business oriented applications. This post describes how to implement the business logic of an application based on CQRS Event Sourcing. 

If you are new to CQRS Event Sourcing and are ready for a deep dive, watch this [6.5 hours CQRS workshop](https://www.youtube.com/watch?v=whCk1Q87_ZI) by the inventor of CQRS Event Sourcing, Greg Young.

gonats/eventsourcing is a library for applying the event sourcing pattern in Golang using durable event streams from NATS streaming. Load aggregate roots from event streams stored in NATS streaming. Persist new events from the aggregates which can be used to build read models. 

The library is inspired by the GetEventStore.com client API. 

# How to create your first aggregate with gonats/eventsourcing
Define the struct for your aggregate and embed the eventsourcing.AggregateCommonImpl struct.
```
type BitcoinWallet struct {
	eventsourcing.AggregateCommonImpl
	address   string
	ownerName string
}
```
Write the first event that will be published by the BitcoinWallet aggregate. Events are defined as Protocol Buffers types. Add a new file with the name bitcoinwalletcontracts.proto with the following contents:
```
syntax = "proto3";
package bitcoinwalletcontracts;

message BitcoinWalletCreated {
    string address = 1;
    string nameOfOwner = 2;
}
```
The .proto file can be compiled to a go file with go types with the command `protoc bitcoinwalletcontracts.proto --go_out=plugins=grpc:.`  

Implement the HandleStateChange method to satisfy the eventsourcing.AggregateCommon interface
```
func (b *BitcoinWallet) HandleStateChange(event interface{}) {
	if castedEvent, ok := event.(*bitcoinwalletcontracts.BitcoinWalletCreated); ok {
		fmt.Printf("Received BitcoinWalletCreated %+v", castedEvent)
		b.address = castedEvent.Address
		b.ownerName = castedEvent.NameOfOwner
	} else {
		fmt.Errorf("BitcoinWallet: could not handle unknown event of type %s\n", reflect.TypeOf(event))
	}
}
```
Create a new instance of the BitCoinWallet aggregate 
```
myBitcoinWallet := &BitcoinWallet{
		AggregateCommonImpl: eventsourcing.AggregateCommonImpl{
			IdImpl:     "1",
			Repository: repo,
		},
	}
```
Instantiate a NATS Streaming reposistory with the following line.
```
repo, err := eventsourcing.NewRepository("gonatseventsourcing_cluster", "test_client1")
```
Using the repository, the BitcoinWalletCreated event is published as follows.
```
evt := &bitcoinwalletcontracts.BitcoinWalletCreated{Address: address, NameOfOwner: ownerName}
b.Repository.ApplyEvent(myBitcoinWallet, evt)
```
The BitcoinWallet event is added to the stream of this instance of the BitcoinWallet aggregate.  

The aggregate state can later be retrieved with the following line.
```
var loadedAgg = &UnitTestAggregate{
  	AggregateCommonImpl: AggregateCommonImpl{
	  IdImpl: "1",
  },
}
repo.LoadAggregate(loadedAgg)
```
