package service

import (
	"github.com/security-testing-api/model"
	"github.com/security-testing-api/repository"
)

type KYCServiceImpl struct {
	kycRepo repository.KYCRepository
}

func NewKYCServiceImpl(kycRepo repository.KYCRepository) KYCService {
	return &KYCServiceImpl{kycRepo: kycRepo}
}

func (s *KYCServiceImpl) SaveKYC(kyc *model.KYC) (*model.KYC, error) {
	return s.kycRepo.Save(kyc)
}
