package mappers

import (
	"github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/models"
)

func CreateUser(dto *auth.RegistrationRequestDto) *models.UsersModel {
	return &models.UsersModel{
		Email:     dto.Email,
		Phone:     dto.Phone,
		Password:  dto.Password,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}
}
