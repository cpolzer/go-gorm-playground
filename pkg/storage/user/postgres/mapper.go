package postgres

import (
	userModel "cpolzer.de/m/v2/pkg/storage/user"
	"gorm.io/gorm"
)

func mapToUserDto(u user) *userModel.UserDto {
	return &userModel.UserDto{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Address: userModel.AddressDto{
			ID:         u.Address.ID,
			Street:     u.Address.Street,
			PostalCode: u.Address.PostalCode,
		},
	}
}

func mapToUserEntity(dto userModel.UserDto) *user {
	return &user{
		Model: gorm.Model{
			ID: dto.ID,
		},
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Address: address{
			Model: gorm.Model{
				ID: dto.Address.ID,
			},
			Street:     dto.Address.Street,
			PostalCode: dto.Address.PostalCode,
		},
	}
}
