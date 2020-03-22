package user

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"time"
)

// Base model containing the common columns for all tables
type BaseModel struct {
	ID        uuid.UUID `sql:"type:char(36);primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Before storing into the DB, mutate the ID to UUID
func (bm *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}
