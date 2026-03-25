package integration

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/genmodels"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/participantsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/secretfriendsctrl"
	"github.com/jictyvoo/amigonimo_api/test/integration/stdrunners"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
)

func TestSecretFriendMountedRoutes(t *testing.T) {
	engine := NewEngine(t)
	const ownerPassword = "owner-routes-password"

	ownerBuilder := fixtures.NewUser().
		WithEmail("routes-owner@example.com").
		WithFullname("Routes Owner").
		WithPassword(ownerPassword)
	owner := ownerBuilder.Build()
	ownerProfile := ownerBuilder.BuildProfile()

	targetBuilder := fixtures.NewUser().
		WithEmail("routes-target@example.com").
		WithFullname("Target Person").
		WithPassword(ownerPassword)
	targetUser := targetBuilder.Build()
	targetProfile := targetBuilder.BuildProfile()

	otherBuilder := fixtures.NewUser().
		WithEmail("routes-other@example.com").
		WithFullname("Other Person").
		WithPassword(ownerPassword)
	otherUser := otherBuilder.Build()
	otherProfile := otherBuilder.BuildProfile()

	openEvent := fixtures.NewSecretFriend().
		WithOwner(owner).
		WithName("Open Routes Event").
		Build()
	openEvent.InviteCode = "openrt01"
	openEvent.Location = sql.NullString{String: "Original Place", Valid: true}
	openEvent.Datetime = sql.NullTime{Time: time.Now().Add(48 * time.Hour), Valid: true}

	drawnEvent := fixtures.NewSecretFriend().
		WithOwner(owner).
		WithName("Drawn Routes Event").
		Build()
	drawnEvent.InviteCode = "drawrt01"
	drawnEvent.Location = sql.NullString{String: "Drawn Place", Valid: true}
	drawnEvent.Datetime = sql.NullTime{Time: time.Now().Add(72 * time.Hour), Valid: true}
	drawnEvent.Status = string(entities.StatusDrawn)

	openOwnerParticipant := fixtures.NewParticipant().
		WithUser(owner).
		WithSecretFriend(openEvent).
		Build()
	openTargetParticipant := fixtures.NewParticipant().
		WithUser(targetUser).
		WithSecretFriend(openEvent).
		Build()
	openOtherParticipant := fixtures.NewParticipant().
		WithUser(otherUser).
		WithSecretFriend(openEvent).
		Build()

	drawnOwnerParticipant := fixtures.NewParticipant().
		WithUser(owner).
		WithSecretFriend(drawnEvent).
		Build()
	drawnTargetParticipant := fixtures.NewParticipant().
		WithUser(targetUser).
		WithSecretFriend(drawnEvent).
		Build()
	drawnOtherParticipant := fixtures.NewParticipant().
		WithUser(otherUser).
		WithSecretFriend(drawnEvent).
		Build()

	drawResultID := uuid.Must(uuid.NewV7())
	drawResult := &genmodels.DrawResult{
		ID:                    drawResultID[:],
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		GiverParticipantID:    drawnOwnerParticipant.ID,
		ReceiverParticipantID: drawnTargetParticipant.ID,
		SecretFriendID:        drawnEvent.ID,
	}

	wishlistID := uuid.Must(uuid.NewV7())
	wishlistItem := &genmodels.WishlistItem{
		ID:            wishlistID[:],
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Label:         "Coffee voucher",
		Comments:      sql.NullString{String: "Medium roast", Valid: true},
		ParticipantID: drawnTargetParticipant.ID,
	}

	if err := engine.Seed(
		owner,
		ownerProfile,
		targetUser,
		targetProfile,
		otherUser,
		otherProfile,
		openEvent,
		drawnEvent,
		openOwnerParticipant,
		openTargetParticipant,
		openOtherParticipant,
		drawnOwnerParticipant,
		drawnTargetParticipant,
		drawnOtherParticipant,
		drawResult,
		wishlistItem,
	); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	openEventID, _ := entities.NewHexIDFromBytes(openEvent.ID)
	drawnEventID, _ := entities.NewHexIDFromBytes(drawnEvent.ID)
	ownerID, _ := entities.NewHexIDFromBytes(owner.ID)
	openOwnerParticipantID, _ := entities.NewHexIDFromBytes(openOwnerParticipant.ID)
	openTargetParticipantID, _ := entities.NewHexIDFromBytes(openTargetParticipant.ID)
	openOtherParticipantID, _ := entities.NewHexIDFromBytes(openOtherParticipant.ID)
	targetUserID, _ := entities.NewHexIDFromBytes(targetUser.ID)
	otherUserID, _ := entities.NewHexIDFromBytes(otherUser.ID)

	updatedName := "Updated Open Routes Event"
	updatedLocation := "Updated Place"

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			stdrunners.LoginRunner(engine.BaseURL(), owner.Email, ownerPassword),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/", struct{}{}),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.ExpectStatus(http.StatusInternalServerError),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(
					http.MethodGet,
					"/secret-friends/invites/description/{code}",
					struct{}{},
				),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("code", openEvent.InviteCode),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					secretfriendsctrl.InviteInfoResponse{
						SecretFriendID: openEventID.String(),
						Name:           openEvent.Name,
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}", struct{}{}),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", openEventID),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					secretfriendsctrl.GetSecretFriendResponse{
						ID:                openEventID.String(),
						Name:              openEvent.Name,
						Location:          openEvent.Location.String,
						OwnerID:           ownerID.String(),
						ParticipantsCount: 3,
						Status:            openEvent.Status,
					},
					func(expected, actual *secretfriendsctrl.GetSecretFriendResponse) error {
						expected.Datetime = actual.Datetime
						return nil
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(
					http.MethodPatch,
					"/secret-friends/{id}",
					secretfriendsctrl.UpdateSecretFriendRequest{
						Name:     updatedName,
						Location: updatedLocation,
					},
				),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", openEventID),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					secretfriendsctrl.UpdateSecretFriendResponse{
						Success: true,
						Message: "secret friend updated successfully",
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}", struct{}{}),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", openEventID),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					secretfriendsctrl.GetSecretFriendResponse{
						ID:                openEventID.String(),
						Name:              updatedName,
						Location:          updatedLocation,
						OwnerID:           ownerID.String(),
						ParticipantsCount: 3,
						Status:            openEvent.Status,
					},
					func(expected, actual *secretfriendsctrl.GetSecretFriendResponse) error {
						expected.Datetime = actual.Datetime
						return nil
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(
					http.MethodGet,
					"/secret-friends/{id}/participants/",
					struct{}{},
				),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", openEventID),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					[]participantsctrl.ParticipantResponse{
						{
							ParticipantID: openOwnerParticipantID.String(),
							UserID:        ownerID.String(),
							Fullname:      ownerProfile.Fullname.String,
						},
						{
							ParticipantID: openTargetParticipantID.String(),
							UserID:        targetUserID.String(),
							Fullname:      targetProfile.Fullname.String,
						},
						{
							ParticipantID: openOtherParticipantID.String(),
							UserID:        otherUserID.String(),
							Fullname:      otherProfile.Fullname.String,
						},
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/draw-result", struct{}{}),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", drawnEventID),
				netoche.ExpectStatus(http.StatusInternalServerError),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}", struct{}{}),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", "not-a-uuid"),
				netoche.ExpectStatus(http.StatusBadRequest),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestDrawSecretFriendRoute(t *testing.T) {
	engine := NewEngine(t)
	const ownerPassword = "draw-routes-password"

	ownerBuilder := fixtures.NewUser().
		WithEmail("draw-owner@example.com").
		WithPassword(ownerPassword)
	owner := ownerBuilder.Build()
	ownerProfile := ownerBuilder.BuildProfile()

	userTwoBuilder := fixtures.NewUser().
		WithEmail("draw-user-two@example.com").
		WithPassword(ownerPassword)
	userTwo := userTwoBuilder.Build()
	userTwoProfile := userTwoBuilder.BuildProfile()

	userThreeBuilder := fixtures.NewUser().
		WithEmail("draw-user-three@example.com").
		WithPassword(ownerPassword)
	userThree := userThreeBuilder.Build()
	userThreeProfile := userThreeBuilder.BuildProfile()

	drawEvent := fixtures.NewSecretFriend().
		WithOwner(owner).
		WithName("Draw Execution Event").
		Build()
	drawEvent.Status = string(entities.StatusOpen)
	drawEvent.Datetime = sql.NullTime{Time: time.Now().Add(24 * time.Hour), Valid: true}

	ownerParticipant := fixtures.NewParticipant().
		WithUser(owner).
		WithSecretFriend(drawEvent).
		Build()
	userTwoParticipant := fixtures.NewParticipant().
		WithUser(userTwo).
		WithSecretFriend(drawEvent).
		Build()
	userThreeParticipant := fixtures.NewParticipant().
		WithUser(userThree).
		WithSecretFriend(drawEvent).
		Build()

	if err := engine.Seed(
		owner,
		ownerProfile,
		userTwo,
		userTwoProfile,
		userThree,
		userThreeProfile,
		drawEvent,
		ownerParticipant,
		userTwoParticipant,
		userThreeParticipant,
	); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	drawEventID, _ := entities.NewHexIDFromBytes(drawEvent.ID)

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			stdrunners.LoginRunner(engine.BaseURL(), owner.Email, ownerPassword),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodPost, "/secret-friends/{id}/draw", struct{}{}),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", drawEventID),
				netoche.ExpectStatus(http.StatusInternalServerError),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
