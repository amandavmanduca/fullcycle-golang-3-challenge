package web

import (
	"encoding/json"
	"net/http"

	"github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/usecase"
)

type WebOrderHandler struct {
	useCaseContainer usecase.OrderContainer
}

func NewWebOrderHandler(
	useCaseContainer usecase.OrderContainer,
) *WebOrderHandler {
	return &WebOrderHandler{
		useCaseContainer: useCaseContainer,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := h.useCaseContainer.CreateOrderUseCase
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	output, err := h.useCaseContainer.GetOrdersUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
