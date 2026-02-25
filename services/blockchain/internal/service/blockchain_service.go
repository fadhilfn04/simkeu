package service

import (
	"simkeu/service-blockchain/internal/repository"
)

type BlockchainService struct {
	Repo *repository.BlockchainRepository
}
