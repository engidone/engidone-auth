package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"engidoneauth/util/crypto"
)

// TableFilter handles filtering and validation of table names
type TableFilter struct {
	SpecifiedTables []string
	AvailableTables []string
}

// NewTableFilter creates a new table filter
func NewTableFilter(tableFlag string, availableTables []string) *TableFilter {
	return &TableFilter{
		SpecifiedTables: parseTableNames(tableFlag),
		AvailableTables: availableTables,
	}
}

// ShouldProcessTable determines if a table should be processed based on the filter
func (tf *TableFilter) ShouldProcessTable(tableName string) bool {
	// If no tables specified, process all
	if len(tf.SpecifiedTables) == 0 {
		return true
	}

	// Check if table name is in the specified list
	for _, specified := range tf.SpecifiedTables {
		if strings.EqualFold(tableName, specified) {
			return true
		}
	}
	return false
}

// ValidateSpecifiedTables checks if all specified tables exist in available tables
func (tf *TableFilter) ValidateSpecifiedTables() []string {
	if len(tf.SpecifiedTables) == 0 {
		return nil
	}

	var invalidTables []string
	for _, specified := range tf.SpecifiedTables {
		found := false
		for _, available := range tf.AvailableTables {
			if strings.EqualFold(specified, available) {
				found = true
				break
			}
		}
		if !found {
			invalidTables = append(invalidTables, specified)
		}
	}

	return invalidTables
}

// GetTargetTables returns the list of tables that will be processed
func (tf *TableFilter) GetTargetTables() []string {
	if len(tf.SpecifiedTables) > 0 {
		return tf.SpecifiedTables
	}
	return tf.AvailableTables
}

// parseTableNames parses comma-separated table names from flag
func parseTableNames(tableFlag string) []string {
	if tableFlag == "" {
		return []string{}
	}

	tables := strings.Split(tableFlag, ",")
	var result []string
	for _, table := range tables {
		trimmed := strings.TrimSpace(table)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// StringConverter handles conversion of various types to strings
type StringConverter struct{}

// ToString converts any value to its string representation
func (sc *StringConverter) ToString(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// FunctionProcessor handles processing of special function strings in seed data
type FunctionProcessor struct {
	converter *StringConverter
}

// NewFunctionProcessor creates a new function processor
func NewFunctionProcessor() *FunctionProcessor {
	return &FunctionProcessor{
		converter: &StringConverter{},
	}
}

// ProcessValue processes a value that may contain function calls
func (fp *FunctionProcessor) ProcessValue(value any) (any, error) {
	valueStr := fp.converter.ToString(value)
	splitValues := strings.Split(valueStr, "|")

	if len(splitValues) == 1 {
		return fp.processSingleFunction(splitValues[0])
	} else if len(splitValues) == 2 {
		return fp.processFunctionWithArg(splitValues[0], splitValues[1])
	}

	return valueStr, nil
}

// processSingleFunction processes functions without arguments
func (fp *FunctionProcessor) processSingleFunction(functionName string) (any, error) {
	switch functionName {
	case "uuid_string":
		return handleUUIDString(), nil
	default:
		return functionName, nil
	}
}

// processFunctionWithArg processes functions with arguments
func (fp *FunctionProcessor) processFunctionWithArg(functionName, arg string) (any, error) {
	switch functionName {
	case "uuid_parse":
		return handleUUIDParse(arg), nil
	case "hash_password":
		return handleHashPassword(arg), nil
	case "add_minutes_to_now":
		minutes, err := strconv.Atoi(arg)
		if err != nil {
			minutes = 0
		}
		return handleAddMinutesToNow(minutes), nil
	default:
		return arg, nil
	}
}

// ProcessingResult contains the results of processing seed items
type ProcessingResult struct {
	TableName     string
	FileName      string
	TotalInserted int
	ItemCount     int
	SQL           string
	Values        [][]any
}

// SeedProcessor handles the main seed processing logic
type SeedProcessor struct {
	functionProcessor *FunctionProcessor
	converter         *StringConverter
}

// NewSeedProcessor creates a new seed processor
func NewSeedProcessor() *SeedProcessor {
	return &SeedProcessor{
		functionProcessor: NewFunctionProcessor(),
		converter:         &StringConverter{},
	}
}

// ProcessSeedFile processes a single seed file and prepares data for insertion
func (sp *SeedProcessor) ProcessSeedFile(seedFile *SeedFile) (*ProcessingResult, error) {
	fields, sequences := GetFieldsAndSequences(seedFile.Data.Items)

	// Build SQL statement
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)",
		seedFile.Data.Name, fields, sequences)

	// Process all items
	allValues := make([][]any, 0, len(seedFile.Data.Items))
	for _, item := range seedFile.Data.Items {
		values, err := sp.processItemValues(item, fields)
		if err != nil {
			return nil, fmt.Errorf("error processing item values: %w", err)
		}
		allValues = append(allValues, values)
	}

	return &ProcessingResult{
		TableName: seedFile.Data.Name,
		FileName:  seedFile.FileName,
		ItemCount: len(allValues),
		SQL:       sql,
		Values:    allValues,
	}, nil
}

// processItemValues processes all values in a single seed item
func (sp *SeedProcessor) processItemValues(item map[string]any, fieldsStr string) ([]any, error) {
	fieldOrder := strings.Split(fieldsStr, ", ")
	values := make([]any, 0, len(fieldOrder))

	for _, field := range fieldOrder {
		value := item[field]
		processedValue, err := sp.functionProcessor.ProcessValue(value)
		if err != nil {
			return nil, fmt.Errorf("error processing field %s: %w", field, err)
		}
		values = append(values, processedValue)
	}

	return values, nil
}

// Function handlers

// handleUUIDString generates a new UUID string
func handleUUIDString() string {
	return uuid.NewString()
}

// handleUUIDParse parses and validates a UUID string
func handleUUIDParse(id string) string {
	if id == "" {
		return uuid.NewString()
	}
	return uuid.MustParse(id).String()
}

// handleHashPassword hashes a password using bcrypt
func handleHashPassword(password string) string {
	return crypto.HashPassword(password)
}

// handleAddMinutesToNow adds the specified minutes to the current time
func handleAddMinutesToNow(minutes int) time.Time {
	return time.Now().Add(time.Duration(minutes) * time.Minute)
}