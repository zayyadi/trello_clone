# Trello Clone Backend (Go)

This is a sophisticated backend for a Trello-like task management application, built with Go. It provides a robust foundation with core features, a layered architecture, and database integration.

## Features
-   **User Authentication:** Secure registration and login using JWT.
-   **Board Management:**
    -   CRUD operations for boards.
    -   Each board has an owner.
    -   Board membership: Users can be added to or removed from boards.
-   **List Management:**
    -   CRUD operations for lists within a board.
    -   List reordering within a board.
-   **Card Management:**
    -   CRUD operations for cards within a list.
    -   Card details: title, description, due date, assigned user.
    -   Card reordering within a list.
    -   Moving cards between lists on the same board.
-   **Layered Architecture:** Handlers, Services, Repositories for clear separation of concerns.
-   **Database:** PostgreSQL with GORM, including automatic schema migrations.
-   **Configuration:** Environment variable based configuration (`.env` file).
-   **API Design:** RESTful with DTOs for request validation and response formatting.
-   **Authorization:**
    -   Board access controlled by ownership and membership.
    -   Operations on lists and cards require appropriate board access.
-   **Docker Support:** Includes a `Dockerfile` for easy containerization.

## Setup

1.  **PostgreSQL:**
    *   Ensure you have a PostgreSQL server running.
    *   Create a database (e.g., `trello_clone_db`).

2.  **Environment Variables:**
    *   Copy `.env.example` to `.env`: `cp .env.example .env`
    *   Edit `.env` and fill in your database credentials and a strong `JWT_SECRET_KEY`.

3.  **Go Dependencies:**
    ```bash
    go mod tidy
    ```

4.  **Run the Application:**
    ```bash
    go run main.go
    ```
    The server will start, typically on port 8080 (configurable in `.env`).

## API Endpoints

Base URL: `http://localhost:{SERVER_PORT}` (e.g., `http://localhost:8080`)

All `/api/*` routes require JWT authentication via an `Authorization: Bearer <token>` header.

### Authentication (`/auth`)
-   `POST /auth/register` - Register a new user.
    -   Body: `{"username": "user", "email": "user@example.com", "password": "password123"}`
-   `POST /auth/login` - Login an existing user.
    -   Body: `{"email": "user@example.com", "password": "password123"}`

### Boards (`/api/boards`)
-   `POST /api/boards` - Create a new board.
    -   Body: `{"name": "My Project Board", "description": "Board for project X"}`
-   `GET /api/boards` - Get all boards the user owns or is a member of.
-   `GET /api/boards/:boardID` - Get a specific board by ID (if user has access).
-   `PUT /api/boards/:boardID` - Update a board (only owner).
    -   Body: `{"name": "Updated Project Board", "description": "New description"}`
-   `DELETE /api/boards/:boardID` - Delete a board (only owner).

### Board Members (`/api/boards/:boardID/members`)
-   `POST /api/boards/:boardID/members` - Add a member to a board (only owner).
    -   Body: `{"email": "member@example.com"}` (or `{"userID": <member_user_id>}`)
-   `GET /api/boards/:boardID/members` - List members of a board.
-   `DELETE /api/boards/:boardID/members/:userID` - Remove a member from a board (only owner).

### Lists (`/api/boards/:boardID/lists` and `/api/lists/:listID`)
-   `POST /api/boards/:boardID/lists` - Create a new list on a board.
    -   Body: `{"name": "To Do", "position": 1}` (position is for ordering)
-   `GET /api/boards/:boardID/lists` - Get all lists for a specific board.
-   `PUT /api/lists/:listID` - Update a list (name, position).
    -   Body: `{"name": "Doing", "position": 2}`
-   `DELETE /api/lists/:listID` - Delete a list.

### Cards (`/api/lists/:listID/cards` and `/api/cards/:cardID`)
-   `POST /api/lists/:listID/cards` - Create a new card in a list.
    -   Body: `{"title": "Setup project", "description": "Initial setup tasks", "position": 1, "dueDate": "2024-12-31T23:59:59Z", "assignedUserID": null}`
-   `GET /api/lists/:listID/cards` - Get all cards for a specific list.
-   `GET /api/cards/:cardID` - Get a specific card by ID.
-   `PUT /api/cards/:cardID` - Update a card.
    -   Body: (any fields from create, e.g., `{"title": "Updated Task", "description": "...", "dueDate": "..."}`)
-   `DELETE /api/cards/:cardID` - Delete a card.
-   `PATCH /api/cards/:cardID/move` - Move a card to a different list and/or position.
    -   Body: `{"targetListID": <new_list_id>, "newPosition": <new_position_in_target_list>}`

## Further Improvements
-   **Real-time Updates:** Implement WebSockets (e.g., using Gorilla WebSocket or Nhooyr WebSocket) for live updates across connected clients when boards, lists, or cards are modified.
-   **Advanced Authorization & Roles:** Introduce more granular roles for board members (e.g., admin, editor, viewer) with specific permissions.
-   **Comprehensive Input Validation:** Enhance validation rules for all DTOs.
-   **Testing:** Add unit tests for services and repositories, and integration tests for handlers.
-   **File Attachments:** Allow users to attach files to cards (e.g., storing in S3 or local filesystem).
-   **Notifications:** Implement in-app or email notifications for mentions, assignments, due date reminders.
-   **Activity Logging:** Track user actions (e.g., card creation, moves, comments) for an audit trail.
-   **Search Functionality:** Implement search across boards, lists, and cards.
-   **Card Details:** Add features like labels/tags, checklists.
-   **Soft Deletes & Archiving:** Implement soft deletes for data and an archiving mechanism for boards, lists, and cards.
-   **Advanced Reordering:** Use more robust reordering algorithms (e.g., fractional indexing) to avoid re-indexing large numbers of items.
-   **API Versioning.**
-   **Rate Limiting & Security Headers.**
-   **OAuth2 Integration:** Allow login with Google, GitHub, etc.