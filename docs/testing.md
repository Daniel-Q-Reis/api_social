# Testing

This project includes both unit tests and integration tests to ensure the quality and correctness of the code.

## Unit Tests

Unit tests are located in the `internal/usecase` directory alongside the code they test. They use mocks to isolate the business logic from external dependencies.

### Running Unit Tests

To run all unit tests:

```bash
go test ./internal/usecase/...
```

To run unit tests with coverage:

```bash
go test -cover ./internal/usecase/...
```

## Integration Tests

Integration tests verify that the application works correctly with external dependencies, particularly the PostgreSQL database. They are located in the `integration-test` directory.

### Running Integration Tests

1. Start the test database:
   ```bash
   docker-compose up -d db
   ```

2. Run the integration tests:
   ```bash
   go test ./integration-test/...
   ```

Alternatively, you can run the integration tests using Docker:

```bash
docker-compose -f docker-compose-integration-test.yml up --build
```

### Test Database

The integration tests use the same database schema as the main application. The tests automatically clean up data before each test run to ensure isolation.

### Environment Variables

The integration tests use the following environment variables:

- `PG_URL`: Database connection string (default: `postgres://postgres:password@localhost:5432/social_db?sslmode=disable`)
- `SKIP_INTEGRATION_TESTS`: Set to `true` to skip integration tests

## Test Structure

- `internal/usecase/*_test.go`: Unit tests for business logic
- `integration-test/postgres_test.go`: Integration tests for user and post repositories
- `integration-test/postgres_interaction_test.go`: Integration tests for comment, like, and follow repositories