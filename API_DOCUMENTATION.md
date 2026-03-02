# Chat App API Documentation

## Overview

Comprehensive REST API for a real-time chat application with advanced features including user presence tracking, analytics, admin controls, message management, and user preferences.

## Base URL

```
http://localhost:8080
```

## Authentication

All endpoints (except `/health` and public routes) require an `Authorization` header:

```
Authorization: Bearer {jwt_token}
```

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user

### Chat Management

- `GET /api/chats` - List all chats for authenticated user
- `POST /api/chats` - Create a private chat
- `POST /api/chats/group` - Create a group chat
- `GET /api/chats/{id}/messages` - Get messages in a chat
- `GET /api/chats/{id}/settings` - Get chat settings
- `PUT /api/chats/{id}/settings` - Update chat settings
- `GET /api/chats/{id}/pinned` - Get pinned messages
- `GET /api/chats/{id}/forwarded` - Get forwarded messages

### Messages

- `POST /api/messages/{id}/pin` - Pin a message
- `DELETE /api/messages/{id}/pin` - Unpin a message
- `GET /api/messages/{id}/pin-status` - Check pin status
- `POST /api/messages/{id}/forward` - Forward to single chat
- `POST /api/messages/{id}/forward-to-multiple` - Forward to multiple chats
- `POST /api/messages/{id}/react` - Add reaction to message
- `GET /api/messages/{id}/reactions` - Get message reactions

### File Management

- `POST /api/files/upload` - Upload file
- `GET /api/files/download/{filename}` - Download file
- `DELETE /api/files/delete/{filename}` - Delete file
- `GET /api/files/list` - List uploaded files

### Message Search

- `POST /api/search/messages` - Search messages with pagination
- `GET /api/search/messages/chat/{id}` - Search in specific chat
- `POST /api/search/advanced` - Advanced search with filters
- `GET /api/search/users` - Search for users
- `GET /api/search/chats` - Search for chats
- `GET /api/search/trending` - Get trending messages
- `GET /api/search/history/{userID}` - Get user search history

### User Presence & Status

- `POST /api/presence/update` - Update user presence status
- `GET /api/presence/user/{id}` - Get user presence
- `GET /api/presence/chat/{id}/members` - Get online members in chat
- `GET /api/presence/online` - Get count of online users
- `POST /api/presence/set-away` - Set user as away
- `POST /api/presence/history` - Log presence history

### User Preferences

- `GET /api/users/{id}/preferences/{chatID}` - Get chat preferences
- `PUT /api/users/{id}/preferences/{chatID}` - Update chat preferences
- `GET /api/users/{id}/notifications` - Get notification preferences
- `PUT /api/users/{id}/notifications` - Update notification preferences
- `POST /api/chats/{id}/mute/{userID}` - Mute chat for user

### User Blocking & Muting

- `POST /api/users/{id}/block/{targetID}` - Block a user
- `DELETE /api/users/{id}/block/{targetID}` - Unblock a user
- `GET /api/users/{id}/blocked` - List blocked users
- `POST /api/users/{id}/mute/{targetID}` - Mute a user
- `DELETE /api/users/{id}/mute/{targetID}` - Unmute a user
- `GET /api/users/{id}/muted` - List muted users
- `GET /api/users/{id}/is-blocked-by/{targetID}` - Check block status

### Admin - User Management

- `GET /api/admin/users` - List all users
- `GET /api/admin/users/{id}` - Get user details
- `POST /api/admin/users/{id}/suspend` - Suspend user
- `POST /api/admin/users/{id}/unsuspend` - Unsuspend user
- `DELETE /api/admin/users/{id}` - Delete user

### Admin - Chat Management

- `GET /api/admin/chats` - List all chats
- `GET /api/admin/chats/{id}` - Get chat details
- `DELETE /api/admin/chats/{id}` - Delete chat
- `POST /api/admin/chats/{id}/remove-member/{memberID}` - Remove member
- `POST /api/admin/chats/{id}/mute` - Mute chat

