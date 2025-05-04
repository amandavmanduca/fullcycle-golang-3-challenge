package service

import (
	"context"

	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/infra/grpc/pb"
	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	OrderContainer usecase.OrderContainer
}

func NewOrderService(orderContainer usecase.OrderContainer) *OrderService {
	return &OrderService{
		OrderContainer: orderContainer,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.OrderContainer.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) GetOrders(ctx context.Context, b *pb.Blank) (*pb.OrderList, error) {
	orders, err := s.OrderContainer.GetOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}
	var ordersResponse []*pb.Order
	for _, order := range orders {
		ordersResponse = append(ordersResponse, &pb.Order{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		})
	}
	return &pb.OrderList{
		Orders: ordersResponse,
	}, nil
}
