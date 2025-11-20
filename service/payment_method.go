package service

import (
	"fmt"

	"github.com/RaFYWStud/BackendBookingStudio/config/pkg/errs"
	"github.com/RaFYWStud/BackendBookingStudio/contract"
	"github.com/RaFYWStud/BackendBookingStudio/database"
	"github.com/RaFYWStud/BackendBookingStudio/dto"
	"gorm.io/gorm"
)

type paymentMethodService struct {
    paymentMethodRepo contract.PaymentMethodRepository
}

func ImplPaymentMethodService(paymentMethodRepo contract.PaymentMethodRepository) contract.PaymentMethodService {
    return &paymentMethodService{
        paymentMethodRepo: paymentMethodRepo,
    }
}

// GetActivePaymentMethods - Public: Customer lihat metode pembayaran aktif
func (s *paymentMethodService) GetActivePaymentMethods() (*dto.PaymentMethodListResponse, error) {
    methods, err := s.paymentMethodRepo.FindActive()
    if err != nil {
        return nil, errs.InternalServerError("failed to fetch payment methods")
    }

    methodDataList := make([]dto.PaymentMethodData, len(methods))
    for i, method := range methods {
        methodDataList[i] = s.mapPaymentMethodToDTO(&method)
    }

    return &dto.PaymentMethodListResponse{
        Success: true,
        Data:    methodDataList,
    }, nil
}

// GetAllPaymentMethods - Admin: Get all payment methods
func (s *paymentMethodService) GetAllPaymentMethods() (*dto.PaymentMethodListResponse, error) {
    methods, err := s.paymentMethodRepo.FindAll()
    if err != nil {
        return nil, errs.InternalServerError("failed to fetch payment methods")
    }

    methodDataList := make([]dto.PaymentMethodData, len(methods))
    for i, method := range methods {
        methodDataList[i] = s.mapPaymentMethodToDTO(&method)
    }

    return &dto.PaymentMethodListResponse{
        Success: true,
        Data:    methodDataList,
    }, nil
}

// GetPaymentMethodByID - Admin: Get payment method detail
func (s *paymentMethodService) GetPaymentMethodByID(id int) (*dto.PaymentMethodResponse, error) {
    method, err := s.paymentMethodRepo.FindByID(id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("payment method not found")
        }
        return nil, errs.InternalServerError("failed to fetch payment method")
    }

    return &dto.PaymentMethodResponse{
        Success: true,
        Data:    s.mapPaymentMethodToDTO(method),
    }, nil
}

// CreatePaymentMethod - Admin: Create new payment method
func (s *paymentMethodService) CreatePaymentMethod(req dto.CreatePaymentMethodRequest) (*dto.PaymentMethodResponse, error) {
    method := &database.PaymentMethod{
        Name:          req.Name,
        BankName:      req.BankName,
        AccountNumber: req.AccountNumber,
        AccountName:   req.AccountName,
        IsActive:      true,
    }

    if err := s.paymentMethodRepo.Create(method); err != nil {
        return nil, errs.InternalServerError("failed to create payment method")
    }

    return &dto.PaymentMethodResponse{
        Success: true,
        Message: "Payment method created successfully",
        Data:    s.mapPaymentMethodToDTO(method),
    }, nil
}

// UpdatePaymentMethod - Admin: Update payment method
func (s *paymentMethodService) UpdatePaymentMethod(id int, req dto.UpdatePaymentMethodRequest) (*dto.PaymentMethodResponse, error) {
    method, err := s.paymentMethodRepo.FindByID(id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("payment method not found")
        }
        return nil, errs.InternalServerError("failed to fetch payment method")
    }

    // Update fields if provided
    if req.Name != nil {
        method.Name = *req.Name
    }
    if req.BankName != nil {
        method.BankName = *req.BankName
    }
    if req.AccountNumber != nil {
        method.AccountNumber = *req.AccountNumber
    }
    if req.AccountName != nil {
        method.AccountName = *req.AccountName
    }
    if req.IsActive != nil {
        method.IsActive = *req.IsActive
    }

    if err := s.paymentMethodRepo.Update(method); err != nil {
        return nil, errs.InternalServerError("failed to update payment method")
    }

    return &dto.PaymentMethodResponse{
        Success: true,
        Message: "Payment method updated successfully",
        Data:    s.mapPaymentMethodToDTO(method),
    }, nil
}

// DeletePaymentMethod - Admin: Delete payment method
func (s *paymentMethodService) DeletePaymentMethod(id int) (*dto.DeletePaymentMethodResponse, error) {
    _, err := s.paymentMethodRepo.FindByID(id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errs.NotFound("payment method not found")
        }
        return nil, errs.InternalServerError("failed to fetch payment method")
    }

    if err := s.paymentMethodRepo.Delete(id); err != nil {
        return nil, errs.InternalServerError("failed to delete payment method")
    }

    return &dto.DeletePaymentMethodResponse{
        Success: true,
        Message: fmt.Sprintf("Payment method with ID %d deleted successfully", id),
    }, nil
}

// Helper: Map payment method to DTO
func (s *paymentMethodService) mapPaymentMethodToDTO(method *database.PaymentMethod) dto.PaymentMethodData {
    return dto.PaymentMethodData{
        ID:            method.ID,
        Name:          method.Name,
        BankName:      method.BankName,
        AccountNumber: method.AccountNumber,
        AccountName:   method.AccountName,
        IsActive:      method.IsActive,
        CreatedAt:     method.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt:     method.UpdatedAt.Format("2006-01-02 15:04:05"),
    }
}