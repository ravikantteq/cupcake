package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Validator interface for message validation
type Validator interface {
	Validate(value interface{}) bool
}

// ValidationResult represents the result of a validation
type ValidationResult struct {
	Valid   bool     `json:"valid"`
	Errors  []string `json:"errors,omitempty"`
	Message string   `json:"message,omitempty"`
}

// MessageValidator handles intelligent message validation
type MessageValidator struct {
	stepValues map[string]interface{} // Store values from previous steps
}

// NewMessageValidator creates a new message validator
func NewMessageValidator() *MessageValidator {
	return &MessageValidator{
		stepValues: make(map[string]interface{}),
	}
}

// SetStepValue stores a value from a previous step
func (mv *MessageValidator) SetStepValue(stepID, field string, value interface{}) {
	key := fmt.Sprintf("%s.%s", stepID, field)
	mv.stepValues[key] = value
}

// ValidateMessage validates a message against expected patterns
func (mv *MessageValidator) ValidateMessage(actual, expected map[string]interface{}) *ValidationResult {
	result := &ValidationResult{Valid: true, Errors: []string{}}

	for key, expectedValue := range expected {
		actualValue, exists := actual[key]
		if !exists {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("field '%s' is missing", key))
			continue
		}

		if err := mv.validateField(key, actualValue, expectedValue); err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, err.Error())
		}
	}

	if result.Valid {
		result.Message = "All validations passed"
	} else {
		result.Message = fmt.Sprintf("Validation failed with %d errors", len(result.Errors))
	}

	return result
}

