# Testing Guide

This document provides comprehensive instructions for running E2E (End-to-End) functional tests and performance tests for the platform-go-challenge GraphQL API.

## Table of Contents

1. [Overview](#overview)
2. [Setup](#setup)
3. [E2E Tests](#e2e-tests)
4. [Performance Tests](#performance-tests)

---

## Overview

The testing suite includes:

- **E2E Tests**: Functional tests that validate the GraphQL API behavior with a real database
- **Benchmark Tests**: Go benchmarks that measure query performance and memory allocation

### Test Structure

```
tests/
├── e2e/                      # End-to-end functional tests
│   ├── setup_test.go         # Test infrastructure and helpers
│   └── userinterface_test.go # UserInterface query tests
├── performance/              # Performance benchmarks
│   └── userinterface_bench_test.go
└── README.md                 # This file
```

---

## Setup

### Prerequisites

1. **Go 1.21+** installed
2. **PostgreSQL** database running
3. **Test Database** created (separate from production)

### Environment Setup

1. Create a test database (if you want to use a separate test database):

```bash
createdb mydb_test
```

2. Configure database connection in `.env` file at the project root:

```bash
# .env file
DB_URL=postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable
TEST_DB_URL=postgres://myuser:mypassword@localhost:5432/mydb_test?sslmode=disable
```

**Note:** Tests will automatically load the `.env` file from the project root. If `TEST_DB_URL` is not set, tests will fall back to using `DB_URL`.

3. Install test dependencies:

```bash
go mod download
```

---

## E2E Tests

### Running E2E Tests

E2E tests validate the complete GraphQL query flow including database interactions.

#### Run all E2E tests:

```bash
go test ./tests/e2e/... -v
```

#### Run specific test:

```bash
go test ./tests/e2e/... -v -run TestUserInterface_WithFavourites
```

#### Run with verbose output:

```bash
go test ./tests/e2e/... -v -count=1
```

### Test Coverage

Current E2E tests cover:

- ✅ Empty favourites (user has no favourites)
- ✅ Single favourite per asset type
- ✅ Multiple favourites of the same type
- ✅ User isolation (only fetching specific user's data)
- ✅ Invalid user ID error handling

### E2E Test Output Example

```
=== RUN   TestUserInterface_EmptyFavourites
--- PASS: TestUserInterface_EmptyFavourites (0.02s)
=== RUN   TestUserInterface_WithFavourites
--- PASS: TestUserInterface_WithFavourites (0.05s)
=== RUN   TestUserInterface_MultipleFavouritesOfSameType
--- PASS: TestUserInterface_MultipleFavouritesOfSameType (0.04s)
=== RUN   TestUserInterface_OnlySpecificUser
--- PASS: TestUserInterface_OnlySpecificUser (0.03s)
=== RUN   TestUserInterface_InvalidUserID
--- PASS: TestUserInterface_InvalidUserID (0.01s)
PASS
```

### Adding New E2E Tests

To add new E2E tests:

1. Create a new test function in `tests/e2e/userinterface_test.go`
2. Use the helper functions:
   - `ExecuteGraphQL(t, query, variables)` - Execute GraphQL queries
   - `CleanupTestData(testDB)` - Clean database before test
   - `SeedTestData(t, testDB)` - Create sample data
3. Follow the AAA pattern: Arrange, Act, Assert

Example:

```go
func TestUserInterface_YourNewTest(t *testing.T) {
    // Arrange
    CleanupTestData(testDB)
    audienceID, _, _ := SeedTestData(t, testDB)

    // Act
    query := `...`
    resp := ExecuteGraphQL(t, query, variables)

    // Assert
    if len(resp.Errors) > 0 {
        t.Fatalf("expected no errors, got: %v", resp.Errors)
    }
}
```

---

## Performance Tests

### Running Benchmark Tests

Benchmark tests measure query performance under different data loads.

#### Run all benchmarks:

```bash
go test ./tests/performance/... -bench=. -benchmem
```

#### Run specific benchmark:

```bash
go test ./tests/performance/... -bench=BenchmarkUserInterface_SmallDataset -benchmem
```

#### Run with custom duration:

```bash
go test ./tests/performance/... -bench=. -benchtime=10s -benchmem
```

#### Run and save results for comparison:

```bash
go test ./tests/performance/... -bench=. -benchmem | tee bench-results.txt
```

### Available Benchmarks

| Benchmark | Description | Data Size |
|-----------|-------------|-----------|
| `BenchmarkUserInterface_SmallDataset` | 5 items per user | Small |
| `BenchmarkUserInterface_MediumDataset` | 50 items per user | Medium |
| `BenchmarkUserInterface_LargeDataset` | 200 items per user | Large |
| `BenchmarkUserInterface_MinimalFields` | Only IDs queried | 50 items |
| `BenchmarkUserInterface_AllFields` | All fields queried | 50 items |
| `BenchmarkUserInterface_Parallel` | Concurrent requests | 50 items/user, 100 users |

### Benchmark Output Example

```
BenchmarkUserInterface_SmallDataset-8         1000    1234567 ns/op    123456 B/op    1234 allocs/op
BenchmarkUserInterface_MediumDataset-8         500    2345678 ns/op    234567 B/op    2345 allocs/op
BenchmarkUserInterface_LargeDataset-8          200    5678901 ns/op    456789 B/op    3456 allocs/op
```

**Reading the output:**
- `1000` - Number of iterations
- `1234567 ns/op` - Nanoseconds per operation
- `123456 B/op` - Bytes allocated per operation
- `1234 allocs/op` - Number of allocations per operation

### Comparing Benchmarks

To compare performance before and after changes:

```bash
# Before changes
go test ./tests/performance/... -bench=. -benchmem > old.txt

# After changes
go test ./tests/performance/... -bench=. -benchmem > new.txt

# Compare (requires benchstat tool)
go install golang.org/x/perf/cmd/benchstat@latest
benchstat old.txt new.txt
```

---

## Performance Optimization Tips

Based on test results, consider these optimizations:

### 1. Database Indexing

Ensure proper indexes exist:

```sql
CREATE INDEX idx_user_favourites_user_id ON user_favourites(user_id);
CREATE INDEX idx_user_favourites_asset_id ON user_favourites(asset_id);
CREATE INDEX idx_user_favourites_type ON user_favourites(type);
```

### 2. Query Optimization

Use `EXPLAIN ANALYZE` to identify slow queries:

```sql
EXPLAIN ANALYZE
SELECT * FROM user_favourites WHERE user_id = 1;
```

### 3. Connection Pooling

Configure GORM connection pool in `db/db.go`:

```go
sqlDB, err := GormDB.DB()
if err == nil {
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
}
```

### 4. Caching

Consider adding Redis caching for frequently accessed data:

```go
// Pseudo-code
func (r *queryResolver) Userinterface(ctx context.Context, userID string) (*model.UserInterface, error) {
    // Try cache first
    if cached, found := redis.Get("userinterface:" + userID); found {
        return cached, nil
    }

    // Query database
    result := // ... existing logic

    // Cache result
    redis.Set("userinterface:" + userID, result, 5*time.Minute)

    return result, nil
}
```

### 5. Batch Loading

Implement DataLoader pattern to prevent N+1 queries if querying multiple users.

---

## Troubleshooting

### Common Issues

**Issue:** Tests fail with "connection refused"
**Solution:** Ensure PostgreSQL is running and TEST_DB_URL is correct

**Issue:** Benchmark results are inconsistent
**Solution:** Run benchmarks multiple times and use longer `-benchtime` duration

**Issue:** "database does not exist"
**Solution:** Create the test database: `createdb mydb_test`

---

## Best Practices

1. **Run tests before committing:** Always run E2E tests before pushing code
2. **Isolate test data:** Use a separate test database, never test against production
3. **Clean between tests:** Use `CleanupTestData()` to ensure test isolation
4. **Monitor performance trends:** Track benchmark results over time to catch regressions
5. **Set appropriate thresholds:** Define SLAs based on business requirements

---

## Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [GORM Documentation](https://gorm.io/docs/)
- [gqlgen Documentation](https://gqlgen.com/)

---

## Contributing

When adding new features:

1. Write E2E tests for new GraphQL queries/mutations
2. Add benchmarks if performance-critical
3. Update this README with new test information
4. Ensure all tests pass before creating PR

---

## Questions?

If you have questions about the testing infrastructure, please open an issue or contact the development team.
