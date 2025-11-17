# GraphQL Caching Layer

Your GraphQL API now includes a comprehensive caching layer with the following features:

## 1. DataLoader (N+1 Query Prevention)

DataLoader batches and caches database requests within a single GraphQL operation to prevent N+1 queries.

**How it works:**
- When multiple resolvers request the same audience ID within one query, DataLoader batches them into a single database call
- Results are cached for the duration of the request
- Configured in [graph/dataloader.go](graph/dataloader.go)

**Example scenario where DataLoader helps:**
```graphql
query {
  audiences {
    id
    gender
    birthcountry
    agegroup
  }
}
```

If you later add relationships (e.g., audiences with purchases), DataLoader prevents making N database queries by batching them.

## 2. Automatic Persisted Queries (APQ)

APQ reduces bandwidth by caching query strings on the server.

**How it works:**
- Client sends a SHA-256 hash of the query instead of the full query string
- Server caches up to 100 queries (configurable in [main.go](main.go#L30))
- On first request, full query is sent and cached
- Subsequent requests only send the hash

**Benefits:**
- Reduces payload size
- Faster query parsing
- Better performance for repeated queries

## 3. Configuration

All caching is configured in [main.go](main.go):

```go
// DataLoader middleware - batches DB queries
h.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
    ctx = graph.DataLoaderMiddleware(api.DB)(ctx)
    return next(ctx)
})

// APQ - caches up to 100 query strings
h.Use(extension.AutomaticPersistedQuery{
    Cache: lru.New[string](100),
})
```

## 4. Cache Settings

You can adjust these settings in [main.go](main.go):

- **APQ Cache Size:** Change `lru.New[string](100)` to cache more/fewer queries
- **DataLoader Batch Size:** Modify `WithBatchCapacity` in [graph/dataloader.go](graph/dataloader.go#L74)
- **DataLoader Wait Time:** Adjust `WithWait` in [graph/dataloader.go](graph/dataloader.go#L75)

## 5. Testing the Cache

You can test caching behavior by:

1. Running the same query multiple times
2. Checking server logs for database query counts
3. Using browser DevTools to see APQ hashes in network requests

## 6. Cache Behavior

**DataLoader:**
- Cache lifetime: Single GraphQL request
- Scope: Per request
- Clears after each operation

**APQ:**
- Cache lifetime: Until server restart or cache eviction (LRU)
- Scope: Server-wide
- Persists across requests