// validateField validates a single field
func (mv *MessageValidator) validateField(fieldName string, actual, expected interface{}) error {
	// Handle nested objects
	if expectedMap, ok := expected.(map[string]interface{}); ok {
		actualMap, ok := actual.(map[string]interface{})
		if !ok {
			return fmt.Errorf("field '%s' should be an object", fieldName)
		}

		for nestedKey, nestedExpected := range expectedMap {
			nestedActual, exists := actualMap[nestedKey]
			if !exists {
				return fmt.Errorf("nested field '%s.%s' is missing", fieldName, nestedKey)
			}

			if err := mv.validateField(fmt.Sprintf("%s.%s", fieldName, nestedKey), nestedActual, nestedExpected); err != nil {
				return err
			}
		}
		return nil
	}

	// Handle arrays
	if expectedArray, ok := expected.([]interface{}); ok {
		actualArray, ok := actual.([]interface{})
		if !ok {
			return fmt.Errorf("field '%s' should be an array", fieldName)
		}

		if len(expectedArray) > 0 {
			// Validate each element against the first expected pattern
			expectedPattern := expectedArray[0]
			for i, actualItem := range actualArray {
				if err := mv.validateField(fmt.Sprintf("%s[%d]", fieldName, i), actualItem, expectedPattern); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Handle validation functions
	if expectedStr, ok := expected.(string); ok {
		return mv.validateWithFunction(fieldName, actual, expectedStr)
	}

	// Direct value comparison
	if !reflect.DeepEqual(actual, expected) {
		return fmt.Errorf("field '%s' expected '%v' but got '%v'", fieldName, expected, actual)
	}

	return nil
}

// validateWithFunction validates using built-in validation functions
func (mv *MessageValidator) validateWithFunction(fieldName string, actual interface{}, pattern string) error {
	// Check for function patterns
	if strings.HasSuffix(pattern, "()") {
		functionName := strings.TrimSuffix(pattern, "()")
		return mv.validateBuiltinFunction(fieldName, actual, functionName, "")
	}

	// Check for function patterns with parameters
	if strings.Contains(pattern, "(") && strings.Contains(pattern, ")") {
		parts := strings.SplitN(pattern, "(", 2)
		functionName := parts[0]
		params := strings.TrimSuffix(parts[1], ")")
		return mv.validateBuiltinFunction(fieldName, actual, functionName, params)
	}

	// Check for step references
	if strings.HasPrefix(pattern, "match(") && strings.HasSuffix(pattern, ")") {
		reference := strings.TrimPrefix(pattern, "match(")
		reference = strings.TrimSuffix(reference, ")")
		return mv.validateStepReference(fieldName, actual, reference)
	}

	// Direct string comparison
	actualStr := fmt.Sprintf("%v", actual)
	if actualStr != pattern {
		return fmt.Errorf("field '%s' expected '%s' but got '%s'", fieldName, pattern, actualStr)
	}

	return nil
}

// validateBuiltinFunction validates using built-in functions
func (mv *MessageValidator) validateBuiltinFunction(fieldName string, actual interface{}, functionName, params string) error {
	switch functionName {
	case "uuid":
		return mv.validateUUID(fieldName, actual)
	case "timestamp":
		return mv.validateTimestamp(fieldName, actual)
	case "number":
		return mv.validateNumber(fieldName, actual, params)
	case "enum":
		return mv.validateEnum(fieldName, actual, params)
	case "regex":
		return mv.validateRegex(fieldName, actual, params)
	case "any":
		return nil // any() always passes
	case "string":
		return mv.validateString(fieldName, actual, params)
	case "array":
		return mv.validateArray(fieldName, actual, params)
	default:
		return fmt.Errorf("unknown validation function: %s", functionName)
	}
}

// validateUUID validates UUID format
func (mv *MessageValidator) validateUUID(fieldName string, actual interface{}) error {
	str, ok := actual.(string)
	if !ok {
		return fmt.Errorf("field '%s' should be a string for UUID validation", fieldName)
	}

	if _, err := uuid.Parse(str); err != nil {
		return fmt.Errorf("field '%s' is not a valid UUID: %s", fieldName, str)
	}

	return nil
}

// validateTimestamp validates timestamp format
func (mv *MessageValidator) validateTimestamp(fieldName string, actual interface{}) error {
	switch v := actual.(type) {
	case string:
		// Try parsing various timestamp formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05Z",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}

		for _, format := range formats {
			if _, err := time.Parse(format, v); err == nil {
				return nil
			}
		}

		// Try parsing as Unix timestamp
		if _, err := strconv.ParseInt(v, 10, 64); err == nil {
			return nil
		}

		return fmt.Errorf("field '%s' is not a valid timestamp: %s", fieldName, v)

	case int64, int, float64:
		// Unix timestamp
		return nil

	default:
		return fmt.Errorf("field '%s' should be a string or number for timestamp validation", fieldName)
	}
}

// validateNumber validates number with optional range
func (mv *MessageValidator) validateNumber(fieldName string, actual interface{}, params string) error {
	var num float64
	switch v := actual.(type) {
	case int:
		num = float64(v)
	case int64:
		num = float64(v)
	case float32:
		num = float64(v)
	case float64:
		num = v
	default:
		return fmt.Errorf("field '%s' should be a number", fieldName)
	}

	if params == "" {
		return nil // No range specified
	}

	// Parse parameters: min=X, max=Y
	parts := strings.Split(params, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "min=") {
			minStr := strings.TrimPrefix(part, "min=")
			min, err := strconv.ParseFloat(minStr, 64)
			if err != nil {
				return fmt.Errorf("invalid min parameter: %s", minStr)
			}
			if num < min {
				return fmt.Errorf("field '%s' value %v is less than minimum %v", fieldName, num, min)
			}
		} else if strings.HasPrefix(part, "max=") {
			maxStr := strings.TrimPrefix(part, "max=")
			max, err := strconv.ParseFloat(maxStr, 64)
			if err != nil {
				return fmt.Errorf("invalid max parameter: %s", maxStr)
			}
			if num > max {
				return fmt.Errorf("field '%s' value %v is greater than maximum %v", fieldName, num, max)
			}
		}
	}

	return nil
}

// validateEnum validates against enumerated values
func (mv *MessageValidator) validateEnum(fieldName string, actual interface{}, params string) error {
	actualStr := fmt.Sprintf("%v", actual)

	values := strings.Split(params, ",")
	for _, value := range values {
		if strings.TrimSpace(value) == actualStr {
			return nil
		}
	}

	return fmt.Errorf("field '%s' value '%s' is not in allowed values: [%s]", fieldName, actualStr, params)
}

// validateRegex validates against regex pattern
func (mv *MessageValidator) validateRegex(fieldName string, actual interface{}, params string) error {
	str, ok := actual.(string)
	if !ok {
		return fmt.Errorf("field '%s' should be a string for regex validation", fieldName)
	}

	regex, err := regexp.Compile(params)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %s", params)
	}

	if !regex.MatchString(str) {
		return fmt.Errorf("field '%s' value '%s' does not match pattern '%s'", fieldName, str, params)
	}

	return nil
}

// validateString validates string with optional parameters
func (mv *MessageValidator) validateString(fieldName string, actual interface{}, params string) error {
	_, ok := actual.(string)
	if !ok {
		return fmt.Errorf("field '%s' should be a string", fieldName)
	}

	if params == "" {
		return nil
	}

	// Handle regex parameter
	if strings.HasPrefix(params, "regex=") {
		pattern := strings.TrimPrefix(params, "regex=")
		return mv.validateRegex(fieldName, actual, pattern)
	}

	return nil
}

// validateArray validates array with optional element validation
func (mv *MessageValidator) validateArray(fieldName string, actual interface{}, params string) error {
	arr, ok := actual.([]interface{})
	if !ok {
		return fmt.Errorf("field '%s' should be an array", fieldName)
	}

	if params == "" {
		return nil
	}

	// Parse parameters: type, minLength=X, maxLength=Y
	parts := strings.Split(params, ",")
	var elementType string
	var minLength, maxLength int = -1, -1

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "minLength=") {
			lenStr := strings.TrimPrefix(part, "minLength=")
			if l, err := strconv.Atoi(lenStr); err == nil {
				minLength = l
			}
		} else if strings.HasPrefix(part, "maxLength=") {
			lenStr := strings.TrimPrefix(part, "maxLength=")
			if l, err := strconv.Atoi(lenStr); err == nil {
				maxLength = l
			}
		} else {
			elementType = part
		}
	}

	// Check length constraints
	if minLength >= 0 && len(arr) < minLength {
		return fmt.Errorf("field '%s' array length %d is less than minimum %d", fieldName, len(arr), minLength)
	}
	if maxLength >= 0 && len(arr) > maxLength {
		return fmt.Errorf("field '%s' array length %d is greater than maximum %d", fieldName, len(arr), maxLength)
	}

	// Validate element types if specified
	if elementType != "" {
		for i, element := range arr {
			if err := mv.validateBuiltinFunction(fmt.Sprintf("%s[%d]", fieldName, i), element, elementType, ""); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateStepReference validates against a value from a previous step
func (mv *MessageValidator) validateStepReference(fieldName string, actual interface{}, reference string) error {
	expectedValue, exists := mv.stepValues[reference]
	if !exists {
		return fmt.Errorf("step reference '%s' not found for field '%s'", reference, fieldName)
	}

	if !reflect.DeepEqual(actual, expectedValue) {
		return fmt.Errorf("field '%s' expected '%v' (from %s) but got '%v'", fieldName, expectedValue, reference, actual)
	}

	return nil
}
