# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Environment

This project uses devbox and direnv for isolated development environments. Always work within the devbox environment:

```bash
devbox shell
```

Any tools not part of Go dependencies should be installed via devbox.

## Common Commands

### Running the Application
```bash
go run main.go <amount> <from> <to>
# Example: go run main.go 100 USD EUR
```

### Testing
```bash
go test
go test -v  # verbose output
```

### Building
```bash
go build -o conv
```

## Architecture

This is a simple CLI currency converter with clean separation of concerns:

- **main.go**: Entry point, argument parsing, and CLI logic
- **converter.go**: Currency conversion logic with interface-based design
- **main_test.go**: Comprehensive unit tests with mocks

### Key Components

- `Converter` interface: Abstraction for currency conversion implementations
- `ApiCurrencyConverter`: HTTP-based converter using Fawaz Ahmed's Currency API
- `FawazConversion`: Custom JSON unmarshaling for the API response format
- `Currency` type: Type-safe currency representation with validation
- `Input` struct: Parsed command-line arguments

### Design Principles

- Interface-based design for testability
- Custom JSON unmarshaling for external API integration
- Comprehensive error handling and validation
- Type safety with custom Currency type
- Clean separation between parsing, conversion, and presentation logic

## Testing Philosophy

- Unit tests are mandatory for every new feature
- Use dependency injection and mocks for external dependencies
- Run tests immediately after changes: `go test`
- Tests cover both happy path and error scenarios