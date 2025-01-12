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

   go build -o logging-ms.exe main.go
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
   docker volume create logger-data
   docker run -d \
      --name logging-service \
      -p 8080:80 \
      -v logger-data:/app/logs.db \
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

---

## API Endpoints

1. **POST /logs**
   **Purpose:** Create a new log entry.

   - **Request Body** (JSON):

   ```json
   {
     "timestamp": "2025-01-12T00:00:00Z",  // optional, defaults to current time
     "logger": "myLogger",
     "level": "info",
     "message": "This is a test message",
     "meta": { "user": "Alice", "action": "testing" }
   }
   ```

   **Response** (JSON):

   ```json
   {
      "id": 1,
      "timestamp": "2025-01-12T00:00:00Z",
      "logger": "myLogger",
      "logger_id": 1,
      "level": "info",
      "message": "This is a test message",
      "meta": {
         "user": "Alice",
         "action": "testing"
      }
   }
   ```

2. **GET /logs**
   **Purpose:** Retrieve a list of logs with optional filters.

   - **Query Params:**
      - mintime, maxtime (e.g. ?mintime=2025-01-01&maxtime=2025-01-02)
      - search (e.g. ?search=message)
      - includeLoggers, excludeLoggers (repeatable, e.g. ?includeLoggers=app1&includeLoggers=app2)
      - offset, limit (e.g. ?offset=10&limit=50)
   - Response (JSON array):

   ```json
   [
      {
         "id": 1,
         "timestamp": "2025-01-12T00:00:00Z",
         "logger": "myLogger",
         "logger_id": 1,
         "level": "info",
         "message": "This is a test message"
      }
   ]
   ```
3. **GET /logs/{id}**
   **Purpose:** Retrieve a single log by ID.

   **Example:** GET /logs/1
   **Response** (JSON):

   ```json
   {
      "id": 1,
      "timestamp": "2025-01-12T00:00:00Z",
      "logger": "myLogger",
      "logger_id": 1,
      "level": "info",
      "message": "This is a test message",
      "meta": {
         "user": "Alice",
         "action": "testing"
      }
   }
   ```
4. **POST /config**
**Purpose:** Update a logger’s level.

   **Request Body** (JSON):

   ```json
   {
      "name": "myLogger",
      "level": "error"
   }
   ```
   This sets myLogger’s level to error. Logs below error (like info, warn) will be ignored for that logger.

   **Response:**

   ```text
   Logger level updated successfully
   ```

---

## Example Usage

1. **Create a log:**

   ```bash
   curl -X POST http://localhost:8080/logs \
      -H "Content-Type: application/json" \
      -d '{
      "logger": "appLogger",
      "level": "info",
      "message": "Just testing logs",
      "meta": { "username": "bob" }
      }'
   ```

2. **View Logs:**

   ```bash
   curl -X GET "http://localhost:8080/logs?search=testing&limit=10"
   ```

3. **Update logger level:**

   ```bash
   curl -X POST http://localhost:8080/config \
      -H "Content-Type: application/json" \
      -d '{
      "name": "appLogger",
      "level": "warn"
      }'
   ```

