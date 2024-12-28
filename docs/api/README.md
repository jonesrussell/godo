# Godo API Documentation

This document describes the HTTP API for the Godo task management application.

## Base URL

The API is served at `http://localhost:8080` by default.

## API Versioning

The API is versioned using URL prefixes. The current version is `v1` and is accessed at `/api/v1/`.

## Authentication

Authentication is not currently implemented. All endpoints are publicly accessible.

## Endpoints

### Health Check

```
GET /health
```

Returns the health status of the API server.

**Response**
```json
{
  "status": "ok"
}
```

### List Tasks

```
GET /api/v1/tasks
```

Returns a list of all tasks.

**Response**
```json
[
  {
    "id": "string",
    "title": "string",
    "completed": boolean
  }
]
```

### Create Task

```
POST /api/v1/tasks
```

Creates a new task.

**Request Body**
```json
{
  "title": "string",
  "completed": boolean
}
```

**Response**
```json
{
  "id": "string",
  "title": "string",
  "completed": boolean
}
```

Status: 201 Created

### Update Task

```
PUT /api/v1/tasks/{id}
```

Updates an existing task.

**Parameters**
- `id`: Task ID (string, required)

**Request Body**
```json
{
  "title": "string",
  "completed": boolean
}
```

**Response**
```json
{
  "id": "string",
  "title": "string",
  "completed": boolean
}
```

### Delete Task

```
DELETE /api/v1/tasks/{id}
```

Deletes a task.

**Parameters**
- `id`: Task ID (string, required)

**Response**
Status: 204 No Content

## Error Responses

The API uses standard HTTP status codes to indicate the success or failure of requests:

- `200 OK`: Request succeeded
- `201 Created`: Resource was successfully created
- `204 No Content`: Request succeeded with no response body
- `400 Bad Request`: Invalid request body or parameters
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

Error responses include a message in the response body:

```json
{
  "error": "string"
}
```

## Future Enhancements

1. Authentication and authorization
2. Rate limiting
3. Request validation
4. Pagination for list endpoints
5. Search and filtering
6. WebSocket support for real-time updates 