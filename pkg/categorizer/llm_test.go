package categorizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSingleItemPrompt(t *testing.T) {
	// Arrange
	input := "Uber Ride"

	// Act
	prompt := buildSingleItemPrompt(input)

	// Assert
	assert.Contains(t, prompt, input, "Expected prompt to contain input")
	// Check if categories are mentioned (using the variable CategoryList)
	// Just check a sample category
	assert.Contains(t, prompt, string(CategoryRideHailing), "Expected prompt to contain category")
}

func TestBuildMultipleItemsPrompt(t *testing.T) {
	// Arrange
	inputs := []string{"Uber Ride", "Salary Deposit"}

	// Act
	prompt := buildMultipleItemsPrompt(inputs)

	// Assert
	for _, input := range inputs {
		assert.Contains(t, prompt, input, "Expected prompt to contain input")
	}
}

func TestParseSingleItemResponse(t *testing.T) {
	// Arrange
	originalInput := "Uber Ride"
	validJSON := `{
		"category": "ride_hailing",
		"merchant_or_counterparty": "Uber",
		"direction": "out",
		"amount": 25.50,
		"currency": "BRL",
		"confidence": 0.95
	}`

	// Act
	record, err := parseSingleItemResponse([]byte(validJSON), originalInput)
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, originalInput, record.Original)
	assert.Equal(t, CategoryRideHailing, record.Category)
	assert.Equal(t, 25.50, record.Amount)
}

func TestParseSingleItemResponse_InvalidJSON(t *testing.T) {
	// Arrange
	invalidJSON := `{"category": "ride_hailing",`

	// Act
	_, err := parseSingleItemResponse([]byte(invalidJSON), "input")

	// Assert
	assert.Error(t, err, "Expected error for invalid JSON")
}

func TestParseMultipleItemsResponse(t *testing.T) {
	// Arrange
	validJSON := `[
		{
			"original": "Uber Ride",
			"category": "ride_hailing",
			"merchant_or_counterparty": "Uber",
			"direction": "out",
			"amount": 25.50,
			"currency": "BRL",
			"confidence": 0.95
		},
		{
			"original": "Salary",
			"category": "salary",
			"merchant_or_counterparty": "Employer",
			"direction": "in",
			"amount": 5000.00,
			"currency": "BRL",
			"confidence": 0.99
		}
	]`

	// Act
	records, err := parseMultipleItemsResponse([]byte(validJSON))
	// Assert
	assert.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, CategoryRideHailing, records[0].Category)
	assert.Equal(t, CategorySalary, records[1].Category)
}

func TestParseMultipleItemsResponse_InvalidJSON(t *testing.T) {
	// Arrange
	invalidJSON := `[{"category": "ride_hailing"}`

	// Act
	_, err := parseMultipleItemsResponse([]byte(invalidJSON))

	// Assert
	assert.Error(t, err, "Expected error for invalid JSON")
}
