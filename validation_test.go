package ga4m

import "testing"

func TestValidateEventName_Valid(t *testing.T) {
	eventName := "validEventName_123"

	err := validateEventName(eventName)
	if err != nil {
		t.Errorf("Expected event name to be valid, got error: %v", err)
	}
}

func TestValidateEventName_Invalid(t *testing.T) {
	eventName := "1InvalidEventName"

	err := validateEventName(eventName)
	if err == nil {
		t.Errorf("Expected event name to be invalid, got no error")
	}
}

func TestValidateParams_Valid(t *testing.T) {
	params := map[string]string{
		"param_one": "value1",
		"paramTwo":  "value2",
	}

	err := validateParams(params)
	if err != nil {
		t.Errorf("Expected parameters to be valid, got error: %v", err)
	}
}

func TestValidateParams_Invalid(t *testing.T) {
	params := map[string]string{
		"param-one": "value1", // Invalid character '-'
	}

	err := validateParams(params)
	if err == nil {
		t.Errorf("Expected parameters to be invalid, got no error")
	}
}
