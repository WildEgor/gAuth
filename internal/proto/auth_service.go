package proto

import (
	"context"
	"github.com/WildEgor/gAuth/internal/repositories"
)

type AuthService struct {
	ur *repositories.UserRepository
}

func NewAuthService(
	ur *repositories.UserRepository,
) *AuthService {
	return &AuthService{
		ur,
	}
}

func (s *AuthService) ValidateToken(ctx context.Context, request *ValidateTokenRequest) (*UserData, error) {
	// TODO: impl
	ur, err := s.ur.FindById("")
	if err != nil {
		return nil, err
	}

	// TODO: add mapper
	return &UserData{
		Id:        ur.Id.Hex(),
		FirstName: ur.FirstName,
		LastName:  ur.LastName,
		Email:     ur.Email,
		Phone:     ur.Phone,
		IsActive:  true, // TODO
	}, nil
}

func (s *AuthService) FindByIds(ctx context.Context, request *FindByIdsRequest) (*FindByIdsResponse, error) {
	var response FindByIdsResponse

	if len(request.Ids) <= 0 {
		return &response, nil
	}

	users, err := s.ur.FindByIds(request.Ids)
	if err != nil {
		return &response, nil
	}

	// TODO: add mapper
	for _, model := range *users {
		response.Users = append(response.Users, &UserData{
			Id:        model.Id.Hex(),
			FirstName: model.FirstName,
			LastName:  model.LastName,
			Email:     model.Email,
			Phone:     model.Phone,
			IsActive:  true, // TODO
		})
	}

	return &response, nil
}

func (s *AuthService) mustEmbedUnimplementedAuthServiceServer() {
	//TODO implement me
	panic("implement me")
}
