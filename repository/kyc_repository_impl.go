package repository

import (
	"github.com/security-testing-api/model"
	"gorm.io/gorm"
)

type KYCRepositoryImpl struct {
	Db *gorm.DB
}

func NewKYCRepositoryImpl(Db *gorm.DB) KYCRepository {
	return &KYCRepositoryImpl{Db: Db}
}

func (r *KYCRepositoryImpl) Save(kyc model.KYC) (model.KYC, error) {
	err := r.Db.Create(&kyc).Error
	return kyc, err
}
