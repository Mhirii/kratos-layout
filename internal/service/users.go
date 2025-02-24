package service

import (
	"context"
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	pb "github.com/go-kratos/kratos-layout/api/users/v1"
)

type UsersService struct {
	pb.UnimplementedUsersServer
	uc  *biz.UsersUsecase
	log *log.Helper
}

func NewUsersService(uc *biz.UsersUsecase, logger log.Logger) *UsersService {
	return &UsersService{uc: uc, log: log.NewHelper(logger)}
}

func (s *UsersService) CreateUsers(ctx context.Context, req *pb.CreateUsersRequest) (*pb.CreateUsersReply, error) {
	_, span := otel.Tracer("users").Start(ctx, "CreateUsers")
	defer span.End()
	username := req.GetUsername()
	email := req.GetEmail()
	phone := req.GetPhone()

	kv := attribute.KeyValue{
		Key:   "email",
		Value: attribute.StringValue(email),
	}
	span.SetAttributes(kv)

	res, err := s.uc.CreateUsers(ctx, &biz.Users{
		Username: &username,
		Email:    &email,
		Phone:    &phone,
	})
	if err != nil {
		s.log.WithContext(ctx).Warnf("CreateUsers: %s", err)
		return nil, err
	}

	resp := &pb.CreateUsersReply{
		Id:       res.ID,
		Username: *res.Username,
		Email:    *res.Email,
		Phone:    res.Phone,
	}
	s.log.WithContext(ctx).Infof("CreateUsers: %s", resp.Id)
	return resp, nil
}
func (s *UsersService) UpdateUsers(ctx context.Context, req *pb.UpdateUsersRequest) (*pb.UpdateUsersReply, error) {
	_, span := otel.Tracer("users").Start(ctx, "UpdateUsers")
	defer span.End()
	id := req.GetId()
	username := req.GetUsername()
	email := req.GetEmail()
	phone := req.GetPhone()

	if id == "" {
		err := errors.BadRequest("users.update", "id is required")
		s.log.WithContext(ctx).Warnf("UpdateUsers: %s", err)
		return nil, err
	}
	if username == "" && email == "" && phone == "" {
		err := errors.BadRequest("users.update", "You must provide at least one field to update")
		s.log.WithContext(ctx).Warnf("UpdateUsers: %s", err)
		return nil, err
	}

	res, err := s.uc.UpdateUsers(ctx, &biz.Users{
		ID:       id,
		Username: &username,
		Email:    &email,
		Phone:    &phone,
	})
	if err != nil {
		s.log.WithContext(ctx).Warnf("UpdateUsers: %s", err)
		return nil, err
	}
	resp := &pb.UpdateUsersReply{
		Id:       res.ID,
		Username: *res.Username,
		Email:    *res.Email,
		Phone:    res.Phone,
	}
	s.log.WithContext(ctx).Infof("UpdateUsers: id %s", resp.Id)
	return resp, nil
}
func (s *UsersService) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersReply, error) {
	_, span := otel.Tracer("users").Start(ctx, "DeleteUsers")
	defer span.End()
	id := req.GetId()
	_, err := s.uc.DeleteUsers(ctx, id)
	if err != nil {
		s.log.WithContext(ctx).Warnf("DeleteUsers: %s", err)
		return nil, err
	}
	resp := &pb.DeleteUsersReply{
		Id: id,
	}
	s.log.WithContext(ctx).Infof("DeleteUsers: id %s", resp.Id)
	return resp, nil
}
func (s *UsersService) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersReply, error) {
	_, span := otel.Tracer("users").Start(ctx, "GetUsers")
	defer span.End()
	id := req.GetId()
	res, err := s.uc.GetByID(ctx, id)
	if err != nil {
		s.log.WithContext(ctx).Warnf("GetUsers: %s", err)
		return nil, err
	}
	resp := &pb.GetUsersReply{
		Id:       res.ID,
		Username: *res.Username,
		Email:    *res.Email,
		Phone:    res.Phone,
	}
	s.log.WithContext(ctx).Infof("GetUsers: id %s", resp.Id)
	return resp, nil
}
func (s *UsersService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersReply, error) {
	_, span := otel.Tracer("users").Start(ctx, "ListUsers")
	defer span.End()
	page := req.GetPage()
	pageSize := req.GetPageSize()
	reverse := req.GetReverse()

	if page < 0 {
		page = 0
		s.log.WithContext(ctx).Warn("page must be positive 0, default to 0")
	}
	if pageSize <= 0 {
		pageSize = 20
		s.log.WithContext(ctx).Warn("pageSize must be greater than 0, default to 20")
	}
	pp := biz.PaginationParams{
		Page:     int(page),
		PageSize: int(pageSize),
		Reverse:  reverse,
	}

	sp := biz.SortParams{}
	if req.SortBy != nil {
		sp.SortBy = *req.SortBy
	}
	if req.SortOrder != nil {
		sp.SortOrder = *req.SortOrder
	}

	res, err := s.uc.ListUsers(ctx, pp, sp)
	if err != nil {
		s.log.WithContext(ctx).Warnf("ListUsers: %s", err)
		return nil, err
	}
	listUsers := make([]*pb.ListUsersUser, len(res.Users))
	for i, user := range res.Users {
		listUsers[i] = &pb.ListUsersUser{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
		}
	}
	resp := &pb.ListUsersReply{
		Users:      listUsers,
		Page:       int32(res.Page),
		PageSize:   int32(res.PageSize),
		Total:      int32(res.Total),
		TotalPages: int32(res.TotalPages),
		Reverse:    res.Reverse,
	}
	s.log.WithContext(ctx).Infof("ListUsers: list of %d elements", len(resp.Users))
	return resp, nil
}
