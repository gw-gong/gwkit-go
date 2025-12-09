package validator

import (
	"errors"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type TestData struct {
	Name    string       `json:"user_name" binding:"required"`
	Email   string       `form:"email_addr" binding:"required,email"`
	Token   string       `header:"X-Token" binding:"required"`
	Bio     string       `label:"个人简介" binding:"required"`
	Phone   string       `binding:"required"`
	Age     int          `binding:"min=18,max=100"`
	Profile *TestProfile `json:"profile" binding:"required"`
}

type TestProfile struct {
	NickName string `json:"nick_name" binding:"required"`
	Avatar   string `json:"avatar" binding:"required,url"`
}

func TestValidationMessages(t *testing.T) {
	validate := binding.Validator.Engine().(*validator.Validate)

	testStruct := TestData{
		Name:  "",        // required
		Email: "invalid", // email
		Age:   16,        // min
		Profile: &TestProfile{
			NickName: "",        // required (nested)
			Avatar:   "invalid", // url (nested)
		},
	}

	err := validate.Struct(testStruct)
	if err == nil {
		t.Fatal("Expected validation errors")
	}

	msg1 := FmtValidationErrors(err)
	t.Logf("Basic: %s", msg1)

	msg2 := FmtValidationErrors(err, testStruct)
	t.Logf("With tags: %s", msg2)

	expectedFields := []string{"user_name", "email_addr", "X-Token", "个人简介", "Phone"}
	for _, field := range expectedFields {
		if !strings.Contains(msg2, field) {
			t.Errorf("Missing expected field: %s", field)
		}
	}

	if !strings.Contains(msg2, "profile.nick_name") {
		t.Error("Missing nested field display")
	}
}

func TestTagPriority(t *testing.T) {
	type TagTest struct {
		Field1 string `json:"json_tag" form:"form_tag" binding:"required"`
		Field2 string `form:"form_only" uri:"uri_tag" binding:"required"`
		Field3 string `label:"label_only" binding:"required"`
		Field4 string `binding:"required"`
	}

	validate := binding.Validator.Engine().(*validator.Validate)
	err := validate.Struct(TagTest{})

	result := FmtValidationErrors(err, TagTest{})

	expectedTags := []string{"json_tag", "form_only", "label_only", "Field4"}
	for _, tag := range expectedTags {
		if !strings.Contains(result, tag) {
			t.Errorf("Missing expected tag: %s", tag)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		args []interface{}
	}{
		{"nil error", nil, nil},
		{"regular error", errors.New("test"), nil},
		{"wrong struct type", getValidationError(), []interface{}{123}},
		{"nil struct", getValidationError(), []interface{}{nil}},
		{"empty interface", getValidationError(), []interface{}{interface{}(nil)}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Panic occurred: %v", r)
				}
			}()

			var result string
			if tc.args != nil {
				result = FmtValidationErrors(tc.err, tc.args...)
			} else {
				result = FmtValidationErrors(tc.err)
			}

			t.Logf("Result: %s", result)
		})
	}
}

func TestUnknownTag(t *testing.T) {
	validate := binding.Validator.Engine().(*validator.Validate)

	_ = validate.RegisterValidation("custom", func(fl validator.FieldLevel) bool {
		return false
	})

	type CustomStruct struct {
		Field string `binding:"custom"`
	}

	err := validate.Struct(CustomStruct{Field: "test"})
	if err == nil {
		t.Skip("No validation error generated")
	}

	result := FmtValidationErrors(err)

	if !strings.Contains(result, "validation failed") {
		t.Errorf("Expected default message for unknown tag, got: %s", result)
	}
}

func TestUnwrapTypeSafety(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unwrapType should not panic, but got: %v", r)
		}
	}()

	result := unwrapType(nil)
	if result != nil {
		t.Error("Expected nil for nil input")
	}

	t.Log("unwrapType safety test passed")
}

func getValidationError() error {
	validate := binding.Validator.Engine().(*validator.Validate)
	type SimpleStruct struct {
		Name string `binding:"required"`
	}
	return validate.Struct(SimpleStruct{})
}
