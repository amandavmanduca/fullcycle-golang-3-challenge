package graph

import "github.com/amandavmanduca/fullcycle-golang-3-challenge/internal/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	OrderContainer usecase.OrderContainer
}
