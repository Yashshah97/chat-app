# Real-time Chat Application

A comprehensive real-time chat application built with Go backend (WebSocket), Angular frontend, and PostgreSQL database. Features 50+ advanced capabilities including real-time messaging, enterprise security, analytics, and data management.

## 🌟 Core Features

- **Real-time Messaging**: WebSocket-based instant messaging with no delays
- **User Authentication**: Secure JWT-based authentication with 2FA support
- **Private & Group Chats**: One-on-one private messaging and group conversations
- **Online Status**: See which users are currently online with presence tracking
- **Message History**: Persistent storage with search and filtering
- **User Profiles**: Customizable profiles with status management
- **Real-time Notifications**: Multi-channel notifications with scheduling

## 🔐 Security & Enterprise Features (Commits #42-50)

### Authentication & Authorization
- **Two-Factor Authentication (2FA)**: Email, SMS, TOTP, backup codes, security keys
- **Role-Based Access Control (RBAC)**: Granular permission management
- **Session Management**: Active session tracking and revocation
- **Login Attempt Monitoring**: Track suspicious login activities
- **Trusted Device Management**: Mark devices as trusted for streamlined access

### Data Security
- **End-to-End Encryption**: AES-256 encryption for messages and sensitive data
- **Password Policies**: Configurable complexity requirements and expiration
- **Password History**: Prevent password reuse with history tracking
- **Secure Deletion**: Cryptographically secure data deletion
- **Encryption Key Management**: User-managed encryption keys with key rotation
- **OAuth Integration**: Support for third-party authentication

### Compliance & Audit
- **Comprehensive Audit Logging**: Track all user actions and system changes
- **Security Incident Tracking**: Monitor and respond to security incidents
- **GDPR Compliance**: Data export, deletion, and privacy controls
- **Data Backup & Recovery**: Scheduled backups with recovery capabilities
- **Audit Trail Queries**: Advanced filtering for compliance reporting

## 📊 Analytics & Monitoring (Commits #43-47)

### Real-time Analytics
- **Performance Metrics**: CPU, memory, database, and network monitoring
- **Application Metrics**: Request/response times, error rates, throughput
- **Usage Statistics**: Per-user activity tracking and trends
- **Error Tracking**: Automatic error detection and aggregation
- **System Health Dashboard**: Real-time system status monitoring
- **Performance Alerts**: Configurable alerts for metric thresholds

### Advanced Analytics
- **User Engagement Metrics**: Active users, session duration, feature usage
- **Chat Analytics**: Message volume, response times, user interactions
- **Admin Dashboard**: System-wide statistics and insights
- **Custom Reports**: Generate reports by date range, user, or chat

## 🔔 Advanced Notifications (Commit #46)

- **Multi-Channel Delivery**: Email, SMS, push notifications, webhooks
- **Notification Scheduling**: Send notifications at optimal times
- **Batch Notifications**: Bulk send campaigns to multiple users
- **Notification Templates**: Pre-built and custom message templates
- **Delivery Tracking**: Monitor notification delivery status
- **Preference Management**: User-controlled notification settings

## 🎯 Message Management & Filtering (Commit #48)

- **Advanced Message Search**: Search by keyword, date, sender, attachments
- **Filter Rules**: Create reusable filters for automatic message organization
- **Message Templates**: Quick reply templates and message prefabs
- **Message Editing**: Edit messages with full history tracking
- **Message Reactions**: Emoji reactions with custom emoji packs
- **Message Pinning**: Pin important messages to chat
- **Message Forwarding**: Forward messages to other chats
- **Read Receipts**: See when messages are read
- **Typing Indicators**: Real-time typing notifications

## 📦 Data Management (Commits #49-50)

### Batch Operations
- **Bulk Message Operations**: Delete, archive, or process multiple messages
- **Batch Job Monitoring**: Track progress of bulk operations
- **Batch Scheduling**: Schedule batch jobs for off-peak hours
- **Error Recovery**: Automatic retry for failed batch items

