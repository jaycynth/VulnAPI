package model

import "time"

type KYC struct {
	ID             int       `gorm:"primaryKey"`
	UserID         int       `gorm:"index"`
	DocumentType   string    `json:"document_type" validate:"required"`
	DocumentNumber string    `gorm:"unique;not null"`
	IssueDate      time.Time `json:"issue_date"`
	ExpiryDate     time.Time `json:"expiry_date"`
	Status         string    `json:"status"`
	DocumentPath   string    `json:"document_path"`
	Created        time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
