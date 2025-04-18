package repository

import (
	"github.com/facelessEmptiness/inventory_service/internal/domain"
)

type ProductRepository interface {
	Create(p *domain.Product) (string, error)
	GetByID(id string) (*domain.Product, error)
}
