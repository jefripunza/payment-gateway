package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UUID v7 string primary keys — GORM does not fill string PKs automatically,
// so every model gets a BeforeCreate hook.

type User struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:120;not null" json:"name"`
	Email        string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:20;not null;default:admin" json:"role"`
	Active       bool      `gorm:"not null;default:true" json:"active"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Wallet struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:120;not null" json:"name"`
	Currency  string    `gorm:"size:8;not null;default:IDR" json:"currency"`
	Balance   int64     `gorm:"not null;default:0" json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Provider struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:120;not null" json:"name"`
	Method    string    `gorm:"size:80;not null" json:"method"` // "provider|value" e.g. "midtrans|QRIS"
	Creds     string    `gorm:"type:text;not null;default:'{}'" json:"-"` // encrypted JSON map of credential values
	Enabled   bool      `gorm:"not null;default:true" json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	return fillID(&u.ID)
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	return fillID(&w.ID)
}

func (p *Provider) BeforeCreate(tx *gorm.DB) error {
	return fillID(&p.ID)
}

func fillID(id *string) error {
	if *id == "" {
		u, err := uuid.NewV7()
		if err != nil {
			return err
		}
		*id = u.String()
	}
	return nil
}
