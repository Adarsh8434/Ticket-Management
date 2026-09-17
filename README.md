# Ticket Management System

A simple REST API built with Go and Gin for managing user tickets.

## Features

- User registration
- User login
- JWT-based authentication
- Password hashing using bcrypt
- Create tickets
- View own tickets
- View a specific own ticket
- Update ticket status
- Ticket ownership authorization
- SQLite database
- Docker support
- Deployed on Render

## Tech Stack

- Go
- Gin
- SQLite
- JWT
- bcrypt
- Docker

## Architecture

The application follows a layered architecture:

```text
Client
  |
  v
Handler
  |
  v
Service
  |
  v
Repository
  |
  v
SQLite Database
```

### Layers

- **Handler** - Handles HTTP requests and responses.
- **Service** - Contains business logic such as ticket status transitions.
- **Repository** - Handles database operations.
- **Model** - Defines application data structures.
- **Middleware** - Handles JWT authentication.

## Project Structure

```text
ticket-system/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── handler/
│   │   ├── auth_handler.go
│   │   └── ticket_handler.go
│   │
│   ├── middleware/
│   │   └── auth_middleware.go
│   │
│   ├── model/
│   │   ├── user.go
│   │   └── ticket.go
│   │
│   ├── repository/
│   │   ├── database.go
│   │   ├── user_repository.go
│   │   └── ticket_repository.go
│   │
│   └── service/
│       ├── auth_service.go
│       └── ticket_service.go
│
├── Dockerfile
├── .dockerignore
├── go.mod
├── go.sum
├── ticket.db
└── README.md
```

## API Endpoints

| Method | Endpoint             | Authentication |
|--------|-----------------------|-----------------|
| GET    | /health                | No              |
| POST   | /auth/register          | No              |
| POST   | /auth/login             | No              |
| POST   | /tickets                | Yes             |
| GET    | /tickets                | Yes             |
| GET    | /tickets/:id            | Yes             |
| PATCH  | /tickets/:id/status     | Yes             |

## Authentication

Protected endpoints require a JWT Bearer token.

```text
Authorization: Bearer <JWT_TOKEN>
```

The JWT contains the authenticated user's ID, which is used to enforce ticket ownership.

### Register User

**Request**

```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**

```json
{
  "id": 1,
  "email": "user@example.com"
}
```

Passwords are stored as bcrypt hashes and are never stored as plain text.

### Login

**Request**

```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response**

```json
{
  "token": "<JWT_TOKEN>"
}
```

Use the returned token for protected endpoints.

### Create Ticket

**Request**

```http
POST /tickets
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "title": "Unable to login",
  "description": "Login is not working properly"
}
```

**Response**

```json
{
  "id": 1,
  "user_id": 1,
  "title": "Unable to login",
  "description": "Login is not working properly",
  "status": "open"
}
```

New tickets are created with `open` status.

### Get Own Tickets

**Request**

```http
GET /tickets
Authorization: Bearer <JWT_TOKEN>
```

**Response**

```json
[
  {
    "id": 1,
    "user_id": 1,
    "title": "Unable to login",
    "description": "Login is not working properly",
    "status": "open"
  }
]
```

Only tickets belonging to the authenticated user are returned.

### Get Ticket by ID

**Request**

```http
GET /tickets/:id
Authorization: Bearer <JWT_TOKEN>
```

Example:

```http
GET /tickets/1
```

Users can only access tickets that they created.

If the ticket does not belong to the authenticated user:

```json
{
  "error": "ticket not found"
}
```

### Update Ticket Status

**Request**

```http
PATCH /tickets/:id/status
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "status": "in_progress"
}
```

**Response**

```json
{
  "message": "ticket status updated successfully"
}
```

## Ticket Status Flow

Tickets follow the following status flow:

```text
open
  |
  v
in_progress
  |
  v
closed
```

### Valid transitions

- `open` -> `in_progress`
- `in_progress` -> `closed`

### Invalid transitions

- `open` -> `closed`
- `in_progress` -> `open`
- `closed` -> `open`
- `closed` -> `in_progress`

A closed ticket cannot be reopened or changed.

## Authorization and Ownership

Every ticket belongs to the user who created it.

For ticket retrieval and status updates, the authenticated user's ID is checked against the ticket's owner ID.

This prevents one user from accessing or modifying another user's tickets.

Example:

```text
User 1
  |
  +-- Ticket 1

User 2
  |
  +-- Ticket 2
```

User 1 can access Ticket 1 but cannot access or update Ticket 2.

## Error Handling

The API returns meaningful HTTP status codes and JSON error messages.

Examples include:

**Invalid request**

```json
{
  "error": "invalid request"
}
```

**Unauthorized request**

```json
{
  "error": "invalid or expired token"
}
```

**Duplicate email**

```json
{
  "error": "email already registered"
}
```

**Ticket not found**

```json
{
  "error": "ticket not found"
}
```

**Invalid status**

```json
{
  "error": "invalid status"
}
```

## Run Locally

Make sure Go is installed.

Run the application:

```bash
go run ./cmd/server
```

The server runs on:

```text
http://localhost:8080
```

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

### Build the Application

```bash
go build ./...
```

## Run with Docker

Build the Docker image:

```bash
docker build -t ticket-system .
```

Run the container:

```bash
docker run -p 8080:8080 ticket-system
```

The application will be available at:

```text
http://localhost:8080
```

Test the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

### Docker Architecture

```text
Docker Container
      |
      v
Go Application
      |
      v
Gin REST API
      |
      v
SQLite Database
```

## Deployment

The application is deployed on Render.

**Deployed Application**

https://ticket-management-isk6.onrender.com/

**Health Check**

https://ticket-management-isk6.onrender.com/health

The health endpoint returns:

```json
{
  "status": "ok"
}
```

## Assumptions

- Users can only view tickets they created.
- Users can only update the status of their own tickets.
- Every ticket starts with `open` status.
- Valid statuses are `open`, `in_progress`, and `closed`.
- Status transitions follow `open -> in_progress -> closed`.
- A closed ticket cannot be reopened or changed.
- There is no admin role.
- There is no ticket assignment functionality.
- There are no comments on tickets.
- SQLite is used as the database.
- JWT is used for authentication.

## Testing Performed

The following functionality was tested:

- User registration
- Duplicate email validation
- User login
- JWT authentication
- Ticket creation
- Listing user's tickets
- Getting a ticket by ID
- Ticket ownership authorization
- Updating ticket status
  - `open -> in_progress`
  - `in_progress -> closed`
- Preventing closed ticket reopening
- Invalid status validation
- Invalid authentication
- Docker build
- Docker container execution
- Docker health check
- Deployed health check

## API Summary

```text
GET    /health

POST   /auth/register
POST   /auth/login

POST   /tickets
GET    /tickets
GET    /tickets/:id
PATCH  /tickets/:id/status
```
