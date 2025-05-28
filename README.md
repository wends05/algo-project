# Algorithm Final Project

- Small exploration game where going through different nodes gives or costs energy.
- The goal of the game is to reach the end node with the shortest amount of energy consumed
- Uses a randomized graph, with nodes as vertices, and transfers as edges with weights (that can be negative)
- Utilizes the **Bellman Ford Algorithm** as the shortest path algorithm to find the shortest path


## Technical Description

- Gin (Go) as the backend for:
  - generating a graph
  - executing the algorithm
- Godot Engine for the game client

## Installation Requirements

Server:
- Requires Go (1.24) - Language for the backend
- (Optional) Air - Live reload for the backend

Game:
- Godot Engine v4.4.1 - Game Engine

## Installation

Clone the repository
```
git clone https://github.com/wends05/algo-project.git
```

### Running the server
1. Go to the server directory
```
cd bellman_ford_server
```

2. Install Go modules and dependencies
```
go mod tidy
```

3. Run the go server with
```
air 
```
assuming air has been installed for live reload,
or with
```
go run .
```
without any live reload

### Running the game
1. Open godot engine and import the `algo_game` folder to godot
2. Run the game with F5 (Windows) or Cmd+B (Mac)
