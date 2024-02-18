package entity

import (
	"archivit_Backend/src/db/entity/auditlog"
)

type User struct {
	Email     string `gorm:"type:varchar(255);not null" json:"email"`
	Password  string `gorm:"type:varchar(255)" json:"password"`
	LoginType string `gorm:"type:varchar(255)" json:"loginType"`
	auditlog.BaseModel
}
