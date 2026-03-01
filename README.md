# Real-time Chat Application

A modern real-time chat application built with Go backend (WebSocket), Angular frontend, and PostgreSQL database. Features include user authentication, private messaging, group chats, and real-time notifications.

## Features

- **Real-time Messaging**: WebSocket-based instant messaging with no delays
- **User Authentication**: Secure JWT-based authentication
- **Private Chats**: One-on-one private messaging between users
- **Group Chats**: Create and manage group conversations
- **Online Status**: See which users are currently online
- **Message History**: Persistent storage of all messages
- **User Profiles**: Customizable user profiles with avatars
- **Notifications**: Real-time notifications for new messages
- **Message Search**: Search through chat history

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Chi Router
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT
- **WebSocket**: Gorilla WebSocket
- **Logging**: Zap
- **Testing**: Testify

### Frontend
- **Framework**: Angular 17+
- **Language**: TypeScript
- **UI Library**: Angular Material
- **HTTP Client**: HttpClient
- **WebSocket**: Angular WebSocket
- **State Management**: RxJS

## Project Structure

```
chat-app/
├── backend/
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   ├── user.go
│   │   └── message.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── chat.go
│   │   └── websocket.go
│   ├── services/
│   │   ├── user_service.go
│   │   ├── message_service.go
│   │   └── chat_service.go
│   ├── middleware/
│   │   └── auth.go
│   ├── database/
│   │   └── db.go
│   └── tests/
│       ├── auth_test.go
│       ├── chat_test.go
│       └── websocket_test.go
├── frontend/
│   ├── angular.json
│   ├── tsconfig.json
│   ├── package.json
│   ├── src/
│   │   ├── app/
│   │   │   ├── app.component.ts
│   │   │   ├── components/
│   │   │   │   ├── auth/
│   │   │   │   ├── chat/
│   │   │   │   └── sidebar/
│   │   │   ├── services/
│   │   │   │   ├── auth.service.ts
│   │   │   │   ├── chat.service.ts
│   │   │   │   └── websocket.service.ts
│   │   │   └── models/
│   │   │       ├── user.ts
│   │   │       └── message.ts
│   │   ├── main.ts
│   │   └── styles.css
│   └── tests/
│       ├── chat.component.spec.ts
│       └── auth.service.spec.ts
├── .gitignore
├── README.md
└── .github/
    └── workflows/
        └── ci-cd.yml
```

## Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+
- Angular CLI 17+
- PostgreSQL 14+

### Backend Setup

```bash
cd backend
go mod download
go mod tidy

# Create PostgreSQL database
createdb chat_app

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=chat_app
export JWT_SECRET=your_jwt_secret

# Run migrations (if available)
go run main.go migrate

# Start the server
go run main.go
```

The backend server will be available at `http://localhost:8080`

### Frontend Setup

```bash
cd frontend
npm install

# Development server
ng serve

# Build for production
ng build --prod
```

The frontend will be available at `http://localhost:4200`

## API Documentation

### Authentication Endpoints

**Register**
```
POST /api/auth/register
Content-Type: application/json

{
  "username": "user",
  "email": "user@example.com",
  "password": "password123"
}
```

**Login**
```
POST /api/auth/login
Content-Type: application/json

{
  "username": "user",
  "password": "password123"
}
```

### Chat Endpoints

**Get Conversations**
```
GET /api/chats
Authorization: Bearer <token>
```

**Create Group Chat**
```
POST /api/chats/group
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Project Team",
  "members": [1, 2, 3]
}
```

**WebSocket Connection**
```
WS /ws/chat/{chatId}?token=<jwt_token>
```

## Running Tests

### Backend Tests
```bash
cd backend
go test ./...
go test -cover ./...
```

### Frontend Tests
```bash
cd frontend
ng test
ng test --code-coverage
```

## Deployment

### Docker

Build and run with Docker:
```bash
# Backend
docker build -t chat-backend ./backend
docker run -p 8080:8080 chat-backend

# Frontend
docker build -t chat-frontend ./frontend
docker run -p 4200:4200 chat-frontend
```

### Docker Compose

```bash
docker-compose up -d
```

## Contributing

1. Create a feature branch (`git checkout -b feature/amazing-feature`)
2. Commit your changes (`git commit -m 'Add amazing feature'`)
3. Push to the branch (`git push origin feature/amazing-feature`)
4. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, email support@chatapp.com or open an issue on GitHub.
