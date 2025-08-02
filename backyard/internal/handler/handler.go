package handler

import (
	"github.com/ravikantteq/cupcake/backyard/internal/manager"
)

// Handlers holds all HTTP handlers
type Handlers struct {
	Consumer *ConsumerHandler
	Flow     *FlowHandler
	Producer *ProducerHandler
	Health   *HealthHandler
}

// NewHandlers creates all handlers with the given managers
func NewHandlers(mgrs *manager.Managers) *Handlers {
	return &Handlers{
		Consumer: NewConsumerHandler(mgrs.Consumer),
		Flow:     NewFlowHandler(mgrs.Flow),
		Producer: NewProducerHandler(mgrs.Producer),
		Health:   NewHealthHandler(),
	}
}
