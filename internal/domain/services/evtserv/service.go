package evtserv

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type Service struct {
	repository SecretFriendRepository
	user       entities.User
}

func NewService(
	repository SecretFriendRepository,
	user entities.User,
) *Service {
	return &Service{
		repository: repository,
		user:       user,
	}
}

func (s *Service) CreateSecretFriend(
	name string,
	datetime *time.Time,
	location string,
	maxDenyListSize int,
) (entities.SecretFriend, error) {
	inviteCode := uuid.New().String()[:8]

	secretFriend := entities.SecretFriend{
		ID:         entities.HexID(uuid.New()),
		Name:       name,
		Datetime:   datetime,
		Location:   location,
		OwnerID:    s.user.ID,
		InviteCode: inviteCode,
		InviteLink: fmt.Sprintf("/invite/%s", inviteCode),
		Status:     entities.StatusDraft,
	}
	secretFriend.CreatedAt = time.Now()
	secretFriend.UpdatedAt = secretFriend.CreatedAt

	if err := s.repository.CreateSecretFriend(&secretFriend); err != nil {
		return entities.SecretFriend{}, fmt.Errorf("create secret friend: %w", err)
	}

	return secretFriend, nil
}

func (s *Service) GetSecretFriend(id entities.HexID) (entities.SecretFriend, error) {
	secretFriend, err := s.repository.GetSecretFriendByID(id)
	if err != nil {
		return entities.SecretFriend{}, fmt.Errorf("get secret friend: %w", err)
	}

	return secretFriend, nil
}

func (s *Service) UpdateSecretFriend(
	id entities.HexID,
	name *string,
	datetime *time.Time,
	location *string,
) error {
	secretFriend, err := s.repository.GetSecretFriendByID(id)
	if err != nil {
		return fmt.Errorf("get secret friend: %w", err)
	}

	if name != nil {
		secretFriend.Name = *name
	}
	if datetime != nil {
		secretFriend.Datetime = datetime
	}
	if location != nil {
		secretFriend.Location = *location
	}
	secretFriend.UpdatedAt = time.Now()

	if err := s.repository.UpdateSecretFriend(&secretFriend); err != nil {
		return fmt.Errorf("update secret friend: %w", err)
	}

	return nil
}

func (s *Service) DrawSecretFriend(id entities.HexID) (int, error) {
	secretFriend, err := s.repository.GetSecretFriendByID(id)
	if err != nil {
		return 0, fmt.Errorf("get secret friend: %w", err)
	}

	if secretFriend.Status == entities.StatusDrawn || secretFriend.Status == entities.StatusClosed {
		return 0, fmt.Errorf("secret friend already drawn")
	}

	// if err := s.drawService.ExecuteDraw(uuid.UUID(id)); err != nil {
	// 	return 0, fmt.Errorf("execute draw: %w", err)
	// }

	secretFriend.Status = entities.StatusDrawn
	secretFriend.UpdatedAt = time.Now()

	if err = s.repository.UpdateSecretFriend(&secretFriend); err != nil {
		return 0, fmt.Errorf("update secret friend status: %w", err)
	}

	count, err := s.repository.GetParticipantsCount(id)
	if err != nil {
		return 0, fmt.Errorf("get participants count: %w", err)
	}

	return count, nil
}

func (s *Service) GetDrawResultForUser(
	secretFriendID entities.HexID,
) (entities.DrawResultItem, error) {
	result, err := s.repository.GetDrawResultForUser(secretFriendID, s.user.ID)
	if err != nil {
		return entities.DrawResultItem{}, fmt.Errorf("get draw result: %w", err)
	}

	return result, nil
}
