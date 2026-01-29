package httphdlrs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func decodeJsonRequest[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, ErrMsgInvalidRequestBody)
		var emptyT T
		return emptyT, false
	}
	return req, true
}

func validateRequiredStrings(input any) error {
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
