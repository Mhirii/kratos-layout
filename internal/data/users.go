package data

import (
	"context"
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/go-kratos/kratos/v2/log"
)

type Users struct {
	//gorm.Model
	//ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	//Username string    `gorm:"not null;uniqueIndex"`
	//Email    string    `gorm:"not null;uniqueIndex"`
	//Phone    *string   `gorm:"not null;uniqueIndex"`
	//Avatar   *string
	ID       uuid.UUID
	Username string
}

type usersRepo struct {
	data *Data
	log  *log.Helper
}

func NewUsersRepo(data *Data, logger log.Logger) biz.UsersRepo {
	return &usersRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *usersRepo) Save(ctx context.Context, u *biz.Users) (*biz.Users, error) {
	_, span := otel.Tracer("users").Start(ctx, "Data Save")
	defer span.End()
	user := &Users{
		Username: *u.Username,
		//Email:    *u.Email,
		//Phone:    u.Phone,
		//Avatar:   u.Avatar,
	}
	//t := r.data.client.Save(user)

	//if t.Error != nil {
	//	err := MapDBErrors(t.Error)
	//	span.End()
	//	return nil, err
	//}
	resp := &biz.Users{
		ID:       user.ID.String(),
		Username: &user.Username,
		//Email:     &user.Email,
		//Phone:     user.Phone,
		//Avatar:    user.Avatar,
		//CreatedAt: &user.CreatedAt,
		//UpdatedAt: &user.UpdatedAt,
		//DeletedAt: &user.DeletedAt.Time,
	}
	return resp, nil
}

func (r *usersRepo) Update(ctx context.Context, u *biz.Users) (*biz.Users, error) {
	_, span := otel.Tracer("users").Start(ctx, "Data Update")
	defer span.End()

	// uid, err := uuid.Parse(u.ID)
	_, err := uuid.Parse(u.ID)

	if err != nil {
		r.log.Warn(err.Error(), zap.Error(err))
		span.AddEvent(err.Error())
		return nil, err
	}
	user := &Users{
		Username: *u.Username,
		//Email:    *u.Email,
		//Phone:    u.Phone,
		//Avatar:   u.Avatar,
	}
	//t := r.data.client.Model(&Users{}).Where("id = ?", uid).Updates(user)
	//if t.Error != nil {
	//	err := MapDBErrors(t.Error)
	//	return nil, err
	//}

	resp := &biz.Users{
		ID:       user.ID.String(),
		Username: &user.Username,
		//Email:     &user.Email,
		//Phone:     user.Phone,
		//Avatar:    user.Avatar,
		//CreatedAt: &user.CreatedAt,
		//UpdatedAt: &user.UpdatedAt,
		//DeletedAt: &user.DeletedAt.Time,
	}
	return resp, nil
}

func (r *usersRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.Users, error) {
	_, span := otel.Tracer("users").Start(ctx, "Data FindByID")
	defer span.End()
	user := &Users{
		ID: id,
	}
	//t := r.data.client.First(&user)
	//if t.Error != nil {
	//	err := MapDBErrors(t.Error)
	//	return nil, err
	//}
	resp := &biz.Users{
		ID:       user.ID.String(),
		Username: &user.Username,
		//Email:    &user.Email,
		//Phone:    user.Phone,
		//Avatar:   user.Avatar,
	}
	return resp, nil
}

func (r *usersRepo) ListAll(ctx context.Context, pp biz.PaginationParams, sp biz.SortParams) ([]biz.Users, error) {
	_, span := otel.Tracer("users").Start(ctx, "Data ListAll")
	defer span.End()
	//offset := pp.PageSize * pp.Page

	var usersList []Users
	//q := r.data.client.Offset(offset).Limit(pp.PageSize)
	//if sp.SortBy != "" {
	//sortOrder := "asc"
	//if sp.SortOrder != "asc" && sp.SortOrder != "desc" {
	//	sortOrder = "asc"
	//} else {
	//	sortOrder = sp.SortOrder
	//}
	//q = q.Order(sp.SortBy + " " + sortOrder)
	//}
	//q = q.Find(&usersList)

	//if q.Error != nil {
	//	err := MapDBErrors(q.Error)
	//	return nil, err
	//}

	var result []biz.Users
	for _, user := range usersList {
		result = append(result, biz.Users{
			ID:       user.ID.String(),
			Username: &user.Username,
			//Email:     &user.Email,
			//Phone:     user.Phone,
			//Avatar:    user.Avatar,
			//CreatedAt: &user.CreatedAt,
			//UpdatedAt: &user.UpdatedAt,
			//DeletedAt: &user.DeletedAt.Time,
		})
	}

	res := make([]biz.Users, len(result))
	copy(res, result)

	return result, nil
}

func (r *usersRepo) Delete(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	_, span := otel.Tracer("users").Start(ctx, "Data Delete")
	defer span.End()
	//err := r.data.client.Delete(&Users{}, id).Error
	//if err != nil {
	//	err = MapDBErrors(err)
	//	return id, err
	//}
	return id, nil
}

func (r *usersRepo) Count(ctx context.Context) (int, error) {
	_, span := otel.Tracer("users").Start(ctx, "Data Count")
	defer span.End()
	var count int64

	/* DELETE */
	count = 0
	/* DELETE */

	//t := r.data.client.Model(&Users{}).Count(&count)
	//if t.Error != nil {
	//	err := MapDBErrors(t.Error)
	//	return 0, err
	//}
	return int(count), nil
}
