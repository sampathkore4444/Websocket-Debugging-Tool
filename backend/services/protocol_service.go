package services

import (
	"encoding/json"
	"fmt"
)

type ProtocolType string

const (
	ProtocolJSON     ProtocolType = "json"
	ProtocolProtobuf ProtocolType = "protobuf"
	ProtocolMsgPack  ProtocolType = "msgpack"
	ProtocolBinary   ProtocolType = "binary"
	ProtocolUnknown  ProtocolType = "unknown"
)

type ProtocolParser struct{}

func NewProtocolParser() *ProtocolParser {
	return &ProtocolParser{}
}

type ParsedMessage struct {
	Protocol ProtocolType      `json:"protocol"`
	Schema   map[string]interface{} `json:"schema,omitempty"`
	Fields   []FieldInfo        `json:"fields"`
	IsValid  bool              `json:"is_valid"`
	Error    string            `json:"error,omitempty"`
}

type FieldInfo struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	Children []FieldInfo `json:"children,omitempty"`
}

func (p *ProtocolParser) DetectProtocol(payload []byte) ProtocolType {
	// Try JSON
	if p.isValidJSON(string(payload)) {
		return ProtocolJSON
	}

	// Check for protobuf magic bytes
	if len(payload) >= 2 {
		// Common protobuf indicators
		if payload[0] == 0x08 || payload[0] == 0x0a {
			return ProtocolProtobuf
		}
	}

	// Check for msgpack
	if len(payload) > 0 {
		firstByte := payload[0]
		// Msgpack fixint, fixstr, fixarray, fixmap ranges
		if (firstByte >= 0x00 && firstByte <= 0x7f) ||
			(firstByte >= 0xa0 && firstByte <= 0xbf) ||
			(firstByte >= 0x90 && firstByte <= 0x9f) ||
			(firstByte >= 0x80 && firstByte <= 0x8f) {
			return ProtocolMsgPack
		}
	}

	// Check if binary
	if !p.isValidUTF8(payload) {
		return ProtocolBinary
	}

	return ProtocolUnknown
}

func (p *ProtocolParser) Parse(payload []byte) *ParsedMessage {
	result := &ParsedMessage{
		Fields: make([]FieldInfo, 0),
	}

	protocol := p.DetectProtocol(payload)
	result.Protocol = protocol

	switch protocol {
	case ProtocolJSON:
		return p.parseJSON(string(payload), result)
	case ProtocolBinary:
		result.IsValid = true
		result.Fields = []FieldInfo{
			{
				Name:  "raw",
				Type:  "binary",
				Value: fmt.Sprintf("%x", payload),
			},
		}
	case ProtocolProtobuf:
		result.IsValid = true
		result.Fields = []FieldInfo{
			{
				Name:  "protobuf_data",
				Type:  "binary",
				Value: fmt.Sprintf("%x", payload),
			},
		}
	default:
		result.IsValid = true
		result.Fields = []FieldInfo{
			{
				Name:  "text",
				Type:  "string",
				Value: string(payload),
			},
		}
	}

	return result
}

func (p *ProtocolParser) parseJSON(payload string, result *ParsedMessage) *ParsedMessage {
	var data interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		result.IsValid = false
		result.Error = err.Error()
		return result
	}

	result.IsValid = true
	result.Schema = p.extractSchema(data)
	result.Fields = p.extractFields(data, "")

	return result
}

func (p *ProtocolParser) extractSchema(data interface{}) map[string]interface{} {
	if obj, ok := data.(map[string]interface{}); ok {
		schema := make(map[string]interface{})
		for k, v := range obj {
			schema[k] = p.getType(v)
		}
		return schema
	}
	return nil
}

func (p *ProtocolParser) getType(v interface{}) string {
	switch v.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	case nil:
		return "null"
	default:
		return "unknown"
	}
}

func (p *ProtocolParser) extractFields(data interface{}, prefix string) []FieldInfo {
	fields := make([]FieldInfo, 0)

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			fieldName := key
			if prefix != "" {
				fieldName = prefix + "." + key
			}

			field := FieldInfo{
				Name:  fieldName,
				Type:  p.getType(value),
				Value: value,
			}

			if nested, ok := value.(map[string]interface{}); ok {
				field.Children = p.extractFields(nested, fieldName)
			}

			fields = append(fields, field)
		}
	case []interface{}:
		fields = append(fields, FieldInfo{
			Name:  prefix,
			Type:  "array",
			Value: fmt.Sprintf("[%d items]", len(v)),
		})
	}

	return fields
}

func (p *ProtocolParser) isValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func (p *ProtocolParser) isValidUTF8(b []byte) bool {
	i := 0
	for i < len(b) {
		if b[i] <= 0x7F {
			i++
		} else if (b[i] & 0xE0) == 0xC0 {
			if i+1 >= len(b) || (b[i+1]&0xC0) != 0x80 {
				return false
			}
			i += 2
		} else if (b[i] & 0xF0) == 0xE0 {
			if i+2 >= len(b) || (b[i+1]&0xC0) != 0x80 || (b[i+2]&0xC0) != 0x80 {
				return false
			}
			i += 3
		} else if (b[i] & 0xF8) == 0xF0 {
			if i+3 >= len(b) || (b[i+1]&0xC0) != 0x80 || (b[i+2]&0xC0) != 0x80 || (b[i+3]&0xC0) != 0x80 {
				return false
			}
			i += 4
		} else {
			return false
		}
	}
	return true
}
