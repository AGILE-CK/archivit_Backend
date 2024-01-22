package entity

import "archivit_Backend/src/db/entity/auditlog"

type User struct {
	Username string `gorm:"type:varchar(50);not null"`
	Password string `gorm:"type:varchar(50);not null"`
	auditlog.BaseModel
}
