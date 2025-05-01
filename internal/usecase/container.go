package usecase

import (
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/entity"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/pkg/events"
)

type OrderContainer struct {
	CreateOrderUseCase *CreateOrderUseCase
	GetOrdersUseCase   *GetOrdersUseCase
}

func NewOrderContainer(
	OrderRepository entity.OrderRepositoryInterface,
	EventDispatcher events.EventDispatcherInterface,
	OrderCreatedEvent events.EventInterface,
) *OrderContainer {
	return &OrderContainer{
		CreateOrderUseCase: newCreateOrderUseCase(OrderRepository, OrderCreatedEvent, EventDispatcher),
		GetOrdersUseCase:   newGetOrdersUseCase(OrderRepository),
	}
}
