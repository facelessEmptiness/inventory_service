package grpc

import (
	"context"
	"github.com/facelessEmptiness/inventory_service/internal/domain"
	"github.com/facelessEmptiness/inventory_service/internal/usecase"

	pb "github.com/facelessEmptiness/inventory_service/proto"
)

type ProductHandler struct {
	pb.UnimplementedInventoryServiceServer
	uc *usecase.ProductUseCase
}

func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

func (h *ProductHandler) AddProduct(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	p := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryId,
	}
	id, err := h.uc.AddProduct(p)
	if err != nil {
		return nil, err
	}
	p.ID = id
	return &pb.ProductResponse{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CategoryId:  p.CategoryID,
	}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.ProductID) (*pb.ProductResponse, error) {
	p, err := h.uc.GetProduct(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ProductResponse{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CategoryId:  p.CategoryID,
	}, nil
}
