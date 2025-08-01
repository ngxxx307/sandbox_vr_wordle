# VR Wordle Sandbox

This project is a backend server for a Wordle-like game, built with Go. It features both a single-player mode and a real-time, turn-based two-player mode using WebSockets.

The architecture is designed to be robust and scalable, using modern Go practices.

## Features

*   **State Machine Architecture**: The server uses a state machine pattern (`controller.Controller` interface) to manage the user's session, transitioning seamlessly between states like a "Game Lounge" and a "Multiplayer Game".
*   **Real-time Multiplayer**:
    *   A central `Hub` manages matchmaking, placing players into a queue.
    *   When two players are matched, a dedicated `GameSession` is created to manage their game state independently from other sessions.
    *   Communication is fully asynchronous using a channel-based "read pump/write pump" pattern for each WebSocket connection.
*   **Dependency Injection**: Leverages Uber's `fx` framework for clean dependency injection and graceful application lifecycle management (`OnStart`, `OnStop` hooks).
*   **Configuration Management**: Uses Viper to load configuration from environment variables, with a sample `.env` file for easy setup.

## Tech Stack

*   **Language**: Go
*   **Web Framework**: [Echo](https://echo.labstack.com/)
*   **WebSocket Library**: [Gorilla WebSocket](http://www.gorillatoolkit.org/pkg/websocket)
*   **Dependency Injection**: [Uber Fx](https://github.com/uber-go/fx)
*   **Configuration**: [Viper](https://github.com/spf13/viper)
*   **Environment Variables**: [GoDotEnv](https://github.com/joho/godotenv)

## Project Structure

```
.
├── config/         # Configuration loading (Viper) and struct definition.
├── controller/     # Handles WebSocket state transitions (e.g., GameLounge -> Multiplayer).
├── hub/            # Manages matchmaking (Hub) and individual game logic (GameSession).
├── http/           # Main application entrypoint, fx setup, and Echo server lifecycle.
├── routes/         # Defines API routes and registers handlers.
├── service/        # Core game services and data structures (e.g., MultiplayerWordle service).
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

1.  You will need two separate WebSocket clients to test the multiplayer functionality. You can use a tool like the [Simple WebSocket Client](https://chrome.google.com/webstore/detail/simple-websocket-client/pfdhoblngboilpfeibdedpjgfnlcodoo) extension for Chrome, or any other client.

2.  **Connect both clients** to the server at:
    `ws://localhost:8080/ws`

3.  In the first client, send the message `multiplayer`. You will receive a "Waiting for another player..." message.

4.  In the second client, send the message `multiplayer`. The game will start, and both clients will receive the game rules.

5.  Follow the turn-based instructions to play the game.

## ⚠️ Limitations

> **Task 4 (Multiplayer)** is only partially finished due to time limit pressure.