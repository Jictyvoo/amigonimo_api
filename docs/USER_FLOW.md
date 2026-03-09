# User Flow

This document details the flow of using the API to manage a Secret Friend (Secret Santa) game.

1. **Create Group**: A user creates a new Secret Friend group/event.
2. **Invite Participants**: The creator (or authorized members) invites other people to join the group.
3. **Accept Invitation**: The invited persons accept the invitation to join the group.
4. **Manage Preferences**:
    * **Wishlist**: Participants add items they would like to receive (up to a maximum capacity).
    * **Denylist**: Participants add valid peers they *cannot* or *do not want to* be paired with (up to a maximum
      capacity). Users cannot deny themselves or non-participants.
    * **Finalize**: Participants explicitly mark their preferences as finished.
5. **The Draw**:
    * The person in charge requests the draw to be performed.
    * The system calculates the pairings, respecting denylists.
6. **Results**: The results of the draw are effectively sent/revealed to every participant (e.g., via email or creating
   a viewable result for them).

```mermaid
flowchart TD
    A([fa:fa-users Create Group]) --> B([fa:fa-envelope Invite Participants])
    B --> C([Participants Receive Invitation])
    C --> D{Accept Invitation?}
    D -- No --> C
    D -- Yes --> E([fa:fa-user-plus Join Group])
    E --> F([fa:fa-list Manage Preferences])
    F --> F1([fa:fa-gift Add Wishlist])
    F --> F2([fa:fa-ban Add Denylist])
    F1 --> G([fa:fa-check Mark Preferences as Finished])
    F2 --> G
    G --> H([fa:fa-hourglass Waiting Phase])
    H --> I([fa:fa-user-tie Organizer Reviews Status])
    I --> J{All Participants Finished?}
    J -- Yes --> K([fa:fa-play Execute Draw])
    J -- No --> L{Proceed Anyway?}
    L -- Wait --> H
    L -- Execute --> K
    K --> M([fa:fa-cogs System Calculates Pairings])
    M --> N{Denylist Constraints Met?}
    N -- No --> M
    N -- Yes --> O([fa:fa-lock Finalize Pairings])
    O --> P([fa:fa-bullhorn Reveal Results])
    P --> Q([fa:fa-eye View Results / Receive Email])
```
