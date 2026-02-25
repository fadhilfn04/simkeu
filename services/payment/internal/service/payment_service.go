package service

import (
	"simkeu/service-payment/internal/repository"
)

type PaymentService struct {
	Repo *repository.PaymentRepository
}
