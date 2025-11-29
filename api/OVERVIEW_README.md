# Overview of the Fuego‑Go API Server

This document gives a concise, high‑level view of the REST API that powers the Secret Friend (Amigo Secreto)
application.  
The implementation is based on the **Fuego** framework (Go), which simplifies routing, validation, and OpenAPI
generation.

> **Authentication**  
> All protected endpoints require a bearer token: `Authorization: Bearer <token>`.  
> Tokens are issued on login (`POST /auth/login`) and are linked to a specific user by email.

---

## 1. Core Endpoints

| Method    | Path                                | Purpose                                                    | Notes                                                          |
|-----------|-------------------------------------|------------------------------------------------------------|----------------------------------------------------------------|
| **GET**   | `/dashboard`                        | Retrieve the user’s active events (both owned and joined). | Returns two lists: `activeCreated` and `activeParticipant`.    |
| **GET**   | `/invites/{code}`                   | Fetch the secret‑friend event information by invite code.  | Returns the event’s `secretFriendId` and name.                 |
| **PATCH** | `/secret-friends/{id}`              | Update an event’s details (name, datetime, location).      | Only accessible to the event manager.                          |
| **POST**  | `/secret-friends/{id}/draw`         | Execute the draw for the event.                            | Returns the drawn friend’s ID, status, and total participants. |
| **POST**  | `/secret-friends/{id}/participants` | Confirm or cancel participation in an event.               | Requires `confirm: true/false` in the body.                    |
| **GET**   | `/secret-friends/{id}/participants` | List all confirmed participants.                           |
| **GET**   | `/secret-friends/{id}/wishlist`     | Retrieve the current user’s wishlist for the event.        | Only returns the caller’s own wishlist.                        |
| **POST**  | `/auth/login`                       | Login via email (no password).                             | Returns `userId` and JWT `token`.                              |
| **POST**  | `/auth/register`                    | Register a new user.                                       | Optional `inviteCode` field.                                   |

> All operations that involve a secret‑friend event expect the event ID in the path and are scoped to the authenticated
> user.

---

## 2. Drawing Mechanics

When the draw is executed (`POST /secret-friends/{id}/draw`), the system:

1. Randomly assigns each participant a secret friend.
2. Persists the result in the database.
3. Sends an email to each participant with:
    * Their assigned friend’s name.
    * Their own wishlist.
    * Confirmation of the event location, date, and time.
    * This is defined in the README and forms part of the business logic. [1]

---

## 4. OpenAPI Documentation

Fuego automatically generates an OpenAPI spec based on the `option` decorators.  
You can view the live spec by navigating to `/openapi.json` (or `/swagger-ui` if enabled) after the server starts.
