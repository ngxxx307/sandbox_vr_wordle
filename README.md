# VR Wordle Sandbox

This project is a backend server for a Wordle-like game, built with Go. It features both a single-player mode and a real-time, turn-based two-player mode using WebSockets.

The architecture is designed to be robust and scalable, using modern Go practices.

## Tech Stack

*   **Language**: Go
*   **Web Framework**: [Echo](https://echo.labstack.com/)
*   **WebSocket Library**: [Gorilla WebSocket](http://www.gorillatoolkit.org/pkg/websocket)
*   **Dependency Injection**: [Uber Fx](https://github.com/uber-go/fx)
*   **Configuration**: [Viper](https://github.com/spf13/viper)
*   **Environment Variables**: [GoDotEnv](https://github.com/joho/godotenv)

## Game Mode

*   **Single-Player Wordle**: Classic Wordle gameplay where players guess a 5-letter word within a limited number of attempts.
*   **Cheated Host Mode**: Single-player mode with the answer revealed to the host, useful for testing and debugging.
*   **Multiplayer Wordle**: Real-time, turn-based two-player mode where players compete to guess the same word, taking alternate turns.

## Core Concept

**Controller**: Handles WebSocket communication and user session flow. Controllers manage state transitions (e.g., GameLounge → Wordle), process user input from WebSocket connections, and coordinate between different game modes.

**Service**: Contains pure business logic for game mechanics. Services handle core Wordle gameplay (word validation, scoring, game rules) without any knowledge of WebSocket connections or user sessions.

**Match**: Manages individual multiplayer game sessions, coordinating communication between players and the game host.

**MatchMaker**: Contains the multiplayer game matching logic, queueing players and creating matches when sufficient players are available.

**WebSocket**: Contains WebSocket utility functions for connection management, message handling, and communication protocols.

## Project Structure

```
.
├── config/         # Configuration loading (Viper) and struct definition.
├── controller/     # Handles WebSocket state transitions (e.g., GameLounge -> Multiplayer).
├── env/            # Directory for environment files (e.g., .env).
├── http/           # Main application entrypoint, fx setup, and Echo server lifecycle.
├── match/            # Manages individual game sessions/matches.
├── matchMaker/     # Handles matchmaking logic, queueing players for multiplayer games.
├── routes/         # Defines API routes and registers handlers.
├── service/        # Core game services and data structures (e.g., Wordle game logic).
└── websocket/      # Low-level WebSocket connection wrapper and message pump logic.
```

## Getting Started

### Prerequisites

*   Go (version 1.24)

### Installation & Running

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/ngxxx307/sandbox_vr_wordle.git
    cd sandbox_vr_wordle
    ```

2.  **Set up environment variables:**
    Create a `.env` file in the `env/` directory by copying the sample file.
    ```sh
    cp env/.env.sample env/.env
    ```
    You can modify the variables in `env/.env` if needed, but the defaults should work.

3.  **Run the server:**
    ```sh
    go run ./http/http.go
    ```
    The server will start and print a message like `Server listening on port 8080`.

### How to Play

1.  **Connect to the server** using a WebSocket client (e.g. postman)  at:
    `ws://localhost:8080/ws`

2.  **Choose a game mode.** Upon connecting, you'll be in the game lounge. Send one of the following messages to start a game:
    *   `Wordle`: Start a standard single-player game.
    *   `Cheated Host Wordle`: Start a single-player game where the answer is revealed to the host (useful for testing).
    *   `Multiplayer Wordle`: Enter the queue for a two-player match.

3.  **For Multiplayer Games:**
    *   Connect a second client to the same address (`ws://localhost:8080/ws`).
    *   From the second client, also send the `multiplayer` message.
    *   Once both players have joined, the match will begin automatically, and both clients will receive the game rules.

4.  **Follow the on-screen instructions** to play.