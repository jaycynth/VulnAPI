package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/security-testing-api/data/request"
	"github.com/security-testing-api/data/response"
	"github.com/security-testing-api/helper"
	"github.com/security-testing-api/model"
	"github.com/security-testing-api/service"
)

type KYCController struct {
	kycService service.KYCService
}

func NewKYCController(service service.KYCService) *KYCController {
	return &KYCController{kycService: service}
}

func (controller *KYCController) Save(ctx *gin.Context) {

	userId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userID, err := strconv.Atoi(userId.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:   http.StatusBadRequest,
			Status: "Invalid user ID",
			Error:  err.Error(),
		})
		return
	}

	err = ctx.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:   http.StatusBadRequest,
			Status: "Could not parse multipart form",
			Error:  err.Error(),
		})
		return
	}

	file, header, err := ctx.Request.FormFile("document")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:   http.StatusBadRequest,
			Status: "Document file is required",
			Error:  err.Error(),
		})
		return
	}
	defer file.Close()

	documentPath := fmt.Sprintf("uploads/%s", header.Filename)
	if err := ctx.SaveUploadedFile(header, documentPath); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:   http.StatusInternalServerError,
			Status: "Failed to upload file",
			Error:  err.Error(),
		})
		return
	}

	kycData := ctx.Request.FormValue("kyc_data")
	var kycRequest request.SaveKYCRequest
	if err := json.Unmarshal([]byte(kycData), &kycRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:   http.StatusBadRequest,
			Status: "Invalid JSON Data",
			Error:  err.Error(),
		})
		return
	}

	kyc := &model.KYC{
		UserID:         userID,
		DocumentType:   kycRequest.DocumentType,
		DocumentNumber: kycRequest.DocumentNumber,
		ExpiryDate:     time.Now(),
		IssueDate:      time.Now(),
		Status:         kycRequest.Status,
		DocumentPath:   documentPath,
	}

	createdKYC, err := controller.kycService.SaveKYC(kyc)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code:   http.StatusInternalServerError,
			Status: "Error saving kyc to database",
			Error:  err.Error(),
		})
		ctx.JSON(http.StatusInternalServerError, helper.ErrorFormatter(err))
		return
	}

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Created",
		Data:   createdKYC,
	}
	ctx.JSON(http.StatusOK, webResponse)
}
