package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"categorizer/pkg/categorizer"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is not set")
	}

	ctx := context.Background()
	client, err := categorizer.NewLLMClient(ctx, apiKey)
	if err != nil {
		log.Fatalf("Failed to create LLM client: %v", err)
	}

	service := categorizer.NewService(client)

	inputs := []string{
		"Supermarket",
		"Banana",
		"Ifood",
		"Uber",
		"Transfer 100 to Alex",
		"Netflix Subscription",
		"Salary deposit",
		"Dentist appointment",
	}

	fmt.Printf("Categorizing %d records in a single batch...\n", len(inputs))
	results, err := service.CategorizeMany(ctx, inputs)
	if err != nil {
		log.Fatalf("Categorization failed: %v", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		log.Fatalf("Failed to encode results: %v", err)
	}
}
