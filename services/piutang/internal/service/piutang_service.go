package service

import (
	"simkeu/service-piutang/internal/repository"
)

type PiutangService struct {
	Repo *repository.PiutangRepository
}
