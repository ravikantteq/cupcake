package manager

import (
	"github.com/ravikantteq/cupcake/backyard/internal/store"
)

// Managers holds all the business logic managers
type Managers struct {
	Consumer *ConsumerManager
	Flow     *FlowManager
	Producer *ProducerManager
}

// NewManagers creates all managers with the given store and configuration
func NewManagers(store store.Store, kafkaBroker string) *Managers {
	return &Managers{
		Consumer: NewConsumerManager(store, kafkaBroker),
		Flow:     NewFlowManager(store, kafkaBroker),
		Producer: NewProducerManager(store),
	}
}
