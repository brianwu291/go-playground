# Go Playground

A collection of Go examples and implementations demonstrating various Go programming concepts, patterns, and real-world applications.

## Project Structure

- `adder/`: Addition operation implementation
- `groundone/`: Basic Go programming examples
- `handlers/`: 
  - `websocket/`: WebSocket connection handling
- `interview/`: Common Go interview problem solutions
- `randomstrgenerater/`: Random string generation utilities
- `ratelimiter/`: Rate limiting implementation
- `realtimechat/`: Real-time chat application implementation
- `safecounter/`: Thread-safe counter implementation
- `templates/`: HTML templates
  - `chat.html`: Chat interface template
- `unbufferhandle/`: Unbuffered channel handling examples
- `utils/`: Utility functions and helpers
- `workerpool/`: Worker pool pattern implementation

## Prerequisites

- Go 1.23.4 (version specified in `go.mod`)
- Docker (for containerization)

## Getting Started

1. Clone the repository
2. Run: `$ go run main.go` or `$ air` for hot reload

## Deployment

The project includes Docker support and can be deployed using Fly.io:
- `Dockerfile`: Container configuration
- `fly.toml`: Fly.io deployment configuration

## Features

- WebSocket Implementation
- Real-time Chat Functionality
- Rate Limiting
- Thread-safe Operations
- Worker Pool Pattern
- Random String Generation
- Interview Problem Solutions
- Basic Go Concepts

## License

See the [LICENSE](LICENSE) file for details.
