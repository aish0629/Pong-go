# Pong-go

A lightweight Pong clone written in Go. This repository contains a simple implementation of the classic Pong game intended for learning Go game loops, basic graphics, and input handling.

## Features

- Traditional two-paddle Pong gameplay
- Keyboard controls for two players
- Simple, minimal codebase suitable for learning and extension

## Requirements

- Go 1.18+ installed
- (Optional) SDL / Ebiten / raylib or other graphics library if the project uses one — check the code to see which library is used and install it if required.

## Setup

1. Clone the repository:

   git clone https://github.com/your-username/Pong-go1.git
   cd Pong-go1

2. Fetch dependencies (if any). For example, if using Ebiten:

   go get -u github.com/hajimehoshi/ebiten/v2

## Build & Run

To build:

   go build ./...

To run:

   go run .

Or, if the entrypoint is a specific file:

   go run cmd/main.go

## Controls

- Player 1: W (up), S (down)
- Player 2: Up Arrow, Down Arrow
- Pause / Quit: check the main loop for exact keys used

## Contributing

Contributions are welcome. Open an issue or submit a pull request with a clear description of changes.

## License

Specify a license for the project (e.g., MIT). Add a LICENSE file to the repository.

## Notes

Check the source files to confirm which graphics/input library is used and update the README commands accordingly (dependency installation, run instructions).