### Data Migration & Import/Export
- **Chat Export**: Export conversations in multiple formats (JSON, CSV, PDF)
- **Data Import**: Import messages from other chat platforms
- **Migration Tools**: Data migration with field mapping and transformation
- **Scheduled Backups**: Automatic backup schedules with retention policies
- **Archive Management**: Archive old chats and conversations
- **GDPR Data Exports**: User data export for compliance

## 🔗 Integration & Extensions (Previously Implemented)

- **Third-Party Integrations**: Slack, Discord, GitHub, Jira connectors
- **Webhook Support**: Incoming and outgoing webhooks for events
- **API Integration**: RESTful API with versioning and deprecation
- **Rate Limiting**: Per-user and per-endpoint rate limiting
- **Emoji & Sticker Packs**: Curated and user-created emoji packs with reviews
- **Chatbots**: NLP-powered chatbots with intent recognition

## 🎮 Gamification & Social Features

- **Badge System**: Achievement badges for user milestones
- **Chat Invitations**: Invite users to chats with custom tokens
- **Trending Topics**: Trending conversation topics and hashtags
- **Summary Statistics**: Daily/weekly chat summaries
- **User Profiles**: Rich profiles with status, bio, and achievements

## 📊 Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Chi Router
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT + 2FA
- **WebSocket**: Gorilla WebSocket
- **Encryption**: AES-256, bcrypt
- **Logging**: Structured logging
- **Testing**: Comprehensive test coverage

### Frontend
- **Framework**: Angular 17+
- **Language**: TypeScript
- **UI Library**: Angular Material
- **HTTP Client**: HttpClient with interceptors
- **WebSocket**: Real-time communication
- **State Management**: RxJS, NgRx

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
│   ├── main.go                          # Server entry point with 50+ features
│   ├── go.mod / go.sum
│   │
│   ├── Core Models & Handlers
│   ├── models.go                        # Base User, Chat, Message models
│   ├── auth.go                          # Authentication logic
│   ├── websocket.go                     # WebSocket implementation
│   ├── file_upload.go                   # File upload management
│   ├── reactions.go                     # Message reactions
│   ├── read_receipts.go                 # Read receipt tracking
│   │
│   ├── Security (Commits #43-44)
│   ├── security_models.go               # 2FA, session, audit models
│   ├── security_handlers.go             # 2FA and session endpoints
│   ├── security_advanced_models.go      # Encryption, password policies
│   ├── security_advanced_handlers.go    # Encryption and audit endpoints
│   │
│   ├── Analytics (Commit #45)
│   ├── analytics_models.go              # System analytics models
│   ├── analytics_handlers.go            # Analytics endpoints
│   ├── analytics_enhanced_models.go     # Performance metrics, errors
│   │
│   ├── Notifications (Commit #46)
│   ├── notification_models.go           # Core notification models
│   ├── notification_handlers.go         # Notification endpoints
│   ├── notification_improvements_models.go  # Scheduling, batching
│   ├── notification_delivery_handlers.go    # Delivery management
│   ├── notification_templates_handlers.go   # Template management
│   │
│   ├── Performance & Monitoring (Commit #47)
│   ├── performance_monitoring_models.go # CPU/Memory/Database metrics
│   │
│   ├── Filtering (Commit #48)
│   ├── advanced_filtering_models.go     # Message filters, rules
│   ├── advanced_search_handlers.go      # Search implementation
│   │
│   ├── Batch Operations (Commit #49)
│   ├── batch_operations_models.go       # Batch job management
│   │
│   ├── Data Migration (Commit #50)
│   ├── migration_tools_models.go        # Migration and import/export
│   │
│   ├── Previously Implemented Features
│   ├── activity_log_models.go / handlers.go      # Audit logging
│   ├── emoji_models.go / handlers.go             # Emoji & stickers
│   ├── backup_models.go / handlers.go            # Chat backups
│   ├── permission_models.go / handlers.go        # RBAC
│   ├── webhook_models.go / handlers.go           # Webhooks
│   ├── ratelimit_models.go / handlers.go         # Rate limiting
│   ├── call_models.go / handlers.go              # Voice/video calls
│   ├── integration_models.go / handlers.go       # 3rd party integrations
│   ├── template_models.go / handlers.go          # Message templates
│   ├── export_models.go / handlers.go            # Data export
│   ├── profile_models.go / handlers.go           # User profiles
│   ├── gamification_models.go / handlers.go      # Badges, achievements
│   ├── preference_models.go / handlers.go        # User preferences
│   │
│   └── Other Features
│       ├── admin_models.go / handlers.go         # Admin functions
│       ├── analytics_models.go / handlers.go     # Basic analytics
│       ├── blocking_models.go / handlers.go      # User blocking
│       ├── chat_management_handlers.go           # Chat management
│       ├── forward_handlers.go                   # Message forwarding
│       ├── message_edit_models.go / handlers.go  # Message editing
│       ├── moderation_handlers.go                # Moderation
│       ├── notification_handlers.go              # Notifications
│       ├── pinned_messages_handlers.go           # Pinned messages
│       ├── presence_models.go / handlers.go      # Presence tracking
│       ├── reaction_handlers.go                  # Emoji reactions
│       ├── read_receipt_handlers.go              # Read receipts
│       ├── search_handlers.go                    # Message search
│       ├── settings_models.go / handlers.go      # Settings
│
├── frontend/
│   ├── angular.json / tsconfig.json
│   ├── package.json
│   ├── src/
│   │   ├── app/
│   │   │   ├── components/
│   │   │   │   ├── auth/                 # Login, registration
│   │   │   │   ├── chat/                 # Chat interface
│   │   │   │   ├── sidebar/              # Chat list
│   │   │   │   ├── notifications/        # Notification UI
│   │   │   │   ├── security/             # 2FA, security settings
│   │   │   │   └── analytics/            # Analytics dashboard
│   │   │   ├── services/
│   │   │   │   ├── auth.service.ts
│   │   │   │   ├── chat.service.ts
│   │   │   │   ├── websocket.service.ts
│   │   │   │   ├── notification.service.ts
│   │   │   │   ├── security.service.ts
│   │   │   │   └── analytics.service.ts
│   │   │   └── models/
│   │   ├── main.ts / styles.css
│   │   └── assets/
│   └── tests/
│
├── API_DOCUMENTATION.md               # Comprehensive API docs
├── README.md                          # This file
├── .gitignore
├── .github/
│   └── workflows/
│       └── ci-cd.yml                  # CI/CD pipeline
└── LICENSE
```

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Node.js 18+
- Angular CLI 17+
- PostgreSQL 14+
- Docker & Docker Compose (optional)

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
export JWT_SECRET=your_jwt_secret_key_here

# Start the server (migrations run automatically)
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

## 📚 API Documentation

For comprehensive API documentation, see [API_DOCUMENTATION.md](API_DOCUMENTATION.md)

### Quick Examples

**Register User**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

**Login**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "password": "SecurePass123!"
  }'
```

**Enable 2FA**
```bash
curl -X POST http://localhost:8080/api/security/2fa/enable \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"method": "totp"}'
```

**Create Group Chat**
```bash
curl -X POST http://localhost:8080/api/chats/group \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Project Team",
    "members": [2, 3, 4]
  }'
```

**WebSocket Connection**
```javascript
const ws = new WebSocket(
  'ws://localhost:8080/ws/chat/1?token=<jwt_token>'
);

ws.onmessage = (event) => {
  console.log('Message:', event.data);
};

ws.send(JSON.stringify({
  type: 'message',
  body: 'Hello everyone!',
  chatId: 1
}));
```

## ✅ Running Tests

### Backend Tests
```bash
cd backend
go test ./...
go test -cover ./...
go test -v ./...
```

### Frontend Tests
```bash
cd frontend
ng test
ng test --code-coverage
ng test --browsers=Chrome
```

### Integration Tests
```bash
cd backend
go test -tags=integration ./...
```

## 🐳 Deployment

### Docker Compose (Recommended)

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database
- Go backend server (http://localhost:8080)
- Angular frontend (http://localhost:4200)

### Manual Docker Setup

**Backend**
```bash
docker build -t chat-backend ./backend
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=password \
  -e DB_NAME=chat_app \
  -e JWT_SECRET=your_secret \
  chat-backend
```

**Frontend**
```bash
docker build -t chat-frontend ./frontend
docker run -p 4200:4200 chat-frontend
```

### Kubernetes Deployment

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend.yaml
```

### Environment Variables

**Backend**
```
DB_HOST=localhost           # PostgreSQL host
DB_PORT=5432               # PostgreSQL port
DB_USER=postgres           # PostgreSQL user
DB_PASSWORD=password       # PostgreSQL password
DB_NAME=chat_app           # Database name
JWT_SECRET=your_secret     # JWT signing secret
PORT=8080                  # Backend port
ENV=production             # Environment (development/production)
LOG_LEVEL=info             # Logging level
```

**Frontend**
```
API_BASE_URL=http://localhost:8080    # Backend API URL
WS_URL=ws://localhost:8080            # WebSocket URL
ENVIRONMENT=production                # Environment
```

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Workflow

```bash
# Create feature branch
git checkout -b feature/your-feature

# Make changes and test
go test ./...
ng test

# Commit with descriptive message
git commit -m "feat: Add amazing feature"

# Push and create PR
git push origin feature/your-feature
```

### Code Standards

- Follow Go conventions for backend code
- Use TypeScript strict mode for frontend
- Write tests for new features
- Maintain 80%+ code coverage
- Update documentation with changes

## 📊 Commit History

This project includes 50+ commits with comprehensive features:

**Session 1-41**: Core chat functionality and initial features
**Session 42**: User preferences, surveys, recommendation engine
**Session 43**: Two-factor authentication and security management
**Session 44**: Advanced encryption, password policies, data security
**Session 45**: Enhanced analytics, performance metrics, error tracking
**Session 46**: Notification scheduling, batching, advanced channels
**Session 47**: Performance monitoring (CPU/Memory/Database metrics)
**Session 48**: Advanced message filtering and filter rules
**Session 49**: Batch operations system for bulk message processing
**Session 50**: Migration tools for data import/export/sync operations

## 🔐 Security

- All passwords are hashed using bcrypt
- JWT tokens for authentication
- Support for 2FA (Email, SMS, TOTP, Security Keys)
- End-to-end encryption for sensitive data
- SQL injection protection via ORM
- CORS enabled for specified origins
- Rate limiting per endpoint and user
- Input validation and sanitization
- Regular security audits recommended

## 📈 Performance

- Optimized database queries with indexes
- WebSocket connection pooling
- Message caching with Redis (optional)
- Pagination for large datasets
- Lazy loading of images and files
- Database connection pooling
- Load balancing ready
- Horizontal scaling support

## 🐛 Known Issues & Limitations

- Frontend requires modern browser (Chrome, Firefox, Safari, Edge)
- PostgreSQL 14+ required for full feature set
- Real-time features require WebSocket support
- File uploads limited to 100MB per file
- Chat history limited to last 10,000 messages per chat

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💬 Support & Contact

For support and questions:
- Email: support@chatapp.com
- GitHub Issues: [Create an issue](https://github.com/Yashshah97/chat-app/issues)
- Discussions: [GitHub Discussions](https://github.com/Yashshah97/chat-app/discussions)

## 🙏 Acknowledgments

- Go community and libraries
- Angular team
- PostgreSQL developers
- Contributors and users

## 📝 Changelog

See [CHANGELOG.md](CHANGELOG.md) for detailed version history.

## 🚀 Roadmap

- [ ] Mobile app (React Native)
- [ ] Video conferencing
- [ ] AI-powered chat suggestions
- [ ] Advanced search with Elasticsearch
- [ ] Message translation
- [ ] Voice transcription
- [ ] Custom themes and branding
- [ ] Enterprise SSO support
- [ ] GraphQL API
- [ ] Microservices architecture