### Admin - Moderation & Reporting

- `POST /api/reports` - Create user report (public)
- `GET /api/admin/reports` - List all reports
- `GET /api/admin/reports/{id}` - Get report details
- `POST /api/admin/reports/{id}/resolve` - Resolve report
- `POST /api/admin/actions` - Log admin action
- `GET /api/admin/actions` - Get admin action history

### Analytics

- `GET /api/analytics/system` - Get system-wide analytics
- `GET /api/analytics/chat/{id}` - Get chat analytics
- `GET /api/analytics/user/{id}` - Get user analytics
- `POST /api/analytics/compute` - Compute/update analytics
- `GET /api/analytics/dashboard` - Get analytics dashboard summary

### WebSocket

- `GET /ws/chat/{id}` - WebSocket connection for real-time chat

## Request/Response Examples

### Create a Chat

```bash
curl -X POST http://localhost:8080/api/chats \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name": "General Chat", "type": "private"}'
```

### Update Presence

```bash
curl -X POST http://localhost:8080/api/presence/update \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "status": "online",
    "device_type": "web",
    "ip_address": "192.168.1.1"
  }'
```

### Search Messages

```bash
curl -X POST http://localhost:8080/api/search/messages \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "hello",
    "chat_id": 1,
    "limit": 20,
    "offset": 0
  }'
```

### Advanced Search with Filters

```bash
curl -X POST http://localhost:8080/api/search/advanced \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "important",
    "chat_id": 1,
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-12-31T23:59:59Z",
    "message_type": "text",
    "is_edited": false,
    "limit": 20,
    "sort_by": "relevance"
  }'
```

### Pin a Message

```bash
curl -X POST http://localhost:8080/api/messages/1/pin \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": 1,
    "pinned_by": 2,
    "reason": "Important announcement"
  }'
```

### Block a User

```bash
curl -X POST http://localhost:8080/api/users/1/block/2 \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"reason": "Spam"}'
```

### Get Chat Analytics

```bash
curl -X GET http://localhost:8080/api/analytics/chat/1 \
  -H "Authorization: Bearer {token}"
```

## Response Format

All responses are in JSON format:

```json
{
  "id": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "data": {}
}
```

## Error Handling

Errors are returned with appropriate HTTP status codes:

- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Status Codes

- `200` - OK
- `201` - Created
- `204` - No Content
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

## WebSocket Events

Real-time updates are sent via WebSocket connection:

- `message:new` - New message
- `message:edit` - Message edited
- `message:delete` - Message deleted
- `user:typing` - User typing indicator
- `user:online` - User came online
- `user:offline` - User went offline
- `reaction:add` - Reaction added
- `reaction:remove` - Reaction removed
- `presence:update` - User presence updated

## Database Models

The API uses the following models:

- **User** - User account information
- **Chat** - Chat room/conversation
- **Message** - Individual message
- **AdminUser** - Admin with special privileges
- **AdminAction** - Admin action audit log
- **UserReport** - User report for moderation
- **ChatAnalytics** - Chat metrics
- **UserAnalytics** - User metrics
- **SystemAnalytics** - System-wide metrics
- **UserPresence** - Real-time user status
- **PresenceHistory** - Historical presence data
- **ChatSettings** - Chat configuration
- **UserChatPreference** - User chat preferences
- **NotificationPreference** - User notification settings
- **PinnedMessage** - Pinned messages
- **BlockedUser** - User blocking relationships
- **MutedUser** - User muting relationships
- **ForwardedMessage** - Message forwarding

## Rate Limiting

API endpoints are rate-limited. Check response headers for limits:

- `X-RateLimit-Limit` - Request limit
- `X-RateLimit-Remaining` - Remaining requests
- `X-RateLimit-Reset` - Reset timestamp
