package service

import (
	"simkeu/service-debitur/internal/repository"
)

type DebiturService struct {
	Repo *repository.DebiturRepository
}

func (s *DebiturService) GetByID(id string) (map[string]interface{}, error) {
	return s.Repo.FindByID(id)
}

func (s *DebiturService) Create(id int, name string) error {
	return s.Repo.Create(id, name)
}
