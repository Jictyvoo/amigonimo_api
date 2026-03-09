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
- [ ] **Participant Validation (ID Injection)**: Ensure users can only add to their wishlist or denylist if they are
  confirmed participants of the group. Prevent ID injection where a user could modify another user's lists.
- [ ] **Denylist Control**: Add a missing feature to control/limit the maximum size of a user's denylist. Prevent users
  from adding themselves or non-participants to their denylist.
- [ ] **Wishlist Control**: Add an internal limit for the maximum wishlist size (even if not exposed
  configuration-wise).
- [ ] **Wait Phase / Status**: Implement a feature for users to explicitly mark their preferences (wishlist/denylist) as
  finished, allowing the organizer to see who is ready.

## Data Structure Refactoring

- [ ] **User Struct**: Refactor the `User` struct by removing `FullName`.
- [ ] **UserProfile**: Move `FullName` to a new `UserProfile` struct/table (keeping column data intact/migrated).
