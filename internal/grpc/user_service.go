package grpc

import (
	"context"
	"errors"

	"github.com/GDH-Project/auth/internal/domain"
	"github.com/GDH-Project/auth/internal/grpc/userpb"
	"go.uber.org/zap"
)

type userService struct {
	userpb.UnimplementedUserServiceServer
	userUseCase domain.UserUseCase
}

func (s *userService) CheckCreateUser(ctx context.Context, req *userpb.GetCheckCreateUserReqeust) (*userpb.GetCheckCreateUserResponse, error) {
	if err := s.userUseCase.CheckCanCreateByEmailOrName(ctx, req.GetEmail(), req.GetName()); err != nil {
		return &userpb.GetCheckCreateUserResponse{
			Ok: false,
		}, nil
	}

	return &userpb.GetCheckCreateUserResponse{
		Ok: true,
	}, nil
}

func (s *userService) GetUserInfoByEmail(ctx context.Context, req *userpb.GetUserInfoByEmailRequest) (*userpb.GetUserInfoResponse, error) {
	user, err := s.userUseCase.FindUserByEmail(ctx, req.GetEmail())
	if err != nil {
		zap.S().Infow("email에 해당하는 사용자를 찾을 수 없습니다.",
			zap.Error(err),
			"email", req.GetEmail(),
		)
		return nil, errors.New("사용자를 찾을 수 없습니다")
	}

	var userRole userpb.UserRole
	switch user.Role {
	case domain.RoleUser:
		userRole = userpb.UserRole_BASIC_USER
	case domain.RoleDevice:
		userRole = userpb.UserRole_DATA_USER
	case domain.RoleAdmin:
		userRole = userpb.UserRole_ADMIN
	default:
		userRole = userpb.UserRole_BASIC_USER
	}

	return &userpb.GetUserInfoResponse{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   userRole,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	var userRole domain.Role
	switch req.GetType() {
	case userpb.CreateUserType_basic_user:
		userRole = domain.RoleUser
	case userpb.CreateUserType_data_user:
		userRole = domain.RoleDevice
	default:
		userRole = domain.RoleUser
	}

	err := s.userUseCase.CreateUser(ctx, &domain.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     userRole,
	})

	if err != nil {
		zap.S().Infow("사용자 생성중 오류가 발생했습니다.", zap.Error(err))
		return nil, errors.New(err.UserMessage)
	}

	return &userpb.CreateUserResponse{}, nil

}

func (s *userService) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	updateUserInfo := &domain.User{
		ID:       req.GetUserId(),
		Name:     req.GetName(),
		Password: req.GetPassword(),
	}
	err := s.userUseCase.UpdateUserByUserID(ctx, updateUserInfo)
	if err != nil {
		zap.S().Infow("사용자 정보 업데이트중 오류가 발생했습니다.", zap.Error(err))
		return nil, errors.New(err.UserMessage)
	}

	user, err := s.userUseCase.FindUserByUserID(ctx, updateUserInfo.ID)
	if err != nil {
		zap.S().Infow("업데이트된 유저 정보를 불러오는데 실패했습니다.", zap.Error(err))
		return nil, errors.New(err.UserMessage)
	}

	var userRole userpb.UserRole
	switch user.Role {
	case domain.RoleUser:
		userRole = userpb.UserRole_BASIC_USER
	case domain.RoleDevice:
		userRole = userpb.UserRole_DATA_USER
	case domain.RoleAdmin:
		userRole = userpb.UserRole_ADMIN
	default:
		userRole = userpb.UserRole_BASIC_USER
	}

	return &userpb.UpdateUserResponse{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   userRole,
	}, nil
}
func (s *userService) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	err := s.userUseCase.DeleteUserByUserIDAndPassword(ctx, req.GetUserId(), req.GetPassword())
	if err != nil {
		zap.S().Infow("사용자를 삭제하는 중 오류가 발생했습니다.", zap.Error(err))
		return nil, errors.New(err.UserMessage)
	}

	return &userpb.DeleteUserResponse{}, nil
}

func NewGrpcUserService(userUseCase domain.UserUseCase) userpb.UserServiceServer {
	return &userService{userUseCase: userUseCase}
}
