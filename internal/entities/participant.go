package entities

import "time"

type Participant struct {
	Timestamp

	ID             HexID
	RelatedUser    User
	SecretFriendID HexID
	JoinedAt       time.Time
	IsReady        bool
}

func NewParticipant(secretFriendID HexID, relatedUser User) Participant {
	return Participant{SecretFriendID: secretFriendID, RelatedUser: relatedUser}
}
