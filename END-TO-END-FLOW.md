# WSInspect - WebSocket Debugging Tool
## End-to-End Flow Documentation

---

## Table of Contents
1. [Introduction](#introduction)
2. [What is WSInspect?](#what-is-wsinspect)
3. [Architecture Overview](#architecture-overview)
4. [Prerequisites](#prerequisites)
5. [Installation & Setup](#installation--setup)
6. [Running the Application](#running-the-application)
7. [Frontend-Backend Communication](#frontend-backend-communication)
8. [API Endpoints Reference](#api-endpoints-reference)
9. [Go Concepts Explained](#go-concepts-explained)
10. [Testing Process](#testing-process)
11. [Troubleshooting](#troubleshooting)

---

## 1. Introduction

This document provides a comprehensive guide to understanding and using WSInspect - a WebSocket debugging tool similar to Postman but specifically designed for WebSocket APIs.

---

## 2. What is WSInspect?

WSInspect is a powerful WebSocket debugging and testing tool that allows developers to:

- **Inspect WebSocket traffic** - Capture and analyze messages between client and server
- **Create proxies** - Act as a middleman between client and WebSocket server
- **Test and replay messages** - Replay captured messages to test server behavior
- **Fuzz testing** - Automatically generate and test various message inputs
- **Protocol support** - Handle JSON, binary, and custom payload formats

---

## 3. Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         WSInspect                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────┐          ┌──────────────────┐             │
│  │    Frontend      │          │     Backend       │             │
│  │   (React +       │◄────────►│    (Go + Gin)    │             │
│  │    Vite)        │   HTTP/   │                  │             │
│  │                  │   WebSocket│                 │             │
│  └──────────────────┘          └────────┬─────────┘             │
│                                          │                        │
│                                          ▼                        │
│                                 ┌──────────────────┐             │
│                                 │   WebSocket      │             │
│                                 │   Proxy          │             │
│                                 │   (Hub)          │             │
│                                 └────────┬─────────┘             │
│                                          │                        │
│                    ┌─────────────────────┼─────────────────────┐│
│                    ▼                     ▼                     ▼│
│            ┌──────────────┐      ┌──────────────┐      ┌───────┐│
│            │   Target     │      │   SQLite     │      │ File  ││
│            │   Server     │      │   Database   │      │ System││
│            └──────────────┘      └──────────────┘      └───────┘│
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Components:

1. **Frontend (React + Vite + Tailwind)**
   - User interface for interacting with the tool
   - Dashboard, Sessions, Connections, Fuzz Testing, Settings pages
   - WebSocket message editor with Monaco Editor

2. **Backend (Go + Gin)**
   - REST API for CRUD operations
   - WebSocket proxy for traffic inspection
   - Business logic services

3. **Database (SQLite)**
   - Stores sessions, messages, fuzz tests
   - Persistent storage using GORM

---

## 4. Prerequisites

### For Backend (Go):
- Go 1.21+ installed
- Understanding of basic Go syntax

### For Frontend:
- Node.js 18+ installed
- npm or yarn

### For Database:
- SQLite (included, no setup needed)

---

## 5. Installation & Setup

### Step 1: Clone the repository
```bash
git clone <repository-url>
cd "opencode projects/Websockets Debugging Tool"
```

### Step 2: Setup Backend

Navigate to the backend directory:
```bash
cd backend
```

Initialize Go modules:
```bash
go mod init wsinspect/backend
go mod tidy
```

The project uses these key dependencies:
- **Gin** - HTTP web framework (like Express.js in Node.js)
- **GORM** - ORM for database operations
- **SQLite** - Database driver
- **Cobra** - CLI command framework
- **Viper** - Configuration management
- **Gorilla WebSocket** - WebSocket handling

### Step 3: Setup Frontend

Navigate to the frontend directory:
```bash
cd frontend
```

Install dependencies:
```bash
npm install
```

---

## 6. Running the Application

### Option A: Using Docker (Recommended)

```bash
# Build and run all services
docker-compose up --build

# Access the application:
# - Frontend: http://localhost:5173
# - Backend API: http://localhost:8080
# - WebSocket: ws://localhost:8080
```

### Option B: Running Backend Manually

```bash
cd backend

# Run with default settings
go run main.go start

# Or with custom port
go run main.go start --port 9000

# Or with custom target
go run main.go start --port 8080 --target "wss://my-server.com/ws"
```

### Option C: Running Frontend Manually

```bash
cd frontend
npm run dev
```

The frontend will start on http://localhost:5173

---

## 7. Frontend-Backend Communication

### HTTP Communication Flow

```
┌─────────────┐    HTTP Request     ┌─────────────┐
│             │ ──────────────────► │             │
│   Frontend  │    GET /api/sessions│   Backend   │
│   (React)   │ ◄────────────────── │   (Gin)     │
│             │    JSON Response    │             │
└─────────────┘                     └─────────────┘
```

### WebSocket Communication Flow

```
┌─────────────┐    WS Connection    ┌─────────────┐
│             │ ───────────────────► │             │
│   Frontend  │                      │   Backend   │
│   (React)   │   Real-time Messages│   (Hub)     │
│             │ ◄──────────────────► │             │
│             │                      │             │
└─────────────┘                      └──────┬──────┘
                                             │
                                    ┌────────▼────────┐
                                    │ Target WebSocket│
                                    │    Server       │
                                    └─────────────────┘
```

### Vite Proxy Configuration

The frontend is configured to proxy API requests to the backend:

```typescript
// vite.config.ts
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true,
  },
  '/ws': {
    target: 'ws://localhost:8080',
    ws: true,
  },
}
```

This means:
- `http://localhost:5173/api/sessions` → `http://localhost:8080/api/sessions`
- `ws://localhost:5173/ws` → `ws://localhost:8080/ws`

---

## 8. API Endpoints Reference

### Base URL: `http://localhost:8080`

### Health Check

#### GET /health
Check if the server is running.

**Response:**
```json
{
  "status": "healthy",
  "service": "wsinspect"
}
```

**Frontend Usage:**
```typescript
const response = await fetch('/api/health');
const data = await response.json();
```

---

### Proxy Status

#### GET /api/proxy/status
Get the current status of the WebSocket proxy.

**Response:**
```json
{
  "status": "running",
  "active_connections": 5
}
```

**Frontend Usage:**
```typescript
// Dashboard.tsx
const response = await fetch('/api/proxy/status');
const data = await response.json();
setStats(prev => ({
  ...prev,
  activeConnections: data.active_connections || 0
}));
```

---

### Sessions

#### GET /api/sessions
List all WebSocket sessions.

**Query Parameters:**
- `limit` (optional): Number of sessions to return (default: 20)
- `offset` (optional): Number of sessions to skip (default: 0)

**Response:**
```json
{
  "sessions": [
    {
      "id": 1,
      "connection_id": "abc-123",
      "client_ip": "127.0.0.1",
      "server_host": "ws://localhost:3000",
      "status": "active",
      "message_count": 42,
      "start_time": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 1
}
```

**Frontend Usage:**
```typescript
// Sessions.tsx
useEffect(() => {
  fetch('/api/sessions')
    .then(res => res.json())
    .then(data => setSessions(data.sessions || []));
}, []);
```

---

#### GET /api/sessions/:id
Get details of a specific session.

**Response:**
```json
{
  "id": 1,
  "connection_id": "abc-123",
  "client_ip": "127.0.0.1",
  "server_host": "ws://localhost:3000",
  "status": "active",
  "message_count": 42,
  "start_time": "2024-01-15T10:30:00Z"
}
```

---

#### DELETE /api/sessions/:id
Delete a session.

**Response:**
```json
{
  "message": "Session deleted"
}
```

---

### Messages

#### GET /api/messages/session/:session_id
Get all messages for a specific session.

**Response:**
```json
{
  "messages": [
    {
      "id": 1,
      "session_id": 1,
      "direction": "client-to-server",
      "payload": "{\"type\": \"hello\"}",
      "payload_format": "json",
      "timestamp": "2024-01-15T10:30:00Z",
      "opcode": 1
    }
  ],
  "total": 1
}
```

---

#### POST /api/messages/inject
Inject a custom message into the WebSocket stream.

**Request Body:**
```json
{
  "session_id": 1,
  "payload": "{\"type\": \"ping\"}",
  "direction": "client-to-server"
}
```

**Response:**
```json
{
  "message": "Message injected"
}
```

---

### Replay

#### POST /api/replay
Replay a captured message or sequence of messages.

**Request Body:**
```json
{
  "session_id": 1,
  "message_ids": [1, 2, 3],
  "delay_ms": 100
}
```

**Response:**
```json
{
  "message": "Replay started"
}
```

---

### Fuzz Testing

#### GET /api/fuzz
List all fuzz tests.

**Response:**
```json
[
  {
    "id": 1,
    "name": "Test 1",
    "status": "completed"
  }
]
```

---

#### POST /api/fuzz
Create a new fuzz test.

**Request Body:**
```json
{
  "name": "My Fuzz Test",
  "template": "{\"type\": \"{{data}}\"}",
  "strategy": "random",
  "iterations": 100
}
```

**Response:**
```json
{
  "message": "Fuzz test created"
}
```

---

#### POST /api/fuzz/:id/run
Run a specific fuzz test.

**Response:**
```json
{
  "message": "Fuzz test running"
}
```

---

#### GET /api/fuzz/:id
Get details of a fuzz test.

**Response:**
```json
{
  "message": "Fuzz test details"
}
```

---

### WebSocket Proxy

#### GET /ws/*target
Connect to the WebSocket proxy for real-time traffic inspection.

**Query Parameters:**
- `host`: Target WebSocket server URL

**Example:**
```
ws://localhost:8080/ws?host=localhost:3000
```

**Frontend Usage:**
```typescript
const ws = new WebSocket('ws://localhost:8080/ws?host=localhost:3000');

ws.onopen = () => {
  console.log('Connected to WebSocket proxy');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};
```

---

## 9. Go Concepts Explained

### For Beginners - Go Basics

#### Packages and Imports

```go
package main  // Every Go file starts with a package declaration

import (
    "fmt"      // Formatted I/O
    "log"      // Logging
    "net/http" // HTTP client/server
)
```

**Key Concept:** In Go, everything is organized into packages. The `main` package makes the file executable.

---

#### Variables and Types

```go
// Explicit type declaration
var port int = 8080

// Type inference
var target = "ws://localhost:3000"

// Short variable declaration (inside functions)
name := "wsinspect"
```

**Key Concept:** Go is statically typed - variable types are known at compile time.

---

#### Functions

```go
// Basic function
func greet(name string) string {
    return "Hello, " + name
}

// Function with multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

---

#### Structs (Like Classes in OOP)

```go
// Define a struct
type Session struct {
    ID           uint
    ConnectionID string
    ClientIP     string
    ServerHost   string
    Status       string
}

// Create an instance
session := Session{
    ID:           1,
    ConnectionID: "abc-123",
    ClientIP:     "127.0.0.1",
    Status:       "active",
}
```

---

#### Methods (Attached to Structs)

```go
// Method on Session struct
func (s *Session) IsActive() bool {
    return s.Status == "active"
}

// Call the method
if session.IsActive() {
    fmt.Println("Session is active")
}
```

---

#### Interfaces

```go
// Define an interface
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Any type implementing Read() method satisfies the interface
```

---

### WSInspect-Specific Go Concepts

#### 1. GORM (ORM for Database)

GORM is an ORM (Object-Relational Mapping) library that lets you work with databases using Go structs.

```go
// Define a model (struct becomes a database table)
type Session struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    ConnectionID string    `gorm:"uniqueIndex" json:"connection_id"`
    Status       string    `gorm:"default:active" json:"status"`
}

// Create database connection
db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

// Auto-migrate (create tables automatically)
db.AutoMigrate(&Session{})

// Create
session := Session{ConnectionID: "abc"}
db.Create(&session)

// Read
var session Session
db.First(&session, 1)  // Find by ID
db.Where("connection_id = ?", "abc").First(&session)  // Find by condition

// Update
session.Status = "closed"
db.Save(&session)

// Delete
db.Delete(&session)
```

**Tags Explained:**
- `gorm:"primaryKey"` - This field is the primary key
- `gorm:"uniqueIndex"` - Create unique index on this field
- `gorm:"default:active"` - Default value when creating new records
- `gorm:"index"` - Create index for faster queries
- `json:"id"` - JSON serialization field name

---

#### 2. Gin Web Framework

Gin is a web framework similar to Express.js (Node.js) or Flask (Python).

```go
// Create a router
r := gin.Default()

// Basic routing
r.GET("/path", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "Hello"})
})

// Route with parameter
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"user_id": id})
})

// Query parameters
r.GET("/search", func(c *gin.Context) {
    query := c.Query("q")        // Get query param ?q=value
    page := c.DefaultQuery("page", "1")  // With default
    c.JSON(200, gin.H{"query": query, "page": page})
})

// Request body parsing
r.POST("/create", func(c *gin.Context) {
    var requestBody MyStruct
    if err := c.ShouldBindJSON(&requestBody); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, requestBody)
})

// Groups (like Express routers)
api := r.Group("/api")
{
    api.GET("/users", handler)
    api.POST("/users", handler)
}

// Start server
r.Run(":8080")
```

---

#### 3. Gorilla WebSocket

The Gorilla WebSocket library provides WebSocket implementation for Go.

```go
import "github.com/gorilla/websocket"

// Upgrader - HTTP to WebSocket upgrade
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true  // Allow all origins (for development)
    },
}

