# Go HTTP Server

A simple and efficient HTTP server built with Go.

## Features
- RESTful API endpoints
- JSON request/response handling
- Graceful shutdown
- Configurable port
- Health check endpoint

## Getting Started

### Prerequisites
- Go 1.19 or higher

### Installation

1. Clone the repository:
```bash
git clone sy-GO-nergy
cd http/
```

2. Install Dependencies

``` bash
go mod tidy
```

3. Run the server
```bash
go run cmd/server/main.go
```

API Endpoints
GET /health - Health check

GET /api/v1/users - Get all users

GET /api/v1/users/{id} - Get user by ID

POST /api/v1/users - Create new user

Configuration
Server configuration can be modified in configs/config.yaml:

```yaml
server:
  port: 8080
  read_timeout: 15
  write_timeout: 15
```

Development
To build the application:

```bash
go build -o bin/server cmd/server/main.go
```
License
MIT
