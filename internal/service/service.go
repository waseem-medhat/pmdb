package service

import (
	"github.com/wipdev-tech/pmdb/internal/database"
)

type Service struct {
	DB *database.Queries
}

func New() *Service {
	return &Service{}
}
