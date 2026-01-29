package validatorutil

import (
	"fmt"
	"reflect"
)

func RequiredStrings(input any) error {
	return validateRequiredValue(reflect.ValueOf(input))
}

func validateRequiredValue(value reflect.Value) error {
	if !value.IsValid() {
		return nil
	}
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Struct:
		valueType := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := valueType.Field(i)
			if field.PkgPath != "" {
				continue
			}
			fieldValue := value.Field(i)
			if field.Type.Kind() == reflect.String && field.Tag.Get("validate") == "required" {
				if fieldValue.String() == "" {
					return fmt.Errorf("field %s is empty", field.Name)
				}
			}
			if err := validateRequiredValue(fieldValue); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			if err := validateRequiredValue(value.Index(i)); err != nil {
				return err
			}
		}
	}
	return nil
}
