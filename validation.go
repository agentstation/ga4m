package ga4m

import "fmt"

const (
	maxEventNameLength  = 40
	maxParamNameLength  = 40
	maxParamValueLength = 100
	maxEventParams      = 25
)

func validateEventName(name string) error {
	if len(name) > maxEventNameLength {
		return fmt.Errorf("event name must be %d characters or fewer", maxEventNameLength)
	}
	if len(name) == 0 {
		return fmt.Errorf("event name cannot be empty")
	}
	// Check first character is a letter
	if !isLetter(name[0]) {
		return fmt.Errorf("event name must start with a letter")
	}
	// Check remaining characters
	for i := 1; i < len(name); i++ {
		if !isAlphanumericOrUnderscore(name[i]) {
			return fmt.Errorf("event name must contain only alphanumeric characters and underscores")
		}
	}
	return nil
}

func validateParams(params map[string]string) error {
	if len(params) > maxEventParams {
		return fmt.Errorf("events can have a maximum of %d parameters", maxEventParams)
	}

	for name, value := range params {
		if len(name) > maxParamNameLength {
			return fmt.Errorf("parameter name '%s' exceeds maximum length of %d", name, maxParamNameLength)
		}
		if len(name) == 0 {
			return fmt.Errorf("parameter name cannot be empty")
		}
		// Check first character is a letter
		if !isLetter(name[0]) {
			return fmt.Errorf("parameter name '%s' must start with a letter", name)
		}
		// Check remaining characters
		for i := 1; i < len(name); i++ {
			if !isAlphanumericOrUnderscore(name[i]) {
				return fmt.Errorf("parameter name '%s' must contain only alphanumeric characters and underscores", name)
			}
		}

		if len(value) > maxParamValueLength {
			return fmt.Errorf("parameter value for '%s' exceeds maximum length of %d", name, maxParamValueLength)
		}
	}
	return nil
}

// Helper functions
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isAlphanumericOrUnderscore(c byte) bool {
	return isLetter(c) || (c >= '0' && c <= '9') || c == '_'
}
