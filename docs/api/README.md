# Godo API Documentation

## Overview

The Godo API provides a RESTful interface for managing tasks. All endpoints return JSON responses and follow standard HTTP status codes.

## Base URL

```
http://localhost:8080
```

## Authentication

Currently, the API is open and does not require authentication. Authentication will be added in a future update.

## Endpoints

### Health Check

Check if the API is running.

```
GET /health
```

Response:
```json
{
    "status": "ok"
}
```

### Tasks

#### List Tasks

Retrieve all tasks.

```
GET /api/v1/tasks
```

Response:
```json
[
    {
        "id": "uuid-string",
        "title": "Task title",
        "description": "Task description",
        "created_at": "2024-12-28T09:26:00Z",
        "updated_at": "2024-12-28T09:26:00Z",
        "completed_at": "2024-12-28T10:00:00Z"
    }
]
```

#### Create Task

Create a new task.

```
POST /api/v1/tasks
```

Request Body:
```json
{
    "title": "Buy groceries",
    "description": "Milk, bread, eggs",
    "completed_at": "2024-12-28T10:00:00Z"
}
```

Response:
```json
{
    "id": "uuid-string",
    "title": "Buy groceries",
    "description": "Milk, bread, eggs",
    "created_at": "2024-12-28T09:26:00Z",
    "updated_at": "2024-12-28T09:26:00Z",
    "completed_at": "2024-12-28T10:00:00Z"
}
```

#### Update Task (Full Update)

Update an existing task. Note: This endpoint requires all task fields to be provided.
A PATCH endpoint for partial updates is coming soon.

```
PUT /api/v1/tasks/{id}
```

Request Body:
```json
{
    "title": "Updated title",
    "description": "Updated description",
    "completed_at": "2024-12-28T10:00:00Z"
}
```

Response:
```json
{
    "id": "uuid-string",
    "title": "Updated title",
    "description": "Updated description",
    "created_at": "2024-12-28T09:26:00Z",
    "updated_at": "2024-12-28T09:26:15Z",
    "completed_at": "2024-12-28T10:00:00Z"
}
```

#### Delete Task

Delete a task.

```
DELETE /api/v1/tasks/{id}
```

Response: `204 No Content`

## Status Codes

- `200 OK` - Request succeeded
- `201 Created` - Resource created
- `204 No Content` - Request succeeded, no content to return
- `400 Bad Request` - Invalid request body
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Testing with HTTPie

[HTTPie](https://httpie.io/) is a user-friendly command-line HTTP client. Here are examples of using it with the Godo API:

```bash
# Health check
http :8080/health

# List all tasks
http :8080/api/v1/tasks

# Create a task
http POST :8080/api/v1/tasks title="Buy groceries" description="Milk, bread, eggs"

# Update a task
http PUT :8080/api/v1/tasks/{id} title="Updated title" description="New description"

# Delete a task
http DELETE :8080/api/v1/tasks/{id}
```

## Future Enhancements

- Authentication
- Request validation
- Rate limiting
- Pagination
- Sorting and filtering
- Search endpoint
- WebSocket support for real-time updates 