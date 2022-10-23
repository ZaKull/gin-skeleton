package model

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Search struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	Text      string `json:"source"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Searches []Search

func (Search) TableName() string {
	return "searches"
}

func (Searches) TableName() string {
	return "searches"
}

func (u *Search) GetFirstByText() error {
	now := time.Now()
	db := DB().Order("created_at desc").Where("created_at > ?", now.AddDate(0, 0, -1)).First(u, "text = ?", u.Text)

	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	} else if db.Error != nil {
		return db.Error
	}
	return nil
}

func (u *Search) GetLikeText(id string) error {
	db := DB().Where("id=?", id).First(u)

	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return ErrDataNotFound
	} else if db.Error != nil {
		return db.Error
	}
	return nil
}

func (u *Search) CreateSearch() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}

func (u *Searches) CreateSearchs() error {
	db := DB().Clauses(clause.Insert{Modifier: "IGNORE"}).Create(u)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return ErrKeyConflict
	}

	return nil
}
