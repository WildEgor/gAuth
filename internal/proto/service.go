package proto

import (
	kcAdapter "github.com/WildEgor/gAuth/internal/adapters/keycloak"
	"github.com/WildEgor/gAuth/internal/dtos/rpc"
)

type GRPCService struct {
	ka *kcAdapter.KeycloakAdapter
}

func NewGRPCService(ka *kcAdapter.KeycloakAdapter) *GRPCService {
	return &GRPCService{
		ka,
	}
}

func (s *GRPCService) ValidateToken(req *rpc.ValidateTokenPayloadDto) (int, error) {
	return 1, nil
}
