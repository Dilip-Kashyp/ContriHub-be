# ContriHub Backend

The backend engine for ContriHub, built with Go and designed for high performance, scalability, and deep AI integration.

## Tech Stack

- **Language**: [Go (Golang)](https://golang.org/)
- **Framework**: [Gin Gonic](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Cache/Rate Limiting**: [Redis](https://redis.io/)
- **AI Integration**: Groq LLM / Custom AI Handlers

## Core Features

- **GitHub Proxy**: Securely handles GitHub API requests with efficient caching.
- **AI Orchestration**: Endpoints for codebase explanation, roadmap generation, and intelligent project matching.
- **Rate Limiting**: IP and Token-based rate limiting using Redis to ensure service stability.
- **Auto-Migrations**: Seamless database schema updates using GORM.
- **Context-Aware Chat**: Stores and retrieves AI chat history for a persistent developer assistant experience.

## API Architecture

The API is versioned under `/api/v1` and follows strict RESTful principles.

### Key Endpoints

#### Authentication
- `GET /auth/login`: Initiates GitHub OAuth flow.
- `GET /auth/github/callback`: Handles GitHub OAuth callback.

#### AI Services (Rate Limited)
- `POST /api/v1/ai/explain`: Generates a high-level explanation of a repository.
- `POST /api/v1/ai/roadmap`: Creates a personalized learning roadmap.
- `POST /api/v1/ai/find-projects`: Matches developers with suitable projects.
- `GET /api/v1/ai/chat-history`: Retrieves past conversation logs.
- `POST /api/v1/ai/chat`: Submits a new message to the Gibo AI.

## Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL
- Redis

### Installation

1. **Clone and Navigate**
   ```bash
   cd backend
   ```

2. **Environment Configuration**
   Copy `.env.example` to `.env` and fill in your credentials:
   - `DB_URL`: PostgreSQL connection string.
   - `REDIS_URL`: Redis connection string.
   - `GITHUB_CLIENT_ID` & `GITHUB_CLIENT_SECRET`: Your GitHub OAuth app credentials.
   - `GROQ_API_KEY`: API key for AI features.

3. **Run the Server**
   ```bash
   go run main.go
   ```
   The server will start on port `5050` by default.

## License

MIT License.
