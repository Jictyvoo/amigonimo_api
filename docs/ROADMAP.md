# Roadmap

## Testing & Quality Assurance

- [ ] **Unit Tests**: Create more unit tests for each usecase and service to improve coverage.
- [ ] **Integration Tests**: Integrate a custom mock server to allow testing of mailer requests/email sending in
  integration tests.

## Core Features & Services

- [ ] **Draw Service**: Finalize the implementation of the `drawserv` service to ensure correct pairing logic.
- [ ] **Mailer Service**: Implement the mailer service for sending notifications and results.

## Flow & Logic Refactoring

- [ ] **Invite Flow**: Complete revisit of the "invite someone to secret-friend" flow:
    - Review invite code mechanisms.
    - Define access controls (public vs. specific people).
- [ ] **Denylist Control**: Add a missing feature to control/limit the maximum size of a user's denylist.
- [ ] **Wishlist Control**: Add an internal limit for the maximum wishlist size (even if not exposed
  configuration-wise).

## Data Structure Refactoring

- [ ] **User Struct**: Refactor the `User` struct by removing `FullName`.
- [ ] **UserProfile**: Move `FullName` to a new `UserProfile` struct/table (keeping column data intact/migrated).
