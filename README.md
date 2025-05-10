## Event Management API

A simple RESTful API for managing events, built in Go with MySQL as the database.

### Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Database Setup](#database-setup)
- [Configuration](#configuration)
- [Running the Server](#running-the-server)
- [API Endpoints](#api-endpoints)
- [Data Model](#data-model)
- [Examples](#examples)
- [Error Handling](#error-handling)
- [License](#license)

---
## Features

- Create, read, update, and delete events
- JSON-based request and response
- Time parsing with `parseTime=true`
- Structured logging and HTTP status codes

---
## Prerequisites

- Go 1.18+
- MySQL 5.7+ (or compatible)
- `github.com/go-sql-driver/mysql` driver

---
## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/event-management-api.git
   cd event-management-api
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```

---
## Database Setup

1. Start MySQL server and create a database:
   ```sql
   CREATE DATABASE event_management;
   USE event_management;
   ```
2. Create the `events` table:
   ```sql
   CREATE TABLE events (
     id INT AUTO_INCREMENT PRIMARY KEY,
     title VARCHAR(255) NOT NULL,
     description TEXT,
     location VARCHAR(255),
     start_time DATETIME,
     end_time DATETIME,
     created_by VARCHAR(100),
     created_at DATETIME,
     updated_at DATETIME
   );
   ```

---
## Configuration

Update the `dsn` (Data Source Name) in `main.go` to match your MySQL credentials and host:

```go
const dsn = "root:password@tcp(localhost:3306)/event_management?parseTime=true"
```

---
## Running the Server

```bash
go run main.go
```

The server will start on port `8080`:
```
Server running on port: 8080
```

---
## API Endpoints

All endpoints use JSON for requests and responses.

| Method | Endpoint        | Description                    |
|--------|-----------------|--------------------------------|
| GET    | `/events`       | List all events                |
| POST   | `/events`       | Create a new event             |
| GET    | `/events/{id}`  | Retrieve a single event by ID  |
| PUT    | `/events/{id}`  | Update an existing event by ID |
| DELETE | `/events/{id}`  | Delete an event by ID          |

---
## Data Model

```json
{
  "id": 1,
  "title": "Team Meeting",
  "description": "Discuss project milestones",
  "location": "Conference Room A",
  "start_time": "2025-05-15T09:00:00Z",
  "end_time": "2025-05-15T10:00:00Z",
  "created_by": "alice@example.com",
  "created_at": "2025-05-01T12:00:00Z",
  "updated_at": "2025-05-01T12:00:00Z"
}
```

---
## Examples

### Create Event

```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Project Kickoff",
    "description": "Initial project meeting",
    "location": "Zoom",
    "start_time": "2025-05-20T14:00:00Z",
    "end_time": "2025-05-20T15:00:00Z",
    "created_by": "bob@example.com"
 }'
```

**Response**:
```json
{
  "message": "Event created",
  "data": { /* event object */ }
}
```

### List Events

```bash
curl http://localhost:8080/events
```

### Get Event by ID

```bash
curl http://localhost:8080/events/1
```

### Update Event

```bash
curl -X PUT http://localhost:8080/events/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "description": "Updated description",
    "location": "New Location",
    "start_time": "2025-05-15T11:00:00Z",
    "end_time": "2025-05-15T12:00:00Z",
    "created_by": "alice@example.com"
 }'
```

### Delete Event

```bash
curl -X DELETE http://localhost:8080/events/1
```

---
## Error Handling

- Returns `400 Bad Request` for invalid input or malformed IDs
- Returns `404 Not Found` if no event exists with the given ID
- Returns `405 Method Not Allowed` for unsupported HTTP methods
- Returns `500 Internal Server Error` for unexpected errors

---
## License

This project is released under the MIT License. Feel free to use and modify.
