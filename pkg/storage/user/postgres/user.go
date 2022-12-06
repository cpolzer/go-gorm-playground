package postgres

import (
	"context"
	"fmt"

	userModel "cpolzer.de/m/v2/pkg/storage/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const timeZone = "Europe/Berlin"

type pg struct {
	postgreClient *gorm.DB
}

func (p pg) Has(ctx context.Context, userId uint) (bool, error) {
	var exists bool
	err := p.postgreClient.Model(&user{}).Select("count(*) > 0").Where("ID = ?", userId).Find(&exists).Error
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	return false, nil
}

func (p pg) Get(ctx context.Context, userId uint) (*userModel.UserDto, error) {
	user := user{}
	err := p.postgreClient.
		Debug().
		Model(user).
		Preload("Address").
		First(&user, "ID = ?", userId).
		Error
	if err != nil {
		println(ctx, fmt.Errorf("project get() failed. Error was: %s", err))
		return nil, err
	}
	userDto := mapToUserDto(user)

	return userDto, nil
}

func (p pg) Create(ctx context.Context, userDto *userModel.UserDto) (*userModel.UserDto, error) {
	userEntity := mapToUserEntity(*userDto)

	err := p.postgreClient.Debug().Model(userEntity).Create(userEntity).Error
	if err != nil {
		return nil, err
	}
	result := mapToUserDto(*userEntity)
	return result, nil
}

func (p pg) Update(ctx context.Context, userDto *userModel.UserDto) (*userModel.UserDto, error) {
	userEntity := mapToUserEntity(*userDto)

	err := p.postgreClient.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Updates(userEntity).Error
	if err != nil {
		return nil, err
	}
	result := mapToUserDto(*userEntity)
	return result, nil
}

func (p pg) Delete(ctx context.Context, userId uint) error {
	//TODO implement me
	panic("implement me")
}

func New(ctx context.Context, host string, port string) (userModel.Repository, error) {
	dsnTemplate := "host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s sslmode=disable"

	dsn := fmt.Sprintf(dsnTemplate,
		host,
		"postgres",
		"postgres",
		"testdb",
		port,
		timeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&user{}, &address{})
	if err != nil {
		return nil, err
	}
	return pg{
		postgreClient: db,
	}, nil
}
