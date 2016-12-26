package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func GetRequest(r map[string]interface{}, obj interface{}) error {
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i)
		tField := t.Field(i)

		tag := tField.Tag.Get("validate")
		if tag == "" {
			continue
		}

		switch vField.Kind() {
		case reflect.String:
			if tag == "required" {
				name := strings.ToLower(tField.Name)
				if f, ok := r[name]; !ok || reflect.ValueOf(f).Kind() != reflect.String || f == "" {
					return errors.New(fmt.Sprintf("invalid %v", name))
				}
				vField.SetString(r[name].(string))
			}
			if tag == "neglect" {
				name := strings.ToLower(tField.Name)
				if f, ok := r[name]; ok && reflect.ValueOf(f).Kind() == reflect.String && f != "" {
					vField.SetString(r[name].(string))
				}
			}

		case reflect.Int64:
			if tag == "required" {
				name := strings.ToLower(tField.Name)
				if f, ok := r[name]; !ok || reflect.ValueOf(f).Kind() != reflect.Float64 || f.(float64) < 0 {
					return errors.New(fmt.Sprintf("invalid %v", name))
				}
				vField.SetInt(int64(r[name].(float64)))
			}
			if tag == "neglect" {
				name := strings.ToLower(tField.Name)
				if f, ok := r[name]; ok && reflect.ValueOf(f).Kind() == reflect.Float64 && f.(float64) >= 0 {
					vField.SetInt(int64(r[name].(float64)))
				}
			}

		case reflect.Bool:
			if tag == "required" {
				name := strings.ToLower(tField.Name)
				if f, ok := r[name]; !ok || reflect.ValueOf(f).Kind() != reflect.Bool {
					return errors.New(fmt.Sprintf("invalid %v", name))
				}
			}
			if tag == "neglect" {
				name := strings.ToLower(tField.Name)
				if f, ok := r[name]; ok && reflect.ValueOf(f).Kind() == reflect.Bool {
					vField.SetBool(f.(bool))
				}
			}

		default:
			return fmt.Errorf("Unsupported kind '%s'", vField.Kind())
		}
	}
	return nil
}