// Handle WebSocket connection
func handler(w http.ResponseWriter, r *http.Request) {
    // Upgrade HTTP to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()

    // Read message
    messageType, message, err := conn.ReadMessage()
    if err != nil {
        log.Println("Read error:", err)
        return
    }

    // Write message
    err = conn.WriteMessage(messageType, []byte("response"))
    if err != nil {
        log.Println("Write error:", err)
    }
}
```

---

#### 4. Goroutines (Concurrent Execution)

Goroutines are lightweight threads managed by Go's runtime.

```go
// Start a goroutine (like spawning a new thread)
go func() {
    fmt.Println("This runs concurrently")
}()

// Channel - Communication between goroutines
ch := make(chan string)

// Send to channel
go func() {
    ch <- "message"
}()

// Receive from channel
msg := <-ch
```

In WSInspect, goroutines are used for:
- Handling multiple WebSocket connections concurrently
- Running the WebSocket hub's event loop

---

#### 5. Cobra CLI Framework

Cobra is a CLI framework for creating command-line tools.

```go
// Create root command
rootCmd := &cobra.Command{
    Use:   "wsinspect",
    Short: "WebSocket debugging tool",
}

// Create subcommand
startCmd := &cobra.Command{
    Use:   "start",
    Short: "Start the server",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Starting server...")
    },
}

