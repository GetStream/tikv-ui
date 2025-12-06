# TiKV UI

A REST API server for exploring and managing TiKV key-value data with multi-cluster support.

## Project Structure

```
tikv-ui/
├── cmd/
│   └── tikv-ui/
│       └── main.go              # Application entry point
├── pkg/
│   ├── handlers/                # HTTP request handlers
│   │   ├── health.go           # Health check endpoint
│   │   ├── cluster.go          # Cluster management
│   │   ├── get.go              # GET operation
│   │   ├── put.go              # PUT operation
│   │   ├── delete.go           # DELETE operation
│   │   └── scan.go             # SCAN operation
│   ├── server/                  # Server setup and middleware
│   │   ├── server.go           # Server struct with multi-cluster support
│   │   └── middleware.go       # HTTP middleware (logging, etc.)
│   ├── types/                   # Type definitions
│   │   ├── requests.go         # API request types
│   │   └── responses.go        # API response types
│   └── utils/                   # Utility functions
│       ├── http.go             # HTTP helpers (JSON responses, errors)
│       ├── msgpack.go          # Msgpack/JSON parsing
│       └── string.go           # String manipulation utilities
├── go.mod
├── go.sum
└── README.md
```

## Building

```bash
go build -o bin/tikv-ui ./cmd/tikv-ui
```

## Running

Set the `TIKV_PD_ADDRS` environment variable with comma-separated PD addresses for the default cluster:

```bash
export TIKV_PD_ADDRS="127.0.0.1:2379,127.0.0.1:23790"
./bin/tikv-ui
```

The server will start on port 8081 and connect to the default cluster.

## API Endpoints

### Health Check
```
GET /health
```

### Cluster Management

#### Connect to New Cluster
```
POST /api/clusters/connect
Body: {
  "pd_addrs": ["127.0.0.1:2379", "127.0.0.1:23790"],
  "name": "production"  // optional, auto-generated if not provided
}
```

#### List All Clusters
```
GET /api/clusters
Response: {
  "clusters": [
    {
      "name": "default",
      "cluster_id": 123456,
      "pd_addrs": ["127.0.0.1:2379"],
      "active": true
    }
  ]
}
```

#### Switch Active Cluster
```
POST /api/clusters/switch
Body: {"name": "production"}
```

### Raw KV Operations

All operations use the currently active cluster.

#### Get Value
```
POST /api/raw/get
Body: {"key": "mykey"}
Response: {
  "key": "mykey",
  "value": {...},        // parsed msgpack/JSON
  "raw_value": "...",    // raw string
  "found": true
}
```

#### Put Value
```
POST /api/raw/put
Body: {"key": "mykey", "value": "myvalue"}
```

#### Delete Key
```
POST /api/raw/delete
Body: {"key": "mykey"}
```

#### Scan Range
```
POST /api/raw/scan
Body: {"start_key": "a", "end_key": "z", "limit": 100}
Response: {
  "items": [
    {
      "key": "key1",
      "value": {...},      // parsed msgpack/JSON
      "raw_value": "..."   // raw string
    }
  ]
}
```

## Features

### Multi-Cluster Support
- Connect to multiple TiKV clusters dynamically via API
- Switch between clusters without restarting the server
- Each cluster maintains its own connection pool
- Thread-safe cluster management

### Smart Value Parsing
- Automatically detects and parses msgpack-encoded values
- Falls back to JSON parsing if msgpack fails
- Returns both parsed (structured) and raw (string) values
- Handles plain text values correctly

### Graceful Shutdown
- Properly closes all cluster connections on shutdown
- Handles SIGINT and SIGTERM signals

## Architecture

The project follows a clean architecture pattern:

- **cmd/**: Application entry points (main packages)
- **pkg/handlers/**: HTTP handlers organized by operation
- **pkg/server/**: Server configuration, middleware, and multi-cluster management
- **pkg/types/**: Shared type definitions
- **pkg/utils/**: Reusable utility functions

This structure provides:
- Clear separation of concerns
- Easy testing (each handler can be tested independently)
- Maintainability (changes are localized to specific packages)
- Scalability (easy to add new handlers or utilities)
