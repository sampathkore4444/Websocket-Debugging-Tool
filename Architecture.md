# Architecture.md — WSInspect (WebSocket Debugger)

## 1. Overview

WSInspect is a developer tool for **recording, inspecting, modifying, and replaying WebSocket traffic** between clients and servers. It acts as a **WebSocket proxy** between client and server, providing Postman-like debugging capabilities for WebSocket APIs.

---

## 2. Recommended Tech Stack

### 2.1 Backend (Go)

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **Language** | [Go 1.21+](https://go.dev/) | High performance, native concurrency with goroutines |
| **WebSocket Library** | [gorilla/websocket](https://github.com/gorilla/websocket) | Mature, battle-tested WebSocket implementation |
| **Web Framework** | [Gin](https://gin-gonic.com/) | Fast, minimalist HTTP web framework |
| **Database** | [SQLite](https://www.sqlite.org/) (local) / [PostgreSQL](https://www.postgresql.org/) (cloud) | Session storage, configurable per deployment |
| **ORM** | [GORM](https://gorm.io/) | Full-featured ORM for Go |
| **Caching** | [Redis](https://redis.io/) | Real-time communication and caching |
| **Configuration** | [Viper](https://github.com/spf13/viper) | Complete configuration solution |
| **CLI** | [Cobra](https://github.com/spf13/cobra) | CLI argument parsing and subcommands |

### 2.2 Frontend (React + Vite + Tailwind)

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **Framework** | [React 18+](https://react.dev/) | Component-based UI for complex dashboard |
| **Build Tool** | [Vite](https://vitejs.dev/) | Fast development and build times |
| **Language** | [TypeScript](https://www.typescriptlang.org/) | Type safety and better developer experience |
| **State Management** | [Zustand](https://zustand-demo.pmnd.rs/) | Lightweight, reactive state management |
| **UI Components** | [Shadcn UI](https://ui.shadcn.com/) + [Tailwind CSS](https://tailwindcss.com/) | Modern, customizable, accessible components |
| **WebSocket Client** | Native [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket) | Connect to proxy dashboard |
| **Code Editor** | [Monaco Editor](https://microsoft.github.io/monaco-editor/) | Message editing with syntax highlighting |
| **Charts** | [Recharts](https://recharts.org/) | Session analytics visualization |
| **HTTP Client** | [Axios](https://axios-http.com/) | REST API communication |

### 2.3 CLI (Go)

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **Framework** | [Cobra](https://github.com/spf13/cobra) | Robust CLI argument parsing |
| **Terminal UI** | [Bubble Tea](https://github.com/charmbracelet/bubbletea) | Interactive CLI experience |
| **I/O** | Native Go libraries | Clean integration with system |

### 2.4 DevOps & Infrastructure

| Component | Technology | Rationale |
|-----------|------------|-----------|
| **Containerization** | [Docker](https://www.docker.com/) | Easy deployment and distribution |
| **Build Tool** | [Make](https://www.gnu.org/software/make/) | Simplified build commands |
| **Testing** | [Go testing](https://go.dev/testing/) (Unit) + [Playwright](https://playwright.dev/) (E2E) | Comprehensive test coverage |

---

## 3. Core Features

### 3.1 WebSocket Proxy

- Intercept WebSocket connections (ws:// and wss://)
- Inspect, modify, and inject frames in real-time
- Support both client-to-server and server-to-client traffic
- TLS/SSL interception with certificate management

### 3.2 Session Recording

- Capture full WebSocket sessions with metadata:
  - `connection_id`, `timestamp`, `direction`, `message_payload`, `opcode`, `latency`
- Export formats: JSON, NDJSON, Binary
- Configurable retention policies

### 3.3 Message Inspector

- Real-time message display in UI panel
- Display modes: Raw Text, JSON, Binary, Hex
- Syntax highlighting and JSON formatting
- Diff view for comparing messages
- Filters: message type, payload text, direction, time range

### 3.4 Message Editing

- Intercept and modify messages before forwarding
- Support for JSON, Protobuf, MsgPack formats
- Real-time editing with validation

### 3.5 Message Injection

- Manually send custom messages
- Simulate events and test server behavior

### 3.6 Session Replay

- Replay previously recorded sessions
- **Exact Replay**: preserves order, payload, and timing
- **Fast Replay**: removes delays for quick testing
- **Edited Replay**: allows message modification during replay

### 3.7 Session Scripting

- Scripted message flows using YAML, JSON, or JavaScript
- Example script:
  ```yaml
  steps:
    - send: {"type": "login"}
    - wait: 1000
    - send: {"type": "subscribe", "channel": "orders"}
  ```

### 3.8 Message Fuzzing

- Generate malformed messages for testing
- Fuzz cases: missing fields, invalid JSON, large payloads, random strings
- Purpose: test server validation and find bugs

### 3.9 Protocol Awareness

- Detect common message structures
- Support formats: JSON, Protobuf, MsgPack, Binary
- Automatic decoding and schema display

### 3.10 Security Features

- TLS interception with certificate management
- Payload masking for sensitive data
- Configurable redaction (e.g., passwords → ****)

---

## 4. System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐ │
│  │   Browser    │  │   Mobile     │  │    Desktop  │  │   CLI (Terminal)│ │
│  │  Application │  │     App      │  │     App     │  │                  │ │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └────────┬─────────┘ │
└─────────┼─────────────────┼─────────────────┼──────────────────┼──────────┘
          │                 │                 │                  │
          │    WebSocket    │    WebSocket    │   WebSocket      │    STDIN
          │    (ws/wss)      │    (ws/wss)     │    (ws/wss)      │    /API
          │                 │                 │                  │
          ▼                 ▼                 ▼                  ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           WSINSPECT PROXY SERVER                            │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                         API Gateway / REST Server                     │ │
│  │   (Fastify/Express + WebSocket Upgrade Handler)                       │ │
│  └────────────────────────────────┬───────────────────────────────────────┘ │
│                                   │                                          │
│  ┌────────────────────────────────┼───────────────────────────────────────┐ │
│  │                        Core Services Layer                             │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │ │
│  │  │   Proxy     │  │  Session    │  │   Replay    │  │    Fuzz     │   │ │
│  │  │   Engine    │  │  Recorder   │  │   Engine    │  │   Engine    │   │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │ │
│  │  │   Message   │  │   Script    │  │  Protocol   │  │   Message   │   │ │
│  │  │   Editor    │  │   Runner    │  │   Parser    │  │  Injector   │   │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
│                                   │                                          │
│  ┌────────────────────────────────┼───────────────────────────────────────┐ │
│  │                      Data Layer (Storage)                             │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │ │
│  │  │   SQLite    │  │    Redis    │  │    File     │  │   Message   │   │ │
│  │  │  (Sessions) │  │   (Cache)   │  │   System    │  │    Queue    │   │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
                                   │
                                   │ Proxied WebSocket
                                   ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        TARGET WEBSOCKET SERVER                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐ │
│  │    Chat      │  │   Trading    │  │   Gaming     │  │   Real-time      │ │
│  │  Server      │  │   Platform    │  │   Server     │  │   Dashboard      │ │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Component Architecture

### 5.1 Proxy Engine

```
┌─────────────────────────────────────────────────────────────────┐
│                      Proxy Engine Flow                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   Client ──▶ [TCP Listener] ──▶ [Connection Handler]           │
│                          │                                      │
│                          ▼                                      │
│              ┌─────────────────────────┐                       │
│              │   Frame Interceptor     │                       │
│              │  ┌─────────────────────┐ │                       │
│              │  │ • Inspect Frame     │ │                       │
│              │  │ • Modify Frame      │ │                       │
│              │  │ • Drop Frame        │ │                       │
│              │  │ • Inject Frame      │ │                       │
│              │  └─────────────────────┘ │                       │
│              └───────────┬─────────────┘                       │
│                          │                                      │
│                          ▼                                      │
│              ┌─────────────────────────┐                       │
│              │   Upstream Connector   │                       │
│              └───────────┬─────────────┘                       │
│                          │                                      │
│                          ▼                                      │
│                        Server                                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 5.2 Session Recorder

```
┌─────────────────────────────────────────────────────────────────┐
│                    Session Recorder Architecture                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐     │
│  │   Incoming   │    │   Message    │    │    Output    │     │
│  │   Frames     │───▶│   Processor  │───▶│   Writer     │     │
│  └──────────────┘    └──────────────┘    └──────────────┘     │
│         │                   │                   │             │
│         ▼                   ▼                   ▼             │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              Session Metadata Collector                  │  │
│  │  • connection_id    • timestamp    • direction          │  │
│  │  • opcode           • latency      • message_size        │  │
│  └─────────────────────────────────────────────────────────┘  │
│                               │                                 │
│                               ▼                                 │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │              Storage (SQLite / File System)              │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 5.3 Replay Engine

```
┌─────────────────────────────────────────────────────────────────┐
│                      Replay Engine                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐     │
│  │   Session    │    │   Playback   │    │   WebSocket  │     │
│  │   Loader     │───▶│   Controller │───▶│   Emulator   │     │
│  └──────────────┘    └──────────────┘    └──────────────┘     │
│         │                   │                   │             │
│         ▼                   ▼                   ▼             │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │                   Replay Modes                           │  │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────────────┐  │  │
│  │  │   Exact    │  │   Fast     │  │      Edited        │  │  │
│  │  │  Replay    │  │  Replay    │  │     Replay         │  │  │
│  │  │ (preserve  │  │ (skip      │  │  (allow message   │  │  │
│  │  │  timing)   │  │  delays)   │  │   modification)   │  │  │
│  │  └────────────┘  └────────────┘  └────────────────────┘  │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 5.4 Fuzz Engine

```
┌─────────────────────────────────────────────────────────────────┐
│                       Fuzz Engine                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐     │
│  │   Message    │    │    Fuzz      │    │   Test       │     │
│  │   Template   │───▶│   Generator  │───▶│   Runner     │     │
│  └──────────────┘    └──────────────┘    └──────────────┘     │
│                             │                    │             │
│                             ▼                    ▼             │
│                    ┌────────────────┐    ┌──────────────┐     │
│                    │  Fuzz Strategies│    │  Results     │     │
│                    │  • Random      │    │  Collector   │     │
│                    │  • Mutation    │    │              │     │
│                    │  • Boundary    │    │              │     │
│                    │  • Invalid     │    │              │     │
│                    └────────────────┘    └──────────────┘     │
│                                                                  │
│  Example Fuzz Cases:                                            │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │ {"type": null}           │ Invalid JSON                 │  │
│  │ {"type": 999999999999}   │ Out of bounds                │  │
│  │ {"type": ""}             │ Empty string                 │  │
│  │ {}                       │ Missing required fields     │  │
│  │ [1,2,3]                  │ Wrong type (object expected) │  │
│  └─────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 6. Data Models

### 6.1 Session

```typescript
interface Session {
  id: string;
  connection_id: string;
  client_ip: string;
  server_host: string;
  start_time: Date;
  end_time?: Date;
  status: 'active' | 'closed' | 'error';
  message_count: number;
}
```

### 6.2 Message

```typescript
interface WebSocketMessage {
  id: string;
  session_id: string;
  timestamp: Date;
  direction: 'client-to-server' | 'server-to-client';
  opcode: number;
  payload: string | Buffer;
  payload_format: 'text' | 'json' | 'binary' | 'hex';
  latency_ms?: number;
  is_modified: boolean;
  original_payload?: string;
}
```

---

## 7. API Endpoints

### 7.1 Proxy Control

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/proxy/start` | Start proxy server |
| POST | `/api/proxy/stop` | Stop proxy server |
| GET | `/api/proxy/status` | Get proxy status |

### 7.2 Session Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/sessions` | List all sessions |
| GET | `/api/sessions/:id` | Get session details |
| DELETE | `/api/sessions/:id` | Delete session |
| POST | `/api/sessions/:id/export` | Export session |

### 7.3 Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/sessions/:id/messages` | Get session messages |
| POST | `/api/messages/inject` | Inject new message |
| PUT | `/api/messages/:id` | Modify message |

### 7.4 Replay & Fuzz

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/replay` | Start replay |
| POST | `/api/fuzz/start` | Start fuzzing |
| GET | `/api/fuzz/results` | Get fuzz results |

---

## 8. CLI Commands

```bash
# Start proxy server
wsinspect start --port 8080 --target ws://localhost:3000

# List sessions
wsinspect sessions list

# Replay session
wsinspect replay session_id [--mode exact|fast|edited]

# Export session
wsinspect export session_id --format json --output ./session.json

# Start fuzzing
wsinspect fuzz start --template ./template.json

# Run scripted session
wsinspect script run ./script.yaml
```

---

## 9. Performance Targets

| Metric | Target |
|--------|--------|
| Proxy Latency | < 5ms |
| Message Throughput | 10,000 messages/sec |
| Concurrent Connections | 1,000 connections |
| Session Storage | Up to 100GB |

---

## 10. File Structure

```
wsinspect/
├── backend/
│   ├── src/
│   │   ├── core/
│   │   │   ├── proxy/
│   │   │   │   ├── proxy.server.ts
│   │   │   │   ├── connection.handler.ts
│   │   │   │   └── frame.interceptor.ts
│   │   │   ├── session/
│   │   │   │   ├── session.recorder.ts
│   │   │   │   └── session.store.ts
│   │   │   ├── replay/
│   │   │   │   └── replay.engine.ts
│   │   │   └── fuzz/
│   │   │       └── fuzz.engine.ts
│   │   ├── api/
│   │   │   ├── routes/
│   │   │   └── middleware/
│   │   ├── services/
│   │   └── utils/
│   └── package.json
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── stores/
│   │   └── hooks/
│   ├── package.json
│   └── vite.config.ts
├── cli/
│   ├── src/
│   │   ├── commands/
│   │   └── index.ts
│   └── package.json
└── docker/
    └── Dockerfile
```

---

## 11. Extensibility

### Plugin System

```typescript
// Plugin Interface Example
interface WsInspectPlugin {
  name: string;
  version: string;
  
  // Called when a message is intercepted
  onMessage?(message: WebSocketMessage): WebSocketMessage | null;
  
  // Called when a connection is established
  onConnect?(connection: Connection): void;
  
  // Called when a connection is closed
  onDisconnect?(connection: Connection): void;
}

// Example Plugin
const myPlugin: WsInspectPlugin = {
  name: 'json-validator',
  version: '1.0.0',
  onMessage(message) {
    if (message.payload_format === 'json') {
      try {
        JSON.parse(message.payload);
      } catch (e) {
        console.error('Invalid JSON detected!');
      }
    }
    return message;
  }
};
```

---

## 12. Security Considerations

- **TLS Interception**: Support for decrypting wss:// traffic with user-provided certificates
- **Data Masking**: Automatic redaction of sensitive fields (passwords, tokens, API keys)
- **Local Storage**: All session data stored locally by default
- **Optional Encryption**: Additional encryption layer for stored sessions

---

## 13. Future Enhancements

- [ ] Visual Studio Code extension
- [ ] Postman collection import/export
- [ ] Cloud storage integration
- [ ] Team collaboration features
- [ ] GraphQL WebSocket support
- [ ] gRPC-WebSocket tunneling