// Add flags
startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port number")

// Add subcommand to root
rootCmd.AddCommand(startCmd)

// Execute
rootCmd.Execute()
```

Usage:
```bash
wsinspect start --port 9000
```

---

#### 6. Viper Configuration

Viper is a configuration management solution.

```go
// Create config
v := viper.New()

// Set config file
v.SetConfigName("config")
v.SetConfigType("yaml")
v.AddConfigPath(".")

// Set defaults
v.SetDefault("server.port", "8080")
v.SetDefault("database.path", "./db")

// Read config
v.ReadInConfig()

// Get values
port := v.GetInt("server.port")
path := v.GetString("database.path")
```

Configuration file (config.yaml):
```yaml
server:
  port: 8080
  host: 0.0.0.0

database:
  type: sqlite
  path: ./wsinspect.db

proxy:
  target: ws://localhost:3000
  buffer_size: 8192

session:
  retention_days: 30

log:
  level: info
```

---

#### 7. The Hub Pattern (WebSocket Management)

The Hub is a design pattern used in WSInspect to manage multiple WebSocket connections.

```go
type Hub struct {
    clients    map[*Client]bool    // All connected clients
    broadcast  chan []byte         // Messages to broadcast
    register   chan *Client        // New client registrations
    unregister chan *Client        // Client disconnections
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

// Run the hub (event loop)
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            
        case client := <-h.unregister:
            delete(h.clients, client)
            
        case message := <-h.broadcast:
            // Broadcast to all clients
            for client := range h.clients {
                client.send <- message
            }
        }
    }
}
```

---

## 10. Testing Process

### Step 1: Verify Backend is Running

```bash
# Check health endpoint
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy","service":"wsinspect"}
```

### Step 2: Verify Proxy Status

```bash
curl http://localhost:8080/api/proxy/status

