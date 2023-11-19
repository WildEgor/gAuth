package proto

import (
	"context"
	kcAdapter "github.com/WildEgor/gAuth/internal/adapters/keycloak"
	"github.com/WildEgor/gAuth/internal/repositories"
)

type AuthService struct {
	ka *kcAdapter.KeycloakAdapter
	ur *repositories.UserRepository
}

func NewAuthService(
	ka *kcAdapter.KeycloakAdapter,
	ur *repositories.UserRepository,
) *AuthService {
	return &AuthService{
		ka,
		ur,
	}
}

func (s *AuthService) ValidateToken(ctx context.Context, request *ValidateTokenRequest) (*UserData, error) {
	//TODO implement me
	return &UserData{
		Id: "123",
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
