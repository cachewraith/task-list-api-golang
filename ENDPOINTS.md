# API Endpoints Reference

Quick reference for testing all API endpoints in Postman.

## Base URL
```
http://localhost:8080
```

---

## Public Endpoints (No Auth Required)

### 1. Health Check
```
GET /health
```

**Headers:** None

**Body:** None

**Expected Response:**
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "up"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 2. Register User
```
POST /api/v1/auth/register
```

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Expected Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john@example.com",
    "name": "John Doe",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 3. Login
```
POST /api/v1/auth/login
```

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Expected Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john@example.com",
      "name": "John Doe",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 4. Logout
```
POST /api/v1/auth/logout
```

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:** None

**Note:** With JWT authentication, the actual logout happens client-side by deleting the token. This endpoint is provided for completeness and could be extended to implement token blacklisting.

**Expected Response:** `200 OK`
```json
{
  "success": true,
  "message": "Logout successful. Please delete your token on the client side.",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## Protected Endpoints (Auth Required)

**Required Header for all protected endpoints:**
```
Authorization: Bearer YOUR_JWT_TOKEN_HERE
Content-Type: application/json
```

---

### 4. Get All Todos
```
GET /api/v1/todos
```

**Headers:**
```
Authorization: Bearer {{token}}
```

**Body:** None

**Expected Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Buy groceries",
      "description": "Milk, eggs, bread",
      "completed": false,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "meta": {
    "total": 1
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 5. Get Todo by ID
```
GET /api/v1/todos/{{todo_id}}
```

**Headers:**
```
Authorization: Bearer {{token}}
```

**Body:** None

**Expected Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "completed": false,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 6. Create Todo
```
POST /api/v1/todos
```

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread"
}
```

