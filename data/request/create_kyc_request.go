package request

type SaveKYCRequest struct {
	FullName       string `json:"full_name" binding:"required"`
	DocumentType   string `json:"document_type" binding:"required"`
	DocumentNumber string `json:"document_number" binding:"required"`
	IssueDateStr   string `json:"issue_date" binding:"required"`
	ExpiryDateStr  string `json:"expiry_date" binding:"required"`
	Status         string `json:"status" binding:"required"`
}
