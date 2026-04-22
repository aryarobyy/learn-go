# Auth Service - Project Description

This is a comprehensive authentication service built with Go that provides secure user authentication and authorization functionality. The service implements modern security practices including JWT tokens, refresh token rotation, and OAuth integration.

## Features

### 1. User Management
- **User Registration**: Create new user accounts with name, username, and password
- **User Login**: Authenticate users with username and password
- **User Retrieval**: Get user information by ID or username
- **User Listing**: Retrieve multiple users with pagination support

### 2. Authentication & Authorization
- **JWT-based Authentication**: Generate and validate access tokens with configurable expiration
- **Refresh Token System**: Long-lived refresh tokens with rotation and reuse detection
- **Token Refresh**: Ability to refresh access tokens using valid refresh tokens
- **Password Hashing**: Secure password storage using bcrypt

### 3. OAuth Integration
- **Google OAuth**: Integration with Google's authentication system
- **OAuth Redirect**: Handles the OAuth flow with Google's authentication server

### 4. Security Features
- **Redis-based Token Management**: Store and validate refresh tokens in Redis
- **Token Reuse Detection**: Detect and prevent refresh token reuse attacks
- **Session Expiration**: Automatic session expiration after 30 days of inactivity
- **Rate Limiting**: Built-in rate limiting middleware to prevent abuse

### 5. Infrastructure & Deployment
- **Docker Support**: Containerized deployment with Docker and docker-compose
- **PostgreSQL Database**: Persistent storage for user data
- **Redis Cache**: For session and token management
- **Database Migrations**: Schema management with migration support

## API Endpoints

### Authentication Routes (`/auth`)
- `POST /auth` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Token refresh

### User Routes (`/user`)
- `GET /user` - Get multiple users (with pagination)
- `GET /user/{id}` - Get user by ID

### OAuth Routes
- `GET /google_login` - Google OAuth redirect

## Architecture

The project follows a clean architecture pattern with the following layers:

- **App Layer**: Main application entry point and server configuration
- **Config Layer**: Environment configuration and Google OAuth setup
- **Controller Layer**: HTTP request handling and response formatting
- **Service Layer**: Business logic implementation
- **Repository Layer**: Database operations
- **Helper Layer**: Utility functions (JWT handling, password hashing, etc.)
- **Middleware Layer**: Cross-cutting concerns (rate limiting, etc.)

## Security Implementation

### Token Management
- Access tokens with configurable expiration (default 30 minutes)
- Refresh tokens with rotation and reuse detection
- Redis-based token validation to enable revocation
- Session expiration after 30 days of inactivity

### Password Security
- Passwords are hashed using bcrypt with default cost
- Input validation for user names and other fields

### Rate Limiting
- Built-in rate limiting to prevent brute force attacks

## Technology Stack

- **Language**: Go (Golang)
- **Web Framework**: chi router
- **Database**: PostgreSQL
- **Cache**: Redis
- **Authentication**: JWT tokens
- **OAuth**: Google OAuth 2.0
- **Containerization**: Docker and Docker Compose
- **Database Migrations**: migrate/migrate

## Configuration

The service uses environment variables for configuration:

- `JWT_SECRET` - Secret key for signing access tokens
- `JWT_EXPIRED` - Access token expiration time (default: 30m)
- `JWT_REFRESH_SECRET` - Secret key for signing refresh tokens
- `JWT_REFRESH_EXPIRED` - Refresh token expiration time (default: 7d)
- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret
- `REDIS_ADDR` - Redis server address (default: localhost:6379)
- `PORT` - Server port (default: 8080)

## Deployment

The service can be deployed using Docker Compose which includes:
- PostgreSQL database container
- Redis cache container
- Migration container for database schema updates