package user

import (
	"archivit_Backend/src/db"
	"archivit_Backend/src/db/entity"
	"fmt"
	"github.com/jinzhu/gorm"
)

func FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User

	ds := db.InitAndGetDataSource()
	result := ds.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, nil
		}
		return nil, result.Error
	}

	fmt.Printf("CreatedAt value: %#v\n", user.CreatedAt)

	return &user, nil
}

func SaveUser(user *entity.User) error {
	ds := db.InitAndGetDataSource()

	if err := ds.Save(&user).Error; err != nil {
		return err
	}
	return nil
}
