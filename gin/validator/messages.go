package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var ValidationMessages = map[string]string{
	"required": "Field '{field}' is required",
	"oneof":    "Field '{field}' must be one of: {param}",
	"min":      "Field '{field}' must be greater than or equal to {param}",
	"max":      "Field '{field}' must be less than or equal to {param}",
	"email":    "Field '{field}' must be a valid email address",
	"url":      "Field '{field}' must be a valid URL",
	"len":      "Field '{field}' length must be {param}",
	"gt":       "Field '{field}' must be greater than {param}",
	"gte":      "Field '{field}' must be greater than or equal to {param}",
	"lt":       "Field '{field}' must be less than {param}",
	"lte":      "Field '{field}' must be less than or equal to {param}",
	"eq":       "Field '{field}' must be equal to {param}",
	"ne":       "Field '{field}' must not be equal to {param}",
	"numeric":  "Field '{field}' must be numeric",
	"alpha":    "Field '{field}' must contain only letters",
	"alphanum": "Field '{field}' must contain only letters and numbers",
	"boolean":  "Field '{field}' must be a boolean value",
	"datetime": "Field '{field}' must be a valid datetime format",
	"uuid":     "Field '{field}' must be a valid UUID format",
	"json":     "Field '{field}' must be valid JSON",
	"base64":   "Field '{field}' must be valid Base64 encoded",
	"hexcolor": "Field '{field}' must be a valid hexadecimal color",
	"rgb":      "Field '{field}' must be a valid RGB color",
	"rgba":     "Field '{field}' must be a valid RGBA color",
	"ip":       "Field '{field}' must be a valid IP address",
	"ipv4":     "Field '{field}' must be a valid IPv4 address",
	"ipv6":     "Field '{field}' must be a valid IPv6 address",
	"mac":      "Field '{field}' must be a valid MAC address",
}

func getTagValue(tag string) string {
	if tag == "" {
		return ""
	}
	parts := strings.Split(tag, ",")
	if len(parts) > 0 && parts[0] != "" && parts[0] != "-" {
		return parts[0]
	}
	return ""
}

// extractFieldLabel extracts display label from struct field with priority order
// Priority: json -> form -> uri -> query -> header -> label -> field name
func extractFieldLabel(field reflect.StructField) string {
	// 1. json tag - most commonly used for APIs
	if jsonTag := getTagValue(field.Tag.Get("json")); jsonTag != "" {
		return jsonTag
	}
	// 2. form tag
	if formTag := getTagValue(field.Tag.Get("form")); formTag != "" {
		return formTag
	}
	// 3. uri tag
	if uriTag := getTagValue(field.Tag.Get("uri")); uriTag != "" {
		return uriTag
	}
	// 4. query tag
	if queryTag := getTagValue(field.Tag.Get("query")); queryTag != "" {
		return queryTag
	}
	// 5. header tag
	if headerTag := getTagValue(field.Tag.Get("header")); headerTag != "" {
		return headerTag
	}
	// 6. label tag
	if labelTag := field.Tag.Get("label"); labelTag != "" {
		return labelTag
	}
	// 7. field name as fallback
	return field.Name
}

func parseFieldPath(namespace string, structValue reflect.Value) string {
	if namespace == "" || !structValue.IsValid() {
		return ""
	}

	parts := strings.Split(namespace, ".")
	if len(parts) <= 1 {
		return ""
	}

	currentType := unwrapType(structValue.Type())
	if currentType.Name() != parts[0] {
		return ""
	}

	var result []string
	for _, part := range parts[1:] {
		if currentType.Kind() != reflect.Struct {
			break
		}

		fieldName, indexPart := splitFieldAndIndex(part)

		field, found := currentType.FieldByName(fieldName)
		if !found {
			break
		}

		result = append(result, extractFieldLabel(field)+indexPart)

		currentType = unwrapType(field.Type)
	}

	return strings.Join(result, ".")
}

func unwrapType(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}

	for {
		switch t.Kind() {
		case reflect.Ptr, reflect.Slice, reflect.Array:
			elem := t.Elem()
			if elem == nil {
				return t
			}
			t = elem
		default:
			return t
		}
	}
}

func splitFieldAndIndex(part string) (fieldName, indexPart string) {
	if idx := strings.Index(part, "["); idx != -1 {
		return part[:idx], part[idx:]
	}
	return part, ""
}

func getValidStructValue(structData ...interface{}) reflect.Value {
	if len(structData) == 0 || structData[0] == nil {
		return reflect.Value{}
	}

	value := reflect.ValueOf(structData[0])
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	return value
}

func getErrorMessage(tag string) string {
	if message := ValidationMessages[tag]; message != "" {
		return message
	}
	return "Field '{field}' validation failed (rule: " + tag + ")"
}

func getFieldName(fieldErr validator.FieldError, structValue reflect.Value) string {
	if !structValue.IsValid() {
		return fieldErr.Field()
	}

	if fieldName := parseFieldPath(fieldErr.Namespace(), structValue); fieldName != "" {
		return fieldName
	}

	return fieldErr.Field()
}

// FmtValidationErrors formats validation errors to English error messages
// Supports hierarchical and array field names when struct is provided
func FmtValidationErrors(err error, structData ...interface{}) string {
	if err == nil {
		return ""
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	structValue := getValidStructValue(structData...)
	var messages []string

	for _, fieldErr := range validationErrs {
		message := getErrorMessage(fieldErr.Tag())
		fieldName := getFieldName(fieldErr, structValue)

		message = strings.ReplaceAll(message, "{field}", fieldName)
		message = strings.ReplaceAll(message, "{param}", fieldErr.Param())

		messages = append(messages, message)
	}

	return strings.Join(messages, "; ")
}
