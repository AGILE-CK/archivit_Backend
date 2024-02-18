package auditlog

import "time"

// AuditLog audit log
type BaseModel struct {
	ID        int64      `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `sql:"index"`
}
