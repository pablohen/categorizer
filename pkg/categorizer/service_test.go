package categorizer

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockCategorizer is a mock implementation of the Categorizer interface
type MockCategorizer struct {
	CategorizeOneFunc  func(ctx context.Context, input string) (*TransactionRecord, error)
	CategorizeManyFunc func(ctx context.Context, inputs []string) ([]TransactionRecord, error)
}

func (m *MockCategorizer) CategorizeOne(ctx context.Context, input string) (*TransactionRecord, error) {
	if m.CategorizeOneFunc != nil {
		return m.CategorizeOneFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockCategorizer) CategorizeMany(ctx context.Context, inputs []string) ([]TransactionRecord, error) {
	if m.CategorizeManyFunc != nil {
		return m.CategorizeManyFunc(ctx, inputs)
	}
	return nil, nil
}

func TestService_CategorizeOne(t *testing.T) {
	// Arrange
	expectedRecord := &TransactionRecord{
		Original: "Uber",
		Category: CategoryRideHailing,
	}
	mock := &MockCategorizer{
		CategorizeOneFunc: func(ctx context.Context, input string) (*TransactionRecord, error) {
			if input == "Uber" {
				return expectedRecord, nil
			}
			return nil, errors.New("unexpected input")
		},
	}
	service := NewService(mock)
	ctx := context.Background()

	// Act
	result, err := service.CategorizeOne(ctx, "Uber")
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord, result)
}

func TestService_CategorizeOne_Error(t *testing.T) {
	// Arrange
	expectedErr := errors.New("categorization failed")
	mock := &MockCategorizer{
		CategorizeOneFunc: func(ctx context.Context, input string) (*TransactionRecord, error) {
			return nil, expectedErr
		},
	}
	service := NewService(mock)
	ctx := context.Background()

	// Act
	_, err := service.CategorizeOne(ctx, "Uber")

	// Assert
	assert.ErrorIs(t, err, expectedErr)
}

func TestService_CategorizeMany(t *testing.T) {
	// Arrange
	expectedRecords := []TransactionRecord{
		{Original: "Uber", Category: CategoryRideHailing},
		{Original: "Salary", Category: CategorySalary},
	}

	expectedRecordsOriginals := make([]string, len(expectedRecords))
	for i, record := range expectedRecords {
		expectedRecordsOriginals[i] = record.Original
	}

	mock := &MockCategorizer{
		CategorizeManyFunc: func(ctx context.Context, inputs []string) ([]TransactionRecord, error) {
			if slices.Equal(expectedRecordsOriginals, inputs) {
				return expectedRecords, nil
			}
			return nil, errors.New("unexpected inputs")
		},
	}
	service := NewService(mock)
	ctx := context.Background()
	inputs := []string{"Uber", "Salary"}

	// Act
	results, err := service.CategorizeMany(ctx, inputs)
	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, len(expectedRecords))
}

func TestService_CategorizeMany_Error(t *testing.T) {
	// Arrange
	expectedErr := errors.New("batch categorization failed")
	mock := &MockCategorizer{
		CategorizeManyFunc: func(ctx context.Context, inputs []string) ([]TransactionRecord, error) {
			return nil, expectedErr
		},
	}
	service := NewService(mock)
	ctx := context.Background()

	// Act
	_, err := service.CategorizeMany(ctx, []string{"Uber"})

	// Assert
	assert.ErrorIs(t, err, expectedErr)
}
