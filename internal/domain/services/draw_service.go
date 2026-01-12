package services

import (
	"github.com/google/uuid"
)

type DrawService struct{}

func NewDrawService() *DrawService {
	return &DrawService{}
}

// ExecuteDraw performs the secret friend draw algorithm.
func (s *DrawService) ExecuteDraw(secretFriendID uuid.UUID) error {
	return nil
}
