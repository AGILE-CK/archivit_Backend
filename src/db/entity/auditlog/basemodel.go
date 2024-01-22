package auditlog

import "time"

// AuditLog audit log
type BaseModel struct {
	ID        int64 `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
