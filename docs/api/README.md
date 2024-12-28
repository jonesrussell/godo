# Godo API Documentation

## Overview

The Godo HTTP API provides RESTful endpoints for task management and real-time updates via WebSocket.

## API Versions

- v1 (Current): `/api/v1/`

## Authentication

Currently, the API is unauthenticated. Authentication will be added in future versions.

## Endpoints

### Tasks

#### GET /api/v1/tasks
List all tasks.

**Response**
```json
{
  "tasks": [
    {
      "id": "string",
      "title": "string",
      "completed": boolean
    }
  ]
}
```

#### GET /api/v1/tasks/{id}
Get a specific task.

**Response**
```json
{
  "id": "string",
  "title": "string",
  "completed": boolean
}
```

#### POST /api/v1/tasks
Create a new task.

**Request**
```json
{
  "title": "string",
  "completed": boolean
}
```

#### PUT /api/v1/tasks/{id}
Update an existing task.

**Request**
```json
{
  "title": "string",
  "completed": boolean
}
```

#### DELETE /api/v1/tasks/{id}
Delete a task.

### WebSocket

#### WS /api/v1/ws
WebSocket endpoint for real-time updates.

**Events**
- `task.created`
- `task.updated`
- `task.deleted`

## Error Responses

All errors follow this format:
```json
{
  "error": {
    "code": "string",
    "message": "string"
  }
}
```

## Implementation Status

See the [TODO.md](../../TODO.md) file for current implementation status. 