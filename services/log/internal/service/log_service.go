package service

import (
	"simkeu/service-log/internal/repository"
)

type LogService struct {
	Repo *repository.LogRepository
}
