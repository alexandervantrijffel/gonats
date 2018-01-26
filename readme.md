CQRS Event Sourcing with NATS Streaming

Not only is Go an excellent choice for writing system applications and network services, such as the container platform Docker, the container orchestrator Kubernetes and the network proxy Traefik, it is also very suitable to build business oriented applications. This post describes how to implement the business logic of an application based on CQRS Event Sourcing. 

If you are new to CQRS Event Sourcing and want to take a deep dive, watch this 6 hours workshop by the inventor of CQRS Event Sourcing Greg Young. 

Library for applying the event sourcing pattern in Golang using durable event streams from NATS streaming. Load aggregate roots from event streams stored in NATS streaming. Persist new events from the aggregates which can be used to build read models. 

The library is inspired by the GetEventStore.com client API. 