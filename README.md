# Pokedex CLI

A command-line interface (CLI) for interacting with the PokeAPI, allowing users to explore the world of Pokémon, discover new areas, catch Pokémon, and manage their personal Pokedex.

## Features

* **Explore Areas**: Navigate through different locations in the Pokémon world to find wild Pokémon.
* **Catch Pokémon**: Attempt to catch wild Pokémon you encounter and add them to your Pokedex.
* **Inspect Caught Pokémon**: View detailed information about Pokémon you've successfully caught.
* **View Pokedex**: See a list of all the Pokémon you've collected.
* **Command History (New!)**: Easily cycle through previously entered commands using the **Up** and **Down** arrow keys for a more efficient and user-friendly experience.
* **Interactive Input**: Type commands, use backspace, and experience real-time feedback.
* **Caching**: Efficiently caches API responses to reduce network requests and improve performance.

## How to Run

To run this Pokedex CLI, ensure you have Go installed on your system.

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Ikit24/pokedexcli
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

### Command History Navigation:

* Press the **Up Arrow** key to cycle through previously entered commands (from most recent to oldest).
* Press the **Down Arrow** key to cycle forward through commands, or return to an empty input line.
