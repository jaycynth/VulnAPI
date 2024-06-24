package service

import "github.com/security-testing-api/model"

type KYCService interface {
	SaveKYC(kyc *model.KYC) (*model.KYC, error)
}
