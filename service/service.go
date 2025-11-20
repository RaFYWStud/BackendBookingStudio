package service

import "github.com/RaFYWStud/BackendBookingStudio/contract"

func New(repo *contract.Repository) *contract.Service {
    return &contract.Service{
        Auth:          ImplAuthService(repo.Auth),
        Studio:        ImplStudioService(repo.Studio),
        Booking:       ImplBookingService(repo.Booking, repo.Studio),
        PaymentMethod: ImplPaymentMethodService(repo.PaymentMethod), 
        Payment:       ImplPaymentService(repo.Payment, repo.Booking, repo.PaymentMethod),
    }
}