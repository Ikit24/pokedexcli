# Pokedex CLI

A command-line interface (CLI) for interacting with the PokeAPI, allowing users to explore the world of Pokémon, discover new areas, catch Pokémon, and manage their personal Pokedex.

## Features

## 🎮 **Complete Battle System**
* **Turn-based combat** with speed-based turn order
* **Move selection** with input validation and retry logic
* **Realistic damage calculation** using Pokemon stats
* **AI opponent** that uses strongest available moves
* **Multi-turn battles** until one Pokemon faints

### 📈 **XP and Leveling System**
* **Authentic Pokemon XP formula** (level³ progression)
* **XP gain** from defeating opponents based on their base experience
* **Level cap at 100** (like real Pokemon games)
* **Level-up notifications** with colorized output
* **XP persistence** between game sessions

### 💾 **Save/Load System**
* **Automatic save** after battles and captures
* **Manual save command** with retry logic and user confirmation
* **Automatic loading** on game startup
* **Corrupted file recovery** (deletes and creates fresh save)
* **Progress persistence** - Pokemon keep XP, levels, and collection between sessions

### 🎯 **Enhanced User Experience**
* **Input validation loops** - invalid input asks again instead of crashing
* **Battle confirmation** - preview stats before fighting
* **Capture mechanics** - option to catch defeated Pokemon
* **Error handling** - graceful recovery from invalid commands
* **Command dependencies** - must run `map` before `explore`

### 🏗️ **Code Organization**
* **Modular architecture** - separated commands into individual files
* **Clean functions** - extracted `playerTurn()`, `opponentTurn()`, `checkVictory()`
* **Shared utilities** - XP calculation functions available across commands
* **Proper package structure** - commands in dedicated package

### 📦 **Core Features**
* **Explore Areas**: Navigate through different locations in the Pokémon world to find wild Pokémon.
* **Catch Pokémon**: Attempt to catch wild Pokémon you encounter and add them to your Pokedex.
* **Inspect Caught Pokémon**: View detailed information about Pokémon you've successfully caught.
* **View Pokedex**: See a list of all the Pokémon you've collected.
* **Command History**: Easily cycle through previously entered commands using the **Up** and **Down** arrow keys for a more efficient and user-friendly experience.
* **Interactive Input**: Type commands, use backspace, and experience real-time feedback.
* **Caching**: Efficiently caches API responses to reduce network requests and improve performance.

## How to Run

To run this Pokedex CLI, ensure you have Go installed on your system.

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-github-username/your-repo-name.git
   cd your-repo-name
   ```

2. **Install dependencies**:
   This project uses `github.com/eiannone/keyboard` for interactive input.
   ```bash
   go mod tidy
   ```
   (Or, if `go mod tidy` doesn't work, `go get github.com/eiannone/keyboard`)

3. **Run the application**:
   ```bash
   go run .
   ```

## Usage

Once the application is running, you will see a `Pokedex >` prompt. Type commands and press `Enter`.

### Commands:

* `help`: Displays a list of all available commands and their descriptions.
* `exit`: Closes the Pokedex application.
* `map`: Displays the next 20 location areas.
* `mapb`: Displays the previous 20 location areas.
* `explore <location_area_name>`: Lists the Pokémon found in a specific location area (e.g., `explore oreburgh-mine-1f`).
* `catch <pokemon_name>`: Attempts to catch a specified Pokémon (e.g., `catch pikachu`).
* `inspect <pokemon_name>`: Shows detailed information about a Pokémon you've caught (e.g., `inspect charmander`).
* `pokedex`: Displays a list of all Pokémon currently in your Pokedex.
* `battle <pokemon_name>`: Initiate turn-based combat with a wild Pokémon (e.g., `battle pikachu`).
* `save`: Manually save your current progress to disk.

### Command History Navigation:

* Press the **Up Arrow** key to cycle through previously entered commands (from most recent to oldest).
* Press the **Down Arrow** key to cycle forward through commands, or return to an empty input line.

## Ideas for Extending the Project

This project provides a solid foundation, but there are many ways it could be expanded and improved:

* ~~Simulate battles between Pokémon.~~ ✅ **COMPLETED** - Full turn-based battle system implemented
* ~~Add more comprehensive unit tests.~~ 
* ~~Implement a "party" system for Pokémon, allowing them to level up.~~ ✅ **COMPLETED** - XP and leveling system implemented
* Allow caught Pokémon to evolve after certain conditions are met.
* ~~Persist a user's Pokedex to disk (e.g., using JSON files) so progress is saved between sessions.~~ ✅ **COMPLETED** - Full save/load system implemented
* Enhance exploration with random encounters or choice-based navigation.
* Add support for different types of Pokeballs with varying catch rates.
* ~~Refactor for better organization and testability.~~ ✅ **COMPLETED** - Modular architecture implemented
* Add Pokemon type effectiveness system for battles.
* Implement Pokemon evolution system.
* Add more battle moves and abilities.
* Create a tournament or gym leader system.
* Add Pokemon breeding mechanics.
* Implement status effects (poison, burn, paralysis, etc.).
* Add multiplayer battle support.
