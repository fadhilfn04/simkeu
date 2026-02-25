package service

import (
	"simkeu/service-master/internal/repository"
)

type MasterService struct {
	Repo *repository.MasterRepository
}
