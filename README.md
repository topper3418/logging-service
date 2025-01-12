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
