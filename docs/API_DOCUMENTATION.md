# API Documentation

This document describes all available REST and GraphQL endpoints for the platform.

## Base URL
`http://localhost:8080`

---

## REST API Endpoints

### Audiences

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/audience` | Create a new audience |
| GET | `/audiences` | Get all audiences |
| GET | `/audience/:id` | Get audience by ID |
| PUT | `/audience/:id` | Update audience by ID |
| DELETE | `/audience/:id` | Delete audience by ID |

**Audience Model:**
```json
{
  "id": 1,
  "gender": "Male",
  "birthcountry": "USA",
  "agegroup": "25-34",
  "dailyhours": 5,
  "noofpurchases": 10
}
```

### Charts

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/chart` | Create a new chart |
| GET | `/charts` | Get all charts |
| GET | `/chart/:id` | Get chart by ID |
| PUT | `/chart/:id` | Update chart by ID |
| DELETE | `/chart/:id` | Delete chart by ID |

**Chart Model:**
```json
{
  "id": 1,
  "title": "Sales Chart",
  "xaxistitle": "Months",
  "yaxistitle": "Revenue"
}
```

### Insights

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/insight` | Create a new insight |
| GET | `/insights` | Get all insights |
| GET | `/insight/:id` | Get insight by ID |
| PUT | `/insight/:id` | Update insight by ID |
| DELETE | `/insight/:id` | Delete insight by ID |

**Insight Model:**
```json
{
  "id": 1,
  "text": "Sales increased by 20% in Q4"
}
```

### User Favourites

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/userstar` | Create a new user favourite |
| GET | `/userstars` | Get all user favourites |
| GET | `/userstar/:id` | Get user favourite by ID |
| GET | `/userstars/user/:userid` | Get all favourites for a specific user |
| PUT | `/userstar/:id` | Update user favourite by ID |
| DELETE | `/userstar/:id` | Delete user favourite by ID |

**UserStar Model:**
```json
{
  "id": 1,
  "userid": 123,
  "type": "chart",
  "assetid": 456
}
```

**Type field** can be: `"audience"`, `"chart"`, or `"insight"`

---

## GraphQL API

**Endpoint:** `POST /graphql`
**Playground:** `GET /graphql`

### Queries

#### Audiences
```graphql
# Get all audiences
query {
  audiences {
    id
    gender
    birthcountry
    agegroup
    dailyhours
    noofpurchases
  }
}

# Get audience by ID
query {
  audience(id: "1") {
    id
    gender
    birthcountry
  }
}
```

#### Charts
```graphql
# Get all charts
query {
  charts {
    id
    title
    xaxistitle
    yaxistitle
  }
}

# Get chart by ID
query {
  chart(id: "1") {
    id
    title
  }
}
```

#### Insights
```graphql
# Get all insights
query {
  insights {
    id
    text
  }
}

# Get insight by ID
query {
  insight(id: "1") {
    id
    text
  }
}
```

#### User Favourites
```graphql
# Get all user favourites
query {
  userstars {
    id
    userid
    type
    assetid
  }
}

# Get user favourite by ID
query {
  userstar(id: "1") {
    id
    userid
    type
    assetid
  }
}

# Get favourites by user ID
query {
  userstarsByUser(userid: 123) {
    id
    type
    assetid
  }
}
```

### Mutations

#### Audiences
```graphql
# Create audience
mutation {
  createAudience(input: {
    gender: "Male"
    birthcountry: "USA"
    agegroup: "25-34"
    dailyhours: 5
    noofpurchases: 10
  }) {
    id
    gender
  }
}

# Update audience
mutation {
  updateAudience(id: "1", input: {
    dailyhours: 6
  }) {
    id
    dailyhours
  }
}

# Delete audience
mutation {
  deleteAudience(id: "1")
}
```

#### Charts
```graphql
# Create chart
mutation {
  createChart(input: {
    title: "Sales Chart"
    xaxistitle: "Months"
    yaxistitle: "Revenue"
  }) {
    id
    title
  }
}

# Update chart
mutation {
  updateChart(id: "1", input: {
    title: "Updated Sales Chart"
  }) {
    id
    title
  }
}

# Delete chart
mutation {
  deleteChart(id: "1")
}
```

#### Insights
```graphql
# Create insight
mutation {
  createInsight(input: {
    text: "Sales increased by 20% in Q4"
  }) {
    id
    text
  }
}

# Update insight
mutation {
  updateInsight(id: "1", input: {
    text: "Sales increased by 25% in Q4"
  }) {
    id
    text
  }
}

# Delete insight
mutation {
  deleteInsight(id: "1")
}
```

#### User Favourites
```graphql
# Create user favourite
mutation {
  createUserStar(input: {
    userid: 123
    type: "chart"
    assetid: 456
  }) {
    id
    userid
    type
    assetid
  }
}

# Update user favourite
mutation {
  updateUserStar(id: "1", input: {
    type: "insight"
    assetid: 789
  }) {
    id
    type
    assetid
  }
}

# Delete user favourite
mutation {
  deleteUserStar(id: "1")
}
```

---

## Data Relationships

The `UserStar` model allows users to favourite different types of assets:

- **Type:** "audience" → References an Audience (via `assetid`)
- **Type:** "chart" → References a Chart (via `assetid`)
- **Type:** "insight" → References an Insight (via `assetid`)

Example: A user can favourite a specific chart by creating:
```json
{
  "userid": 123,
  "type": "chart",
  "assetid": 456
}
```

Where `456` is the ID of the chart they want to favourite.

---

## Response Format (REST)

All REST endpoints return responses in this format:

```json
{
  "status": 200,
  "message": "Success message",
  "data": { /* response data */ }
}
```

Error responses:
```json
{
  "status": 400,
  "message": "Error message",
  "data": null
}
```

---

## GraphQL Features

- **Caching:** Automatic Persisted Queries (APQ) enabled
- **Playground:** Interactive GraphQL IDE at `/graphql`
- **DataLoader:** Batch loading to prevent N+1 queries
- **Introspection:** Full schema introspection available
