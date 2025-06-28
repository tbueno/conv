# Currency Converter

A fast and simple CLI currency converter written in Go that supports a wide range of currencies using real-time exchange rates from [Fawaz Ahmed's Currency API](https://github.com/fawazahmed0/exchange-api).

## Features

- **340+ Currencies**: Support for fiat currencies, cryptocurrencies, and precious metals
- **Real-time Rates**: Uses live exchange rates from a reliable API
- **Offline Cache**: Automatically caches currency list for faster performance
- **Clean CLI**: Built with Cobra framework for excellent user experience
- **Fast & Lightweight**: Single binary with no external dependencies

## Installation

### Build from Source

```bash
git clone <repository-url>
cd conv
go build
```

## Usage

### Convert Currency

```bash
conv <amount> <from> <to>
```

**Examples:**
```bash
conv 100 USD EUR     # Convert 100 USD to EUR
conv 50 GBP JPY      # Convert 50 GBP to JPY
conv 1000 BTC USD    # Convert 1000 BTC to USD
conv 25.5 EUR BRL    # Convert 25.5 EUR to BRL
```

### List Available Currencies

```bash
conv --list          # Show all supported currencies
conv -l              # Short form
```

### Get Help

```bash
conv --help          # Show usage information
conv -h              # Short form
```

## Supported Currencies

The converter supports 340+ currencies including:
- **Fiat currencies**: USD, EUR, GBP, JPY, CAD, AUD, etc.
- **Cryptocurrencies**: BTC, ETH, ADA, DOT, SOL, etc.
- **Precious metals**: Gold (XAU), Silver (XAG), Platinum (XPT), Palladium (XPD)

Use `conv --list` to see the complete list of supported currencies.

## Development

### Project Structure

```
conv/
├── cmd/                    # Cobra CLI commands
│   └── root.go
├── internal/
│   ├── currency/          # Currency types and validation
│   │   └── types.go
│   └── converter/         # Conversion logic
│       └── converter.go
├── conf/                  # Configuration files
│   └── currencies.json    # Cached currency list
├── main.go               # Application entry point
├── main_test.go          # Tests
└── CLAUDE.md            # Development instructions
```

### Testing

```bash
go test                  # Run all tests
go test -v              # Verbose output
```

### Building

```bash
go build -o conv        # Build binary
```

## Architecture

The application follows clean architecture principles:

- **Interface-based design**: `Converter` interface for testability
- **Separation of concerns**: CLI, currency logic, and conversion logic are separated
- **Type safety**: Custom `Currency` type with validation
- **Error handling**: Comprehensive error handling throughout
- **Caching**: Automatic currency list caching for performance