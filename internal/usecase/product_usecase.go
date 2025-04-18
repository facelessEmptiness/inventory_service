package usecase

import (
	"github.com/facelessEmptiness/inventory_service/internal/domain"
	"github.com/facelessEmptiness/inventory_service/internal/repository"
)

type ProductUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(r repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: r}
}

func (uc *ProductUseCase) AddProduct(p *domain.Product) (string, error) {
	return uc.repo.Create(p)
}

func (uc *ProductUseCase) GetProduct(id string) (*domain.Product, error) {
	return uc.repo.GetByID(id)
}
