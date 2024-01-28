package entity

import "archivit_Backend/src/db/entity/auditlog"

type User struct {
	Email     string `gorm:"type:varchar(50);not null"`
	Password  string `gorm:"type:varchar(50)"`
	LoginType string `gorm:"type:varchar(50)"`
	auditlog.BaseModel
}
