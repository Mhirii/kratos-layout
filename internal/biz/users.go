package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Picture   string `json:"picture"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at,omitempty"`
}

type UsersRepo interface {
	Save(ctx context.Context, u *User) (string, error)
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context, pagination *Pagination) ([]*User, error)
	Update(ctx context.Context, u *User) (*User, error)
	Delete(ctx context.Context, id string) (*User, error)
	Search(ctx context.Context, keyword string, pagination *Pagination) ([]*User, error)
}

type UsersUsecase struct {
	repo UsersRepo
	log  *log.Helper
}

func NewUsersUsecase(repo UsersRepo, logger log.Logger) *UsersUsecase {
	return &UsersUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *UsersUsecase) CreateUser(ctx context.Context, u *User) (string, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersUsecase.CreateUser")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "user",
		Value: attribute.StringValue(u.Username + " " + u.Email + " " + u.Phone),
	})

	res, err := uc.repo.Save(ctx, u)
	if err != nil {
		return "", err
	}
	if res == "" {
		return "", errors.InternalServer("failed to save user", "err was empty but insertions failed")
	}
	return res, nil
}
func (uc *UsersUsecase) GetUser(ctx context.Context, id string) (*User, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersUsecase.GetUser")
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

func (uc *UsersUsecase) ListUsers(ctx context.Context, p *Pagination) ([]*User, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersUsecase.ListUsers")
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

func (uc *UsersUsecase) UpdateUser(ctx context.Context, u *User) (*User, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersUsecase.UpdateUser")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "user",
		Value: attribute.StringValue(u.Username + " " + u.Email + " " + u.Phone),
	})

	res, err := uc.repo.Update(ctx, u)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (uc *UsersUsecase) DeleteUser(ctx context.Context, id string) (*User, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersUsecase.DeleteUser")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "id",
		Value: attribute.StringValue(id),
	})

	res, err := uc.repo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (uc *UsersUsecase) SearchUsers(ctx context.Context, keyword string, p *Pagination) ([]*User, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersUsecase.SearchUsers")
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
