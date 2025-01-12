# Logging Microservice

A lightweight logging microservice written in Go, using SQLite for data persistence.

## Features

- Create and retrieve log entries  
- Flexible queries (filters, search, includes/excludes)  
- Dynamically set logger levels to control which messages are persisted  
- Simple RESTful endpoints  
- Easy Docker deployment with persistent storage  

---

## Table of Contents

1. [Prerequisites](#prerequisites)  
2. [Project Structure](#project-structure)  
3. [Setup and Run Locally](#setup-and-run-locally)  
4. [Environment Variables](#environment-variables)  
5. [Docker Instructions](#docker-instructions)  
6. [API Endpoints](#api-endpoints)  
7. [Example Usage](#example-usage)  
8. [License](#license)

---

## Prerequisites

- Go 1.20+ (with CGO enabled)  
- SQLite library (C compiler needed if building from source; see [CGO requirements](https://github.com/mattn/go-sqlite3/blob/master/README.md))  

If running in Docker, ensure you have [Docker](https://docs.docker.com/get-docker/) installed.

---

## Project Structure

```bash
logging-microservice/
 ┣ db/
 ┃  ┣ database.go
 ┃  ┗ logs.go
 ┣ handlers/
 ┃  ┣ logs.go
 ┃  ┗ config.go
 ┣ models/
 ┃  ┣ logger.go
 ┃  ┗ logentry.go
 ┣ main.go
 ┣ go.mod
 ┣ go.sum
 ┗ Dockerfile
```

- db/: Database logic (initialization, queries, and CRUD operations).
- handlers/: HTTP handlers for logs and config endpoints.
- models/: Data structures (Logger, LogEntry).
- main.go: Application entry point (server and route setup).
- Dockerfile: Multi-stage Docker build file.

---

## Setup and Run Locally

1. **Clone the repository**:

   ```bash
   git clone https://github.com/yourusername/logging-microservice.git
   cd logging-microservice
   ```

2. **(Optional) Update go.mod to match your module path:**:

   ```bash
   go mod tidy
   ```

3. **Build and run**:

   ```bash
   # Make sure CGO is enabled
   export CGO_ENABLED=1

   go build -o logging-ms main.go
   ./logging-ms
   ```

4. **Confirm the service is running at http://localhost:8080.**

---

## Environment Variables

| Variable | Default | Description | 
|----------|---------|-------------|
| PORT     | 8080    | The port on which the server listens |

---

## Docker Instructions

1. **Build the Image**:
   ```bash
   docker build -t logging-ms:latest .
   ```
   **Important:** If you’re using go-sqlite3, ensure your Dockerfile has CGO enabled and a C toolchain (e.g. Alpine’s build-base) in the build stage:
   ```dockerfile
   FROM golang:1.20-alpine AS builder
   RUN apk add --no-cache build-base
   ENV CGO_ENABLED=1
   # ...
   ```

2. **Run with a Persistent Volume**:
   ```bash
   docker run -d \
     --name logging-ms \
     -p 8080:8080 \
     -v $(pwd)/data:/app \
     -e DB_PATH=/app/logs.db \
     logging-ms:latest
   ```
   **Explanation:**

   - -p 8080:8080: Publishes port 8080 from the container to host machine.
   - -v $(pwd)/data:/app: Mounts a local directory ./data into /app in the container, so the SQLite file (logs.db) persists on your host machine.
   - -e DB_PATH=/app/logs.db: Tells the service to store the SQLite DB at /app/logs.db.
   ### Passing Environment Vars at Runtime
   If you want to override the port:
   ```bash
   docker run -d \
     --name logging-ms \
     -e PORT=9000 \
     -e DB_PATH=/app/logs.db \
     -p 9000:9000 \
     -v $(pwd)/data:/app \
     logging-ms:latest
   ```

