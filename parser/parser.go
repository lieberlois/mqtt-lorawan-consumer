package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mqtt_consumer/config"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrInvalidJsonFormat = errors.New("error trying to unmarshal JSON string")
	ErrCastFailed        = errors.New("error parsing to map[string]interface{}")
	ErrInvalidPayload    = errors.New("invalid payload format")
)

type Parser struct {
	config config.Parser
}

func NewParser(cfg config.Parser) *Parser {
	return &Parser{config: cfg}
}

func (parser *Parser) StringToJson(jsonIn string) (map[string]interface{}, error) {
	var parsed interface{}
	bytes := []byte(jsonIn)
	err := json.Unmarshal(bytes, &parsed)

	if err != nil {
		return nil, ErrInvalidJsonFormat
	}

	result, success := parsed.(map[string]interface{})

	if !success {
		return nil, ErrCastFailed
	}

	return result, nil
}

func (parser *Parser) JsonToInfluxLineProtocol(data map[string]interface{}) (string, error) {
	// Target format
	// measurement,tag1=val1,tag2=val2 data1=val1,data2=val2

	measurement := data[parser.config.MeasurementKey]

	// Tags
	tagString := ParseListToLineFormat(data, parser.config.TagKeys)

	// Payload
	payloadString := ParseMapToLineFormat(data, parser.config.ValuesKey)

	if len(payloadString) == 0 {
		log.Println("Invalid payload")
		return "", ErrInvalidPayload
	}

	return fmt.Sprintf("%s%s %s", measurement, tagString, payloadString[1:]), nil
}

func ParseMapToLineFormat(data map[string]interface{}, key string) string {
	var result string

	if val, ok := data[key]; ok {
		dataMap, success := val.(map[string]interface{})

		// Sort alphabetically
		keys := make([]string, 0, len(dataMap))
		for k := range dataMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		if success {
			for _, key := range keys {
				valString := fmt.Sprintf("%v", dataMap[key])

				if strings.Contains(valString, " ") {
					valString = strconv.Quote(valString)
				}

				result += fmt.Sprintf(",%s=%s", key, valString)
			}
		}
	}
	return result
}

func ParseListToLineFormat(data map[string]interface{}, keys []string) string {
	var result string

	sort.Strings(keys)

	for _, key := range keys {
		if val, ok := data[key]; ok {
			valString := fmt.Sprintf("%v", val)

			if strings.Contains(valString, " ") {
				valString = strconv.Quote(valString)
			}
			result += fmt.Sprintf(",%s=%s", key, valString)
		}
	}

	return result
}
