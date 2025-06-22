# SMS OTP Service

A professional SMS OTP (One-Time Password) service built with Go, Fiber, GORM, and PostgreSQL using Clean Architecture principles.

## Features

- 📱 SMS OTP generation and verification
- 🔒 Rate limiting and security controls
- 🏗️ Clean Architecture (Onion Architecture)
- 📊 Swagger/OpenAPI documentation
- 🐳 Docker and Docker Compose support
- 💾 PostgreSQL database with GORM
- 📝 Structured logging
- 🧪 Health checks and monitoring

## Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make (optional)

### Installation

```bash
# Clone repository
git clone https://github.com/your-username/sms-otp-service.git
cd sms-otp-service

# Setup project
make setup

# Start with Docker (recommended)
make docker-run

# Or start locally
make dev
```

### Verify Installation

```bash
curl http://localhost:8080/health
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/otp/send` | Send OTP to phone number |
| POST | `/api/v1/otp/verify` | Verify OTP code |
| POST | `/api/v1/otp/resend` | Resend OTP to phone number |
| GET | `/health` | Service health check |
| GET | `/ready` | Readiness probe |
| GET | `/docs/` | Swagger documentation |

## API Usage

### Send OTP

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/otp/send \
-H "Content-Type: application/json" \
-d '{
"phone_number": "+994501234567",
"purpose": "verification"
}'
```

**Response:**
```json
{
"success": true,
"message": "OTP sent successfully",
"expires_in": 300,
"id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Console Output (Development):**
```
=== MOCK SMS ===
To: +994501234567
Message: Your verification code is: 123456. Valid for 5 minutes.
================
```

### Verify OTP

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/otp/verify \
-H "Content-Type: application/json" \
-d '{
"phone_number": "+994501234567",
"code": "123456",
"purpose": "verification"
}'
```

**Success Response:**
```json
{
"success": true,
"message": "OTP verified successfully",
"verified_at": "2024-01-15T10:30:00Z"
}
```

**Error Response:**
```json
{
"success": false,
"message": "Invalid OTP code. Please try again."
}
```

### Resend OTP

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/otp/resend \
-H "Content-Type: application/json" \
-d '{
"phone_number": "+994501234567",
"purpose": "verification"
}'
```

## OTP Purposes

| Purpose | Description | Use Case |
|---------|-------------|----------|
| `verification` | Phone number verification | Account registration |
| `login` | Two-factor authentication | Secure login |
| `reset` | Password reset | Password recovery |

## Configuration

### Environment Variables

```bash
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=sms_otp_db

# SMS
SMS_PROVIDER=mock
SMS_API_KEY=your_api_key
SMS_SENDER_NAME=OTPService

# OTP
OTP_VALIDITY_MINUTES=5
OTP_RATE_LIMIT_MINUTES=10
OTP_MAX_PER_PERIOD=3
OTP_CODE_LENGTH=6
OTP_MAX_ATTEMPTS=3

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Environment File

Copy and modify the environment template:

```bash
cp .env.example .env
nano .env
```

## Security Features

### Rate Limiting
- 3 OTP requests per 10 minutes per phone number
- 1 minute cooldown between resend requests
- Maximum 3 verification attempts per OTP

### Security Controls
- Cryptographically secure OTP generation
- Automatic OTP expiration (5 minutes)
- Phone number format validation
- Automatic cleanup of expired OTPs

## Development

### Available Commands

```bash
make help                # Show all commands
make setup              # Initial project setup
make dev                # Start development server
make build              # Build application
make test               # Run tests
make test-coverage      # Run tests with coverage
make swagger            # Generate API documentation
make docker-run         # Start with Docker
make docker-logs        # View Docker logs
make api-test           # Test API endpoints
make clean              # Clean build artifacts
```

### Testing

**Automated Testing:**
```bash
make api-test
```

**Manual Testing:**
```bash
# 1. Send OTP
curl -X POST http://localhost:8080/api/v1/otp/send \
-H "Content-Type: application/json" \
-d '{"phone_number": "+994501234567"}'

# 2. Check console for OTP code
# 3. Verify OTP with the code from console
curl -X POST http://localhost:8080/api/v1/otp/verify \
-H "Content-Type: application/json" \
-d '{"phone_number": "+994501234567", "code": "123456"}'
```

## Deployment

### Docker

```bash
# Build and run
docker-compose up -d

# View logs
docker-compose logs -f sms-otp-service

# Stop services
docker-compose down
```

### Production

1. Set production environment variables
2. Change `SMS_PROVIDER` from `mock` to a real provider
3. Use secure database credentials
4. Enable SSL/TLS

## Architecture

The project follows Clean Architecture principles:

```
internal/
├── domain/                 # Business entities and rules
│   ├── entities/          # Domain entities (OTP)
│   ├── repositories/      # Repository interfaces
│   └── services/          # Domain services
├── application/           # Application logic
│   ├── usecases/         # Use cases
│   └── dto/              # Data transfer objects
├── infrastructure/       # External concerns
│   ├── database/         # Database connection
│   ├── repositories/     # Repository implementations
│   ├── sms/             # SMS service implementations
│   └── config/          # Configuration
└── interfaces/           # External interfaces
└── http/            # HTTP handlers and routes
```

## Tech Stack

- **Language:** Go 1.24+
- **Web Framework:** Fiber v2
- **Database:** PostgreSQL with GORM
- **Documentation:** Swagger/OpenAPI
- **Logging:** Logrus
- **Containerization:** Docker & Docker Compose

## Documentation

- **Swagger UI:** http://localhost:8080/docs/
- **Health Check:** http://localhost:8080/health
- **API Documentation:** Available in Swagger UI

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.


