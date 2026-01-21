package categorizer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

const (
	model            = "gemini-3-flash-preview"
	singleItemPrompt = `
Categorize the following transaction record into one of these categories: %s.
Currency is always BRL.
Direction should be "in" or "out".

Input Record:
%q
`
	multipleItemsPrompt = `
Categorize the following list of transaction records into one of these categories: %s.
Currency is always BRL.
Direction should be "in" or "out".
Return a JSON array of objects, one for each input record, in the same order.

Input Records:
%s
`
)

type LLMClient struct {
	client *genai.Client
}

func NewLLMClient(ctx context.Context, apiKey string) (*LLMClient, error) {
	cfg := &genai.ClientConfig{
		APIKey: apiKey,
	}

	client, err := genai.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &LLMClient{
		client: client,
	}, nil
}

func (c *LLMClient) CategorizeOne(ctx context.Context, input string) (*TransactionRecord, error) {
	prompt := buildSingleItemPrompt(input)

	respSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"original":                 {Type: genai.TypeString},
			"category":                 {Type: genai.TypeString},
			"merchant_or_counterparty": {Type: genai.TypeString},
			"direction":                {Type: genai.TypeString},
			"amount":                   {Type: genai.TypeNumber},
			"currency":                 {Type: genai.TypeString},
			"confidence":               {Type: genai.TypeNumber},
		},
		Required: []string{"original", "category", "merchant_or_counterparty", "direction", "amount", "currency", "confidence"},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   respSchema,
	}

	resp, err := c.client.Models.GenerateContent(ctx, model, genai.Text(prompt), config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("no content returned")
	}

	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			sb.WriteString(part.Text)
		}
	}
	jsonBytes := []byte(sb.String())

	return parseSingleItemResponse(jsonBytes, input)
}

func (c *LLMClient) CategorizeMany(ctx context.Context, inputs []string) ([]TransactionRecord, error) {
	prompt := buildMultipleItemsPrompt(inputs)

	itemSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"original":                 {Type: genai.TypeString},
			"category":                 {Type: genai.TypeString},
			"merchant_or_counterparty": {Type: genai.TypeString},
			"direction":                {Type: genai.TypeString},
			"amount":                   {Type: genai.TypeNumber},
			"currency":                 {Type: genai.TypeString},
			"confidence":               {Type: genai.TypeNumber},
		},
		Required: []string{"original", "category", "merchant_or_counterparty", "direction", "amount", "currency", "confidence"},
	}

	respSchema := &genai.Schema{
		Type:  genai.TypeArray,
		Items: itemSchema,
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   respSchema,
	}

	resp, err := c.client.Models.GenerateContent(ctx, model, genai.Text(prompt), config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("no content returned")
	}

	var sb strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			sb.WriteString(part.Text)
		}
	}
	jsonBytes := []byte(sb.String())

	return parseMultipleItemsResponse(jsonBytes)
}

func buildSingleItemPrompt(input string) string {
	return fmt.Sprintf(singleItemPrompt, CategoryList, input)
}

func buildMultipleItemsPrompt(inputs []string) string {
	inputStr := strings.Join(inputs, "\n")
	return fmt.Sprintf(multipleItemsPrompt, CategoryList, inputStr)
}

func parseSingleItemResponse(jsonBytes []byte, input string) (*TransactionRecord, error) {
	var record TransactionRecord
	if err := json.Unmarshal(jsonBytes, &record); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, raw response: %s", err, string(jsonBytes))
	}

	record.Original = input

	return &record, nil
}

func parseMultipleItemsResponse(jsonBytes []byte) ([]TransactionRecord, error) {
	var records []TransactionRecord
	if err := json.Unmarshal(jsonBytes, &records); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, raw response: %s", err, string(jsonBytes))
	}

	return records, nil
}
