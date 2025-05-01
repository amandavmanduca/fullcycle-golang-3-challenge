package usecase

import (
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/entity"
)

type GetOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func newGetOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *GetOrdersUseCase {
	return &GetOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *GetOrdersUseCase) Execute() ([]OrderOutputDTO, error) {
	ordersDTO := []OrderOutputDTO{}
	orders, err := c.OrderRepository.Get()
	if err != nil {
		return ordersDTO, err
	}

	for _, order := range orders {
		dto := OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		}
		ordersDTO = append(ordersDTO, dto)
	}

	return ordersDTO, nil
}
