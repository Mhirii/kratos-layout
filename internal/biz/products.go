package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Product struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float32           `json:"price"`
	Category    string            `json:"category"`
	Tags        []string          `json:"tags"`
	Attributes  map[string]string `json:"attributes"`
	Thumbnail   *string           `json:"thumbnail"`
	Images      []string          `json:"images"`
}

type ProductsRepo interface {
	Save(ctx context.Context, p *Product) (string, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	List(ctx context.Context, pagination *Pagination) ([]*Product, error)
	Update(ctx context.Context, p *Product) (*Product, error)
	Delete(ctx context.Context, id string) (string, error)
	Search(ctx context.Context, keyword string, pagination *Pagination) ([]*Product, error)
}

type ProductsUsecase struct {
	repo ProductsRepo
	log  *log.Helper
}

func NewProductsUsecase(repo ProductsRepo, logger log.Logger) *ProductsUsecase {
	return &ProductsUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *ProductsUsecase) CreateProduct(ctx context.Context, p *Product) (string, error) {
	_, span := otel.Tracer("products").Start(ctx, "CreateProduct")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "product",
		Value: attribute.StringValue(fmt.Sprintf("Name: %s, Desc %s, Category: %s, Price: %f", p.Name, p.Description, p.Category, p.Price)),
	})
	res, err := uc.repo.Save(ctx, p)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (uc *ProductsUsecase) GetProduct(ctx context.Context, id string) (*Product, error) {
	_, span := otel.Tracer("products").Start(ctx, "GetProduct")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "id",
		Value: attribute.StringValue(id),
	})
	res, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (uc *ProductsUsecase) ListProducts(ctx context.Context, p *Pagination) ([]*Product, error) {
	_, span := otel.Tracer("products").Start(ctx, "ListProducts")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "pagination",
		Value: attribute.StringValue(fmt.Sprintf("Page: %d Size: %d", p.Page, p.Size)),
	})
	res, err := uc.repo.List(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (uc *ProductsUsecase) UpdateProduct(ctx context.Context, p *Product) (*Product, error) {
	_, span := otel.Tracer("products").Start(ctx, "UpdateProduct")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "product",
		Value: attribute.StringValue(fmt.Sprintf("Name: %s, Desc %s, Category: %s, Price: %f", p.Name, p.Description, p.Category, p.Price)),
	})
	res, err := uc.repo.Update(ctx, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (uc *ProductsUsecase) DeleteProduct(ctx context.Context, id string) (string, error) {
	_, span := otel.Tracer("products").Start(ctx, "DeleteProduct")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "id",
		Value: attribute.StringValue(id),
	})
	res, err := uc.repo.Delete(ctx, id)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (uc *ProductsUsecase) SearchProducts(ctx context.Context, keyword string, p *Pagination) ([]*Product, error) {
	_, span := otel.Tracer("products").Start(ctx, "SearchProducts")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "keyword",
		Value: attribute.StringValue(keyword),
	})
	span.SetAttributes(attribute.KeyValue{
		Key:   "pagination",
		Value: attribute.StringValue(fmt.Sprintf("Page: %d Size: %d", p.Page, p.Size)),
	})
	res, err := uc.repo.Search(ctx, keyword, p)
	if err != nil {
		return nil, err
	}
	return res, nil
}
