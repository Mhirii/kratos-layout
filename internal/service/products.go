package service

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	pb "layout/api/products/v1"
	"layout/internal/biz"
)

type ProductsService struct {
	pb.UnimplementedProductsServer
	uc  *biz.ProductsUsecase
	log *log.Helper
}

func NewProductsService(uc *biz.ProductsUsecase, logger log.Logger) *ProductsService {
	return &ProductsService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *ProductsService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	_, span := otel.Tracer("products").Start(ctx, "CreateProduct")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "product",
		Value: attribute.StringValue(fmt.Sprintf("Name: %s, Desc %s, Category: %s, Price: %f", req.GetName(), req.GetDescription(), req.GetCategory(), req.GetPrice())),
	})
	bizProd := &biz.Product{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Category:    req.GetCategory(),
		Tags:        req.GetTags(),
		Attributes:  req.GetAttributes(),
		Images:      req.GetImages(),
	}
	thumbnail := req.GetThumbnail()
	if thumbnail != "" {
		bizProd.Thumbnail = &thumbnail
	}
	res, err := s.uc.CreateProduct(ctx, bizProd)
	if err != nil {
		return nil, err
	}
	resp := &pb.CreateProductResponse{
		Id: res,
	}
	return resp, nil
}
func (s *ProductsService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	_, span := otel.Tracer("products").Start(ctx, "GetProduct")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "id",
		Value: attribute.StringValue(req.GetId()),
	})
	res, err := s.uc.GetProduct(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	result := &pb.Product{
		Id:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Price:       res.Price,
		Category:    res.Category,
		Tags:        res.Tags,
		Attributes:  res.Attributes,
		Thumbnail:   res.Thumbnail,
		Images:      res.Images,
	}
	resp := &pb.GetProductResponse{
		Product: result,
	}
	return resp, nil
}
func (s *ProductsService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	_, span := otel.Tracer("products").Start(ctx, "ListProducts")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "pagination",
		Value: attribute.StringValue(fmt.Sprintf("Page: %d Size: %d", req.GetPagination().GetPage(), req.GetPagination().GetPageSize())),
	})
	var page int32 = 0
	var pageSize int32 = 10
	reqPagination := req.GetPagination()
	if reqPagination != nil {
		page = reqPagination.GetPage()
		pageSize = reqPagination.GetPageSize()
	}

	pagination := &biz.Pagination{
		Page: page,
		Size: pageSize,
	}

	res, err := s.uc.ListProducts(ctx, pagination)
	if err != nil {
		return nil, err
	}
	var products []*pb.Product
	for _, p := range res {
		products = append(products, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Category:    p.Category,
			Tags:        p.Tags,
			Attributes:  p.Attributes,
			Thumbnail:   p.Thumbnail,
			Images:      p.Images,
		})
	}
	resp := &pb.ListProductsResponse{
		Products: products,
	}
	return resp, nil
}
func (s *ProductsService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	_, span := otel.Tracer("products").Start(ctx, "UpdateProduct")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "product",
		Value: attribute.StringValue(fmt.Sprintf("Name: %s, Desc %s, Category: %s, Price: %f", req.GetName(), req.GetDescription(), req.GetCategory(), req.GetPrice())),
	})
	bizProd := &biz.Product{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Category:    req.GetCategory(),
		Tags:        req.GetTags(),
		Attributes:  req.GetAttributes(),
		Images:      req.GetImages(),
	}
	thumbnail := req.GetThumbnail()
	if thumbnail != "" {
		bizProd.Thumbnail = &thumbnail
	}
	res, err := s.uc.UpdateProduct(ctx, bizProd)
	if err != nil {
		return nil, err
	}
	resp := &pb.UpdateProductResponse{
		Id: res.ID,
	}
	return resp, nil
}
func (s *ProductsService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	_, span := otel.Tracer("products").Start(ctx, "DeleteProduct")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "id",
		Value: attribute.StringValue(req.GetId()),
	})
	res, err := s.uc.DeleteProduct(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	resp := &pb.DeleteProductResponse{
		Id: res,
	}
	return resp, nil
}
func (s *ProductsService) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	_, span := otel.Tracer("products").Start(ctx, "SearchProducts")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "keyword",
		Value: attribute.StringValue(req.GetQuery()),
	})
	span.SetAttributes(attribute.KeyValue{
		Key:   "pagination",
		Value: attribute.StringValue(fmt.Sprintf("Page: %d Size: %d", req.GetPagination().GetPage(), req.GetPagination().GetPageSize())),
	})
	var page int32 = 0
	var pageSize int32 = 10
	reqPagination := req.GetPagination()
	if reqPagination != nil {
		page = reqPagination.GetPage()
		pageSize = reqPagination.GetPageSize()
	}

	pagination := &biz.Pagination{
		Page: page,
		Size: pageSize,
	}

	res, err := s.uc.SearchProducts(ctx, req.GetQuery(), pagination)
	if err != nil {
		return nil, err
	}
	var products []*pb.Product
	for _, p := range res {
		products = append(products, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Category:    p.Category,
			Tags:        p.Tags,
			Attributes:  p.Attributes,
			Thumbnail:   p.Thumbnail,
			Images:      p.Images,
		})
	}
	resp := &pb.SearchProductsResponse{
		Products: products,
	}
	return resp, nil
}
