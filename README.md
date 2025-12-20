# TeamTime

[![Test](https://github.com/matteo-gildone/teamtime/actions/workflows/test.yml/badge.svg)](https://github.com/matteo-gildone/teamtime/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/matteo-gildone/teamtime)](https://goreportcard.com/report/github.com/matteo-gildone/teamtime)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A CLI tool to view local times for your distributed team members.

## Installation
```bash
go install github.com/matteo-gildone/teamtime@latest
```

Or build from source:
```bash
git clone https://github.com/matteo-gildone/teamtime.git
cd teamtime
go build -o teamtime
```

## Demo

![TeamTime Demo](img/demo.gif)

### Colors disabled

![TeamTime Demo](img/demo-nocolor.gif)

## Quick Start
```bash
# Initialize the app
teamtime init

# Add team members
teamtime add "Alice" "London" "Europe/London"
teamtime add "Priya" "Pune" "Asia/Kolkata"
teamtime add "Lucio" "Poggibonsi" "Europe/Rome"

# View everyone's local time
teamtime check all

# View specific team member
teamtime check Alice
```

## Commands

### `init`
Initialize TeamTime configuration directory (`~/.teamtime`)
```bash
teamtime init
```

### `add`
Add a new team member
```bash
teamtime add <name> <city> <timezone>

# Example
teamtime add "Bob" "Berlin" "Europe/Berlin"
```

Find valid timezone names at [Wikipedia - List of tz database time zones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)

### `check`
Display current local time for all team members
```bash
teamtime check all
```

Output:
```
ID                   | Name                 |  Local Time          
-------------------- | -------------------- | --------------------
1                    | Alice                | 09:30 (Mon 20 Nov)    
2                    | Priya                | 15:00 (Mon 20 Nov)    
3                    | Marco                | 10:30 (Mon 20 Nov)   
```

Display current local time for a particular team members
```bash
teamtime check <name>

# Example
teamtime check Alice
```

Output:
```
ID                   | Name                 | Local Time          
-------------------- | -------------------- | --------------------
1                    | Alice                | 09:30 (Mon 20 Nov)
```

### `remove`
Remove a team member by ID
```bash
teamtime remove <id>

# Example
teamtime remove 2
```

## Configuration

TeamTime stores data in `~/.teamtime/colleagues.json`

## License

MIT