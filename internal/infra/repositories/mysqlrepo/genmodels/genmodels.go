package genmodels

import "github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"

type (
	AuthToken    = dbgen.AuthToken
	Denylist     = dbgen.Denylist
	DrawResult   = dbgen.DrawResult
	Participant  = dbgen.Participant
	SecretFriend = dbgen.SecretFriend
	User         = dbgen.User
	UserProfile  = dbgen.UserProfile
	WishlistItem = dbgen.WishlistItem
)
