package user

import (
	"archivit_Backend/src/db"
	"archivit_Backend/src/db/entity"
	"crypto/sha256"
	"encoding/hex"
	"github.com/jinzhu/gorm"
)

func FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User

	ds := db.GetDataSource()
	if err := ds.Where("email = ?", email).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func SaveUser(user *entity.User) error {
	ds := db.GetDataSource()

	if user.LoginType == "NORMAL" {
		hash := sha256.Sum256([]byte(user.Password))
		user.Password = hex.EncodeToString(hash[:])
	}

	if err := ds.Save(&user).Error; err != nil {
		return err
	}
	return nil
}
