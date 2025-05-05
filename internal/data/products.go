package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"layout/internal/biz"
)

type Products struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"desc"`
	Price       float32            `bson:"price"`
	Category    string             `bson:"category"`
	Tags        []string           `bson:"tags"`
	Attributes  map[string]string  `bson:"attributes"`
	Thumbnail   string             `bson:"thumbnail"`
	Images      []string           `bson:"images"`
}

type productsRepo struct {
	db   *mongo.Database
	log  *log.Helper
	coll *mongo.Collection
	tp   trace.Tracer
}

func NewProductsRepo(data Data, logger log.Logger, tp trace.Tracer) (biz.ProductsRepo, error) {
	m := data.GetMongoDB()
	lg := log.NewHelper(logger)

	if m == nil {
		lg.Warn("No Database configuration found for MongoDB, skipping MongoDB initialization")
		return nil, errors.InternalServer("MongoDB is not configured", "MongoDB is not configured")
	}

	return &productsRepo{
		db:   m,
		log:  lg,
		coll: m.Collection("products"),
		tp:   tp,
	}, nil
}

func (r productsRepo) Save(ctx context.Context, p *biz.Product) (string, error) {
	_, span := otel.Tracer("products").Start(ctx, "Save")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "product",
		Value: attribute.StringValue(fmt.Sprintf("Name: %s, Desc %s, Category: %s, Price: %f", p.Name, p.Description, p.Category, p.Price)),
	})
	product := Products{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Tags:        p.Tags,
		Attributes:  p.Attributes,
		Images:      p.Images,
	}
	if p.Thumbnail != nil {
		product.Thumbnail = *p.Thumbnail
	}
	res, err := r.coll.InsertOne(ctx, product)
	if err != nil {
		r.log.Error("failed to save product", err)
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (r productsRepo) GetByID(ctx context.Context, id string) (*biz.Product, error) {
	_, span := otel.Tracer("products").Start(ctx, "GetByID")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "id",
		Value: attribute.StringValue(id),
	})
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.log.Error("failed to parse product id", err)
		return nil, err
	}

	res := r.coll.FindOne(ctx, bson.M{"_id": idObj})
	if res.Err() != nil {
		r.log.Error("failed to get product", res.Err())
		return nil, res.Err()
	}
	if res.Err() != nil {
		r.log.Error("failed to get product", res.Err())
		return nil, res.Err()
	}
	var p Products
	err = res.Decode(&p)
	if err != nil {
		r.log.Error("failed to decode product", err)
		return nil, err
	}
	return &biz.Product{
		ID:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Tags:        p.Tags,
		Attributes:  p.Attributes,
		Thumbnail:   &p.Thumbnail,
		Images:      p.Images,
	}, nil
}

func (r productsRepo) List(ctx context.Context, pagination *biz.Pagination) ([]*biz.Product, error) {
	_, span := otel.Tracer("products").Start(ctx, "List")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "pagination",
		Value: attribute.StringValue(fmt.Sprintf("Page: %d Size: %d", pagination.Page, pagination.Size)),
	})
	offset := pagination.Page * pagination.Size
	take := pagination.Size
	if offset < 0 {
		offset = 0
	}
	if take < 0 {
		take = 0
	}

	r.log.Infof("ListProducts %d %d", offset, take)
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(take))
	cur, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		r.log.Error("failed to list products", err)
		return nil, err
	}
	var res []*biz.Product
	for cur.Next(ctx) {
		var p Products
		if err := cur.Decode(&p); err != nil {
			r.log.Error("failed to decode product", err)
			return nil, err
		}
		res = append(res, &biz.Product{
			ID:          p.ID.Hex(),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Category:    p.Category,
			Tags:        p.Tags,
			Attributes:  p.Attributes,
			Thumbnail:   &p.Thumbnail,
			Images:      p.Images,
		})
	}
	return res, nil
}

func (r productsRepo) Update(ctx context.Context, p *biz.Product) (*biz.Product, error) {
	_, span := otel.Tracer("products").Start(ctx, "Update")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "product",
		Value: attribute.StringValue(fmt.Sprintf("Name: %s, Desc %s, Category: %s, Price: %f", p.Name, p.Description, p.Category, p.Price)),
	})
	uid, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		r.log.Error("failed to parse product id", err)
		return nil, err
	}
	product := Products{
		ID:          uid,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Tags:        p.Tags,
		Attributes:  p.Attributes,
		Images:      p.Images,
	}
	if p.Thumbnail != nil {
		product.Thumbnail = *p.Thumbnail
	}
	res, err := r.coll.ReplaceOne(ctx, bson.M{"_id": product.ID}, product)
	if err != nil {
		r.log.Error("failed to update product", err)
		return nil, err
	}
	if res.ModifiedCount == 0 {
		r.log.Error("failed to update product", "err was empty but insertions failed")
		return nil, errors.InternalServer("failed to update product", "err was empty but insertions failed")
	}
	return &biz.Product{
		ID:          product.ID.Hex(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		Tags:        product.Tags,
		Attributes:  product.Attributes,
		Thumbnail:   &product.Thumbnail,
		Images:      product.Images,
	}, nil
}

func (r productsRepo) Delete(ctx context.Context, id string) (string, error) {
	_, span := otel.Tracer("products").Start(ctx, "Delete")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "id",
		Value: attribute.StringValue(id),
	})
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.log.Error("failed to parse product id", err)
		return "", err
	}
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": idObj})
	if err != nil {
		r.log.Error("failed to delete product", err)
		return "", err
	}
	if res.DeletedCount == 0 {
		r.log.Error("failed to delete product", "err was empty but insertions failed")
		return "", errors.InternalServer("failed to delete product", "err was empty but insertions failed")
	}
	return id, nil
}

func (r productsRepo) Search(ctx context.Context, keyword string, pagination *biz.Pagination) ([]*biz.Product, error) {
	_, span := otel.Tracer("products").Start(ctx, "Search")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "keyword",
		Value: attribute.StringValue(keyword),
	})
	span.SetAttributes(attribute.KeyValue{
		Key:   "pagination",
		Value: attribute.StringValue(fmt.Sprintf("Page: %d Size: %d", pagination.Page, pagination.Size)),
	})
	var products []Products
	opts := options.Find().SetSkip(int64(pagination.Page * pagination.Size)).SetLimit(int64(pagination.Size))
	cur, err := r.coll.Find(ctx, bson.M{"$text": bson.M{"$search": keyword}}, opts)
	if err != nil {
		r.log.Error("failed to search products", err)
		return nil, err
	}
	if err := cur.All(ctx, &products); err != nil {
		r.log.Error("failed to decode products", err)
		return nil, err
	}
	var res []*biz.Product
	for _, p := range products {
		res = append(res, &biz.Product{
			ID:          p.ID.Hex(),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Category:    p.Category,
			Tags:        p.Tags,
			Attributes:  p.Attributes,
			Thumbnail:   &p.Thumbnail,
			Images:      p.Images,
		})
	}
	return res, nil
}
