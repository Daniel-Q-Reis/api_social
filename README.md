# Social Media API

A RESTful API for a social media platform built with Go, following clean architecture principles.

## Features

- User registration and authentication with JWT
- CRUD operations for posts
- Follow/unfollow users
- Like/unlike posts
- Comment on posts
- Personalized feed
- User search

## Tech Stack

- **Language**: Go
- **Framework**: Chi router
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Database Migration**: golang-migrate

## Project Structure

```
.
├── cmd/app/main.go          # Application entry point
├── config/                  # Configuration files
├── docs/                    # Documentation
├── internal/
│   ├── entity/              # Domain models
│   ├── repo/                # Repository interfaces
│   ├── repo/postgres/       # PostgreSQL implementations
│   ├── usecase/             # Business logic
│   └── controller/
│       └── http/
│           ├── middleware/  # HTTP middleware
│           └── v1/          # HTTP handlers
├── migrations/              # Database migrations
└── go.mod, go.sum           # Go modules
```

## API Endpoints

### Authentication

- `POST /register` - Register a new user
- `POST /login` - Login and get JWT token

### Users & Profiles

- `GET /users/{username}` - Get user profile
- `GET /users/search?q={query}` - Search for users
- `GET /profile` - Get own profile (authenticated)
- `PUT /profile` - Update own profile (authenticated)

### Following

- `POST /users/{username}/follow` - Follow a user (authenticated)
- `DELETE /users/{username}/follow` - Unfollow a user (authenticated)
- `GET /users/{username}/followers` - Get user's followers
- `GET /users/{username}/following` - Get users someone is following

### Posts (Feed)

- `POST /posts` - Create a new post (authenticated)
- `GET /feed` - Get personalized feed (authenticated)
- `GET /posts/{postID}` - Get a single post
- `GET /users/{username}/posts` - Get all posts from a user
- `PUT /posts/{postID}` - Update a post (authenticated)
- `DELETE /posts/{postID}` - Delete a post (authenticated)

### Likes

- `POST /posts/{postID}/like` - Like a post (authenticated)
- `DELETE /posts/{postID}/like` - Unlike a post (authenticated)

### Comments

- `POST /posts/{postID}/comments` - Add a comment (authenticated)
- `GET /posts/{postID}/comments` - Get comments for a post
- `DELETE /posts/{postID}/comments/{commentID}` - Delete a comment (authenticated)

## Setup

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Set environment variables:
   - `PG_URL` - PostgreSQL connection string
   - `JWT_SECRET` - Secret for JWT signing
4. Run database migrations
5. Run the application: `go run cmd/app/main.go`

## Environment Variables

- `PG_URL` - PostgreSQL connection string (required)
- `JWT_SECRET` - Secret for JWT signing (required)
- `PORT` - Server port (default: 8080)

## Database Schema

The database schema is defined in the migrations directory. Run migrations to set up the database.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a pull request