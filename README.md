# LLM-powered Record Categorizer

A Go-based system that uses Google's Gemini Flash model to categorize transaction records into structured data.

## Features

- **LLM Integration**: Uses **Gemini Flash** via the official `google.golang.org/genai` SDK.
- **Structured Output**: Returns strict JSON with fields like `category`, `amount`, `merchant`, etc.
- **Single-Request Batch Processing**: Efficiently categorizes multiple records in a single LLM call to save time and quota.

## Prerequisites

- [Go](https://go.dev/dl/) 1.25+
- A Google Cloud Project with the Gemini API enabled.
- an API Key for Gemini.

## Setup

1.  Clone the repository.
2.  Install dependencies:
    ```bash
    go mod tidy
    ```

## Usage

1.  Create a `.env` file:

    ```bash
    cp .env.example .env
    # Edit .env and paste your API key
    ```

2.  Run the categorizer:
    ```bash
    make run
    # Or directly: go run cmd/categorizer/main.go
    ```

## Testing

Run the test suite:

```bash
make test
```

## Build

Build the binary:

```bash
make build
```

The binary will be created at `bin/categorizer`.

### Example Input

The CLI processes a predefined list of records:

- "Supermarket"
- "Banana"
- "Ifood"
- "Uber"
- "Transfer 100 to Alex"
- "Netflix Subscription"
- "Salary deposit"
- "Dentist appointment"

### Example Output

```json
[
  {
    "original": "Transfer 100 to Alex",
    "category": "transfer",
    "merchant_or_counterparty": "Alex",
    "direction": "out",
    "amount": 100,
    "currency": "BRL",
    "confidence": 0.87
  }
]
```

## Project Structure

- `pkg/categorizer`: Core domain logic and types.
  - `types.go`: Data structures.
  - `llm.go`: Gemini client integration (using `google.golang.org/genai`).
  - `service.go`: Business logic service.
- `cmd/categorizer`: CLI entry point.
