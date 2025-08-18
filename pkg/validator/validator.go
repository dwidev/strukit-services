package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"strukit-services/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func Run() *AppValidator {
	v := validator.New()

	return &AppValidator{
		Validate: v,
	}
}

type AppValidator struct {
	*validator.Validate
}

func (v *AppValidator) Valid(s any) []string {
	results := make([]string, 0)

	filterError := func(e validator.FieldError) {
		fieldPath := strings.Split(e.StructNamespace(), ".")
		typ := reflect.TypeOf(s).Elem()

		var field reflect.StructField
		var jsonTag string

		for _, part := range fieldPath {
			var found bool
			field, found = typ.FieldByName(part)
			if !found {
				continue
			}

			jsonTag = field.Tag.Get("json")
			if len(jsonTag) == 0 {
				panic(fmt.Sprintf("%s doesn't have json tag", e.StructNamespace()))
			}

			typ = field.Type
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
		}

		validate := field.Tag.Get("validate")
		logger.Log.Error("jsonPath", jsonTag)

		switch e.Tag() {
		case "required":
			msg := fmt.Sprintf("%s is required", jsonTag)
			results = append(results, msg)
			return
		case "max":
			maxVal := getValueFromTag(validate, "max=")
			msg := fmt.Sprintf("%s is max %v", jsonTag, maxVal)
			results = append(results, msg)
			return
		case "min":
			minVal := getValueFromTag(validate, "min=")
			msg := fmt.Sprintf("%s is min %v", jsonTag, minVal)
			results = append(results, msg)
			return
		case "gtfield":
			results = append(results, fmt.Sprintf("%s greater than %s", e.Field(), e.Param()))
		default:
			results = append(results, fmt.Sprintf("Validation error %s", e))
		}

	}

	if err := v.Struct(s); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, ev := range ve { // ev is error validation
				filterError(ev)
			}
		}
	}

	return results
}

func getValueFromTag(tag string, prefix string) string {
	tagValues := strings.SplitSeq(tag, ",")

	for value := range tagValues {
		if after, ok := strings.CutPrefix(value, prefix); ok {
			maxValueStr := after
			return maxValueStr
		}
	}

	return ""
}
