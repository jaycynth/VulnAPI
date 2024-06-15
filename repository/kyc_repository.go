package repository

import "github.com/security-testing-api/model"

type KYCRepository interface {
	Save(kyc model.KYC) (model.KYC, error)
}
