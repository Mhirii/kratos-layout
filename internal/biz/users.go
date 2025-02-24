package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Users struct {
	ID        string
	Username  *string
	Email     *string
	Phone     *string
	Avatar    *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type UsersRepo interface {
	Save(context.Context, *Users) (*Users, error)
	Update(context.Context, *Users) (*Users, error)
	FindByID(context.Context, uuid.UUID) (*Users, error)
	ListAll(context.Context, PaginationParams, SortParams) ([]Users, error)
	Delete(context.Context, uuid.UUID) (uuid.UUID, error)
	Count(ctx context.Context) (int, error)
}

type UsersUsecase struct {
	repo UsersRepo
	log  *log.Helper
}

type ListUsersResponse struct {
	PaginationParams
	Total      int
	TotalPages int
	Users      []Users
}

// NewUsersUsecase new a Users usecase.
func NewUsersUsecase(repo UsersRepo, logger log.Logger) *UsersUsecase {
	return &UsersUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UsersUsecase) CreateUsers(ctx context.Context, u *Users) (*Users, error) {
	_, span := otel.Tracer("users").Start(ctx, "Biz CreateUsers")
	defer span.End()
	res, err := uc.repo.Save(ctx, u)
	if err != nil {
		span.AddEvent(err.Error())
		return nil, err
	}
	return res, nil
}

func (uc *UsersUsecase) GetByID(ctx context.Context, id string) (*Users, error) {
	_, span := otel.Tracer("users").Start(ctx, "Biz GetByID")
	defer span.End()
	uid, err := uuid.Parse(id)
	if err != nil {
		span.AddEvent(err.Error())
		return nil, err
	}
	res, err := uc.repo.FindByID(ctx, uid)
	if err != nil {
		span.AddEvent(err.Error())
		return nil, err
	}
	return res, nil
}

func (uc *UsersUsecase) DeleteUsers(ctx context.Context, id string) (string, error) {
	_, span := otel.Tracer("users").Start(ctx, "Biz DeleteUsers")
	defer span.End()
	uid, err := uuid.Parse(id)
	if err != nil {
		span.AddEvent(err.Error())
		return "", err
	}
	res, err := uc.repo.Delete(ctx, uid)
	if err != nil {
		span.AddEvent(err.Error())
		return "", err
	}
	return res.String(), nil
}

func (uc *UsersUsecase) ListUsers(ctx context.Context, pp PaginationParams, sp SortParams) (ListUsersResponse, error) {
	_, span := otel.Tracer("users").Start(ctx, "Biz ListUsers")
	defer span.End()
	if sp.SortOrder != "" && sp.SortBy != "" {
		if sp.SortOrder != "asc" && sp.SortOrder != "desc" {
			span.AddEvent("invalid sort order")
			return ListUsersResponse{}, errors.BadRequest("users.list", "invalid sort order")
		}
	}

	res, err := uc.repo.ListAll(ctx, pp, sp)
	if err != nil {
		span.AddEvent(err.Error())
		return ListUsersResponse{}, err
	}

	totalPages, err := uc.repo.Count(ctx)
	if err != nil {
		span.AddEvent(err.Error())
		return ListUsersResponse{}, err
	}

	return ListUsersResponse{
		Users:      res,
		Total:      len(res),
		TotalPages: totalPages,
	}, nil
}

func (uc *UsersUsecase) UpdateUsers(ctx context.Context, u *Users) (*Users, error) {
	_, span := otel.Tracer("users").Start(ctx, "Biz UpdateUsers")
	defer span.End()
	if (u.Username == nil || *u.Username == "") || (u.Email == nil || *u.Email == "") || (u.Phone == nil || *u.Phone == "") {
		old, err := uc.GetByID(ctx, u.ID)
		if err != nil {
			span.AddEvent(err.Error())
			return nil, err
		}
		if old == nil {
			err := errors.NotFound("users.update", "user not found")
			span.AddEvent(err.Error())
			return nil, err
		}
		if u.Username == nil || *u.Username == "" {
			u.Username = old.Username
		}
		if u.Email == nil || *u.Email == "" {
			u.Email = old.Email
		}
		if u.Phone == nil || *u.Phone == "" {
			u.Phone = old.Phone
		}
	}

	res, err := uc.repo.Update(ctx, u)
	if err != nil {
		span.AddEvent(err.Error())
		return nil, err
	}
	return res, nil
}
