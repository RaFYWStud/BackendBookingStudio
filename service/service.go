package service

import "github.com/RaFYWStud/BackendBookingStudio/contract"

func New(repo *contract.Repository) *contract.Service {
    emailService := ImplEmailService()
    
    return &contract.Service{
        Auth:          ImplAuthService(repo.Auth),
        Studio:        ImplStudioService(repo.Studio),
        Booking:       ImplBookingService(repo.Booking, repo.Studio, emailService),
        PaymentMethod: ImplPaymentMethodService(repo.PaymentMethod), 
        Payment:       ImplPaymentService(repo.Payment, repo.Booking, repo.PaymentMethod, emailService),
        Email:         emailService,
    }
}