package service

import "github.com/unsrat-it-community/back-end-e-voting-2025/contract"

func New(repo *contract.Repository) *contract.Service {
	return &contract.Service{
		Auth: ImplAuthService(repo.Auth),
		Studio: ImplStudioService(repo.Studio),
		Booking: ImplBookingService(repo.Booking, repo.Studio),
		// Add other services here as needed
	}
}
