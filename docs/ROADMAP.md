# Roadmap

Items listed here are outstanding. Completed features are documented in [FEATURES.md](FEATURES.md).

## Testing & Quality Assurance

- **Integration Tests – Mailer**: Add a mock HTTP server to cover the webhook mailer backend and SMTP
  interactions in integration tests.

## Flow & Logic

- **Invite Flow**: The invite-code lookup endpoint exists. A full review of access controls (public
  vs. restricted invites) and the invite acceptance flow is still needed.
- **Unmounted Auth Routes**: `PUT /auth/password/forgot` and `GET /auth/verify/:verify_code` are
  implemented but not wired into the router; the forgot-password and email-verification flows are
  therefore incomplete end-to-end.

## Code Structure

- **Entity Boundary Cleanup**: Several DTO-like types (`UserBasic`, `BasicAuthToken`, `VerifyToken`,
  `DeniedUser`, `DrawResult`) still live in `internal/entities`. They should move to action-local
  packages or a dedicated value-object layer.
- **Action-Split Usecases**: `secretfriend`, `participant`, `wishlist`, and `denylist` packages are
  still service-style aggregates. The target shape is one package per application action (matching the
  pattern already established by `drawfriends/execute` and `drawfriends/getresult`).

