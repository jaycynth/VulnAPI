package repository

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/security-testing-api/model"
	"github.com/stretchr/testify/require"
)

func (s *Suite) Test_Kyc_Save() {
	kyc := &model.KYC{
		UserID:         1,
		DocumentType:   "png",
		DocumentNumber: "827662772",
		IssueDate:      time.Now(),
		ExpiryDate:     time.Now(),
		Status:         "New",
		DocumentPath:   "uploads/like.png",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`INSERT INTO kycs (.+) VALUES (.+)`).
		WithArgs(kyc.UserID, kyc.DocumentNumber).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	savedKYC, err := s.kYCRepository.Save(kyc)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), savedKYC)
	require.Equal(s.T(), kyc.UserID, savedKYC.UserID)
	require.Equal(s.T(), kyc.DocumentNumber, savedKYC.DocumentNumber)
}
