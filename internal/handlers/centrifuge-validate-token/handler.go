package centrifuge_validate_token_handler

import (
	"github.com/WildEgor/gAuth/internal/proto"
)

type CentrifugeValidateTokenHandler struct{}

func NewCentrifugeValidateToken() *CentrifugeValidateTokenHandler {
	return &CentrifugeValidateTokenHandler{}
}

func (hch *CentrifugeValidateTokenHandler) Handle(req *proto.AddRequest) (int, error) {
	return 1, nil
}
