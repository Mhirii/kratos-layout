package service

import (
	"context"

	pb "layout/api/users/v1"
	"layout/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
)

type UsersService struct {
	pb.UnimplementedUsersServer
	uc  *biz.UsersUsecase
	log *log.Helper
}

func NewUsersService(uc *biz.UsersUsecase, logger log.Logger) *UsersService {
	return &UsersService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *UsersService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersService.CreateUser")
	defer span.End()
	reqPr := &biz.User{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Phone:    req.GetPhone(),
		Picture:  req.GetPicture(),
	}
	res, err := s.uc.CreateUser(ctx, reqPr)
	if err != nil {
		return nil, err
	}
	resp := &pb.CreateUserResponse{
		Id: res,
	}

	return resp, nil
}

func (s *UsersService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersService.GetUser")
	defer span.End()

	res, err := s.uc.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	resUser := &pb.User{
		Id:       res.ID,
		Username: res.Username,
		Email:    res.Email,
		Phone:    res.Phone,
		Picture:  &res.Picture,
	}

	resp := &pb.GetUserResponse{User: resUser}
	return resp, nil
}

func (s *UsersService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersService.ListUsers")
	defer span.End()
	reqPr := &biz.Pagination{
		Page: req.GetPagination().GetPage(),
		Size: req.GetPagination().GetPageSize(),
	}
	res, err := s.uc.ListUsers(ctx, reqPr)
	if err != nil {
		return nil, err
	}
	var users []*pb.User
	for _, u := range res {
		users = append(users, &pb.User{
			Id:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Phone:    u.Phone,
			Picture:  &u.Picture,
		})
	}
	resp := &pb.ListUsersResponse{
		Users: users,
	}
	return resp, nil
}

func (s *UsersService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersService.UpdateUser")
	defer span.End()
	reqPr := &biz.User{
		ID:       req.GetId(),
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Phone:    req.GetPhone(),
		Picture:  req.GetPicture(),
	}
	res, err := s.uc.UpdateUser(ctx, reqPr)
	if err != nil {
		return nil, err
	}
	resp := &pb.UpdateUserResponse{
		Id: res.ID,
	}
	return resp, nil
}
func (s *UsersService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersService.DeleteUser")
	defer span.End()
	res, err := s.uc.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	resp := &pb.DeleteUserResponse{
		Id: res.ID,
	}
	return resp, nil
}
func (s *UsersService) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	ctx, span := otel.Tracer("users").Start(ctx, "UsersService.SearchUsers")
	defer span.End()
	reqPr := &biz.Pagination{
		Page: req.GetPagination().GetPage(),
		Size: req.GetPagination().GetPageSize(),
	}
	res, err := s.uc.SearchUsers(ctx, req.GetQuery(), reqPr)
	if err != nil {
		return nil, err
	}
	var users []*pb.User
	for _, u := range res {
		users = append(users, &pb.User{
			Id:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Phone:    u.Phone,
			Picture:  &u.Picture,
		})
	}
	resp := &pb.SearchUsersResponse{
		Users: users,
	}
	return resp, nil
}
