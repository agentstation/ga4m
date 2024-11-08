package ga4m

import (
	"fmt"
	"regexp"
)

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
	matched, err := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, name)
	if err != nil {
		return fmt.Errorf("error validating event name: %w", err)
	}
	if !matched {
		return fmt.Errorf("event name must start with a letter and contain only alphanumeric characters and underscores")
	}
	return nil
}

func validateParams(params map[string]interface{}) error {
	if len(params) > maxEventParams {
		return fmt.Errorf("events can have a maximum of %d parameters", maxEventParams)
	}

	for name, value := range params {
		if len(name) > maxParamNameLength {
			return fmt.Errorf("parameter name '%s' exceeds maximum length of %d", name, maxParamNameLength)
		}
		matched, err := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, name)
		if err != nil {
			return fmt.Errorf("error validating parameter name: %w", err)
		}
		if !matched {
			return fmt.Errorf("parameter name '%s' must start with a letter and contain only alphanumeric characters and underscores", name)
		}

		strValue := fmt.Sprintf("%v", value)
		if len(strValue) > maxParamValueLength {
			return fmt.Errorf("parameter value for '%s' exceeds maximum length of %d", name, maxParamValueLength)
		}
	}
	return nil
}