# Expected response:
# {"status":"running","active_connections":0}
```

### Step 3: Start Frontend

```bash
cd frontend
npm run dev
```

Access http://localhost:5173 in your browser.

### Step 4: Create a Test WebSocket Server

Create a simple WebSocket echo server for testing:

```javascript
// test-server.js (Node.js)
const WebSocket = require('ws');
const server = new WebSocket.Server({ port: 3000 });

server.on('connection', (ws) => {
    console.log('Client connected');
    
    ws.on('message', (message) => {
        console.log('Received:', message.toString());
        // Echo back
        ws.send(`Echo: ${message}`);
    });
});

console.log('WebSocket server on port 3000');
```

Run it:
```bash
node test-server.js
```

### Step 5: Test WebSocket Connection

From the frontend:
1. Navigate to Sessions page
2. Click "New Session"
3. Enter target URL: `ws://localhost:3000`
4. Click Connect

### Step 6: Send Test Messages

Using the frontend:
1. In the message editor, type: `{"type": "hello"}`
2. Click Send
3. You should see the echo response

### Step 7: Verify Messages are Stored

```bash
curl http://localhost:8080/api/messages/session/1
```

### Step 8: Test Replay

1. Go to Session Detail page
2. Select messages to replay
3. Click Replay button
4. Verify messages are resent

### Step 9: Test Fuzz Testing

1. Go to Fuzz Testing page
2. Create a new test with template: `{"type": "{{data}}"}`
3. Set strategy to "random"
4. Run the test
5. View results

---

## 11. Troubleshooting

### Issue: Frontend shows "Cannot connect to backend"

**Solution:** Ensure backend is running on port 8080:
```bash
curl http://localhost:8080/health
```

### Issue: WebSocket connection fails

**Solution:** 
1. Check target server is running
2. Verify CORS settings in backend
3. Check firewall settings

### Issue: Database errors

**Solution:**
1. Check database file permissions
2. Ensure SQLite is installed: `go get gorm.io/driver/sqlite`

### Issue: Port already in use

**Solution:**
```bash
# Find process using port 8080
lsof -i :8080  # Linux/Mac
netstat -ano | findstr :8080  # Windows

# Kill the process or use different port
go run main.go start --port 9000
```

---

## Summary

WSInspect is a comprehensive WebSocket debugging tool built with:

- **Backend:** Go + Gin + GORM + SQLite
- **Frontend:** React + Vite + Tailwind + TypeScript
- **WebSocket:** Gorilla WebSocket library
- **Architecture:** REST API + WebSocket Hub pattern

This guide covered:
- ✅ What WSInspect is and what it does
- ✅ Installation and setup
- ✅ How frontend and backend communicate
- ✅ Every API endpoint with examples
- ✅ Go concepts explained for beginners
- ✅ Step-by-step testing process

For more information, check the [SPEC.md](./SPEC.md) and [Architecture.md](./Architecture.md).
