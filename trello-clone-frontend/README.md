# Trello Clone - Frontend

## Project Overview

This is the frontend for a Trello-like task management application, built with [Next.js](https://nextjs.org) and [Material UI (MUI)](https://mui.com/). It provides a user interface for interacting with a compatible backend API to manage boards, lists, cards, and user collaboration. This application uses [Redux Toolkit](https://redux-toolkit.js.org/) for state management.

## Features

*   **User Authentication:** Secure user registration and login (JWT-based).
*   **Board Management:**
    *   View a list of user-specific boards.
    *   Create new boards with a name and description.
    *   Edit existing board details.
    *   Delete boards (restricted to board owner).
*   **List Management:**
    *   View lists within a board, ordered by position.
    *   Create new lists within a board.
    *   Edit list names.
    *   Delete lists.
    *   Drag and drop to reorder lists within a board.
*   **Card Management:**
    *   View cards within lists, ordered by position.
    *   Create new cards with a title and optional description.
    *   Edit card details (title, description, due date, status, assignee, supervisor).
    *   Delete cards.
    *   Drag and drop to reorder cards within the same list or move cards between different lists on the same board.
*   **Card Collaboration & Permissions:**
    *   Assign users (board members) to cards.
    *   Add/remove collaborators (board members) to/from cards.
    *   **Permission Model:**
        *   Board owners have full control over their boards, including lists, cards, and managing collaborators.
        *   Card assignees and collaborators can view the card and edit specific fields (e.g., description, status, due date, reassign). Editing title is restricted to the board owner.
        *   Card Visibility:
            *   Board owners see all cards on their boards.
            *   Board members who are assignees or collaborators on a card can see that card.
            *   If a card has no specific assignee and no collaborators, it is only visible to the board owner (strict interpretation). Otherwise, general board members who are not owner/assignee/collaborator will not see these restricted cards.
*   **Comments:** Add and delete text-based comments on cards.
*   **State Management:** Uses Redux Toolkit for a centralized and predictable state container.
*   **API Interaction:** Communicates with a separate backend API for data persistence and business logic.

## Setup and Installation

1.  **Clone the Repository:**
    ```bash
    git clone <repository-url>
    cd trello-clone-frontend
    ```

2.  **Install Dependencies:**
    ```bash
    npm install
    # or
    # yarn install
    ```

3.  **Environment Variables:**
    *   This project requires a backend API to function. Create a `.env.local` file in the root of the `trello-clone-frontend` directory by copying the example file:
        ```bash
        cp .env.example .env.local
        ```
    *   Modify `.env.local` with the URLs for your running backend instance:
        ```
        NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api
        NEXT_PUBLIC_AUTH_BASE_URL=http://localhost:8080/auth
        ```
        (Adjust the port and base path if your backend is running elsewhere).

## Running the Development Server

1.  **Start the Backend:** Ensure the Go backend server (from the `trello/` directory of the main project) is running.

2.  **Start the Frontend Development Server:**
    ```bash
    npm run dev
    # or
    # yarn dev
    ```

3.  Open [http://localhost:3000](http://localhost:3000) with your browser to see the application.

## Connecting to Backend

This frontend application is designed to work with its corresponding Go backend. Please ensure the backend server is running and the `NEXT_PUBLIC_API_BASE_URL` and `NEXT_PUBLIC_AUTH_BASE_URL` in your `.env.local` file are correctly pointing to it.

## Mock Data

A `mock_data.json` file is included in the root of this frontend project. It contains sample data for users, boards, lists, cards, collaborators, and comments. This data can be used for:
*   Understanding the expected data structures.
*   Manually seeding your backend database for development and testing purposes.
*   Potentially for frontend development if the backend is temporarily unavailable (though the application is primarily designed to work with a live backend).

To use this mock data, you would typically need a script or utility to parse this JSON and make appropriate API calls to your backend to create the entities, or directly insert them into your database if you have direct access.
