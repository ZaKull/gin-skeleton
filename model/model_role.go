package model

import (
	"time"
)

// User the user model
type Role struct {
	ID          uint8     `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TableName for gorm
func (Role) TableName() string {
	return "roles"
}

// GetFirstByID gets the user by his ID
func (r *Role) GetFirstByID(id string) error {
	db := DB().Where("id=?", id).First(r)

	if db.RecordNotFound() {
		return ErrDataNotFound
	} else if db.Error != nil {
		return db.Error
	}

	return nil
}