**Expected Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "completed": false,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 7. Create Multiple Todos (Bulk)
```
POST /api/v1/todos/bulk
```

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
[
  {
    "title": "Complete Monthly Report",
    "description": "Finish the financial audit report"
  },
  {
    "title": "Leg Day Workout",
    "description": "Focus on squats and deadlifts"
  },
  {
    "title": "Study Go Interfaces",
    "description": "Read official Go documentation"
  }
]
```

**Expected Response:** `201 Created`
```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Complete Monthly Report",
      "description": "Finish the financial audit report",
      "completed": false,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Leg Day Workout",
      "description": "Focus on squats and deadlifts",
      "completed": false,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "meta": {
    "total": 2
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 8. Update Todo
```
PUT /api/v1/todos/{{todo_id}}
```

**Headers:**
```
Authorization: Bearer {{token}}
Content-Type: application/json
```

**Body:**
```json
{
  "title": "Buy groceries and more",
  "completed": true
}
```

**Expected Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Buy groceries and more",
    "description": "Milk, eggs, bread",
    "completed": true,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:05:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 9. Delete Todo
```
DELETE /api/v1/todos/{{todo_id}}
```

**Headers:**
```
Authorization: Bearer {{token}}
```

**Body:** None

**Expected Response:** `200 OK`
```json
{
  "success": true,
  "message": "Todo deleted successfully",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 10. View Application Logs
```
GET /api/v1/logs?limit=100&level=INFO
```

**Headers:**
```
Authorization: Bearer {{token}}
```

**Query Parameters:**
- `limit` (optional): Number of logs to return (default: 100, max: 1000)
- `level` (optional): Filter by level (DEBUG, INFO, WARN, ERROR)

**Body:** None

**Expected Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "entries": [
      {
        "timestamp": "2024-01-15T10:30:00Z",
        "level": "INFO",
        "message": "Successfully created todo: 550e8400...",
        "context": "TODO_SERVICE",
        "function": "Create"
      },
      {
        "timestamp": "2024-01-15T10:29:00Z",
        "level": "INFO",
        "message": "Login attempt for email: john@example.com",
        "context": "AUTH_SERVICE",
        "function": "Login"
      }
    ],
    "total": 2,
    "limit": 100
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 11. Clear Application Logs
```
DELETE /api/v1/logs
```

**Headers:**
```
Authorization: Bearer {{token}}
```

**Body:** None

**Expected Response:** `200 OK`
```json
{
  "success": true,
  "message": "All logs cleared",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## Postman Setup Tips

### 1. Create a Collection
Name it: `Todo API`

### 2. Create Environment Variables
Create an environment (e.g., `Local`) with these variables:
| Variable | Initial Value |
|----------|---------------|
| `base_url` | `http://localhost:8080` |
| `token` | (empty - will be set by login request) |

### 3. Set Token Automatically (Login Request)
In the **Tests** tab of your Login request, add:
```javascript
var jsonData = pm.response.json();
pm.environment.set("token", jsonData.data.token);
```

This automatically saves the token to your environment after login.

### 4. Use Variables in Requests
- URL: `{{base_url}}/api/v1/todos`
- Header: `Authorization: Bearer {{token}}`

### 5. Test Flow
1. Register a user (`POST /api/v1/auth/register`)
2. Login (`POST /api/v1/auth/login`) - token auto-saved
3. Create todos (`POST /api/v1/todos`)
4. Get todos (`GET /api/v1/todos`)
5. Update todo (`PUT /api/v1/todos/{{todo_id}}`)
6. Delete todo (`DELETE /api/v1/todos/{{todo_id}}`)
7. View logs (`GET /api/v1/logs?limit=100`)
8. Clear logs (`DELETE /api/v1/logs`)
9. Logout (`POST /api/v1/auth/logout`) - client deletes token

---

## Summary Table
| # | Method | Endpoint | Auth | Description |
|---|--------|----------|------|-------------|
| 1 | GET | `/health` | ❌ | Health check |
| 2 | POST | `/api/v1/auth/register` | ❌ | Register user |
| 3 | POST | `/api/v1/auth/login` | ❌ | Login (get token) |
| 4 | POST | `/api/v1/auth/logout` | ✅ | Logout (client deletes token) |
| 5 | GET | `/api/v1/todos` | ✅ | List all todos |
| 6 | GET | `/api/v1/todos/{{id}}` | ✅ | Get single todo |
| 7 | POST | `/api/v1/todos` | ✅ | Create todo |
| 8 | POST | `/api/v1/todos/bulk` | ✅ | Create multiple |
| 9 | PUT | `/api/v1/todos/{{id}}` | ✅ | Update todo |
| 10 | DELETE | `/api/v1/todos/{{id}}` | ✅ | Delete todo |
| 11 | GET | `/api/v1/logs` | ✅ | View logs |
| 12 | DELETE | `/api/v1/logs` | ✅ | Clear logs |

---

## Error Responses

### 400 Bad Request (Validation Error)
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "title is required"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 401 Unauthorized (Missing/Invalid Token)
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authorization header required"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 404 Not Found
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "todo not found"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## Telegram Alerts (Optional)

The API can send real-time alerts to Telegram when ERROR or WARN level logs occur.

### Setup

1. **Create Telegram Bot:**
   - Message @BotFather on Telegram
   - Create new bot → get `TELEGRAM_BOT_TOKEN`

2. **Get Group ID:**
   - Add bot to your group
   - Send any message in group
   - Visit: `https://api.telegram.org/bot<token>/getUpdates`
   - Look for `"chat":{"id":-100xxxxxxxxx` → that's your `TELEGRAM_GROUP_ID`

3. **Configure Environment:**
```bash
# Add to .env file
TELEGRAM_BOT_TOKEN=7862482571:AAHxxxxx
TELEGRAM_GROUP_ID=-1003388054167
TELEGRAM_THREAD_ID=305  # Optional: for topic groups
```

### Alert Format

🚨 **ERROR**

**Time:** 2024-01-15 14:30:25  
**Context:** AUTH_SERVICE  
**Function:** Login  

**Message:**
```
Login failed - invalid password for user: john@example.com
```

### Features

- ✅ Non-blocking (async sending)
- ✅ Configurable minimum level (ERROR, WARN, INFO, DEBUG)
- ✅ Only sends ERROR and WARN by default
- ✅ HTML formatted messages
- ✅ Supports topic threads (thread_id)
- ✅ 5 second timeout with automatic fail-safe

---

## Complete API Summary

**Total Endpoints:** 13
- Public: 3 (health, register, login)
- Protected: 10 (logout, todos CRUD, logs)

**Features:**
- JWT Authentication
- User ownership of todos
- Bulk operations
- Application logging with viewing
- Telegram alerts for errors
