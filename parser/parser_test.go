package parser

import (
	"mqtt_consumer/config"
	"reflect"
	"testing"
)

func TestNewParser(t *testing.T) {
	parserCfg := config.Parser{}
	parser := NewParser(parserCfg)

	if parser == nil {
		t.Errorf("Error trying to create a new parser")
	}
}

func TestStringToJson(t *testing.T) {
	parser := Parser{config.Parser{}}

	t.Run("valid json should be parsed", func(t *testing.T) {
		s := `{"app_id":"example-tracker", "key":"value", "payload_fields":{"degreesC":21.39892578125,"humidity":31.298828125}}`

		expected := map[string]interface{}{
			"app_id":         "example-tracker",
			"key":            "value",
			"payload_fields": map[string]interface{}{"degreesC": 21.39892578125, "humidity": 31.298828125},
		}
		actual, _ := parser.StringToJson(s)

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v got %v", expected, actual)
		}
	})

	t.Run("invalid JSON should throw error", func(t *testing.T) {
		s := `{"app_id":"example-tracker", "key":value", "payload_fields":{"degreesC":21.39892578125,"humidity":31.298828125}}`

		_, err := parser.StringToJson(s)

		if err == ErrCastFailed {
			t.Errorf("expected error when parsing invalid JSON string %s", s)
		}
	})
}

func TestJsonToInfluxLineProtocol(t *testing.T) {
	cfg := config.Parser{
		MeasurementKey: "dev_id",
		TagKeys:        []string{"attribute1", "attribute2"},
		ValuesKey:      "payload_fields",
	}
	parser := Parser{cfg}

	t.Run("attributes and payload both parsed correctly", func(t *testing.T) {
		json := map[string]interface{}{
			"app_id":         "some-name",
			"dev_id":         "device_name",
			"key":            "value",
			"attribute1":     "tagval1",
			"attribute2":     "tagval 2",
			"payload_fields": map[string]interface{}{"value1": 21.39892578125, "value2": "hello world"},
		}

		expected := "device_name,attribute1=tagval1,attribute2=\"tagval 2\" value1=21.39892578125,value2=\"hello world\""
		actual, _ := parser.JsonToInfluxLineProtocol(json)

		if expected != actual {
			t.Errorf("expected %s got %s", expected, actual)
		}
	})

	t.Run("only payload gets parsed", func(t *testing.T) {
		json := map[string]interface{}{
			"app_id":         "some-name",
			"dev_id":         "device_name",
			"key":            "value",
			"payload_fields": map[string]interface{}{"value1": 21.39892578125, "value2": "hello"},
		}

		expected := "device_name value1=21.39892578125,value2=hello"
		actual, _ := parser.JsonToInfluxLineProtocol(json)

		if expected != actual {
			t.Errorf("expected %s got %s", expected, actual)
		}
	})

	t.Run("missing payload results in error", func(t *testing.T) {
		json := map[string]interface{}{
			"app_id": "some-name",
			"dev_id": "device_name",
			"key":    "value",
		}

		_, err := parser.JsonToInfluxLineProtocol(json)

		if ErrInvalidPayload != err {
			t.Errorf("expected %v got %v", ErrInvalidPayload, err)
		}
	})
}
