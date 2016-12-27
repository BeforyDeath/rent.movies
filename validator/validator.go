package validator

import (
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

		if tag == "required" {
			res, err := find(tField.Name, vField.Kind(), r)
			if err != nil {
				return fmt.Errorf("%v (required)", err)
			}

			err = write(vField, res)
			if err != nil {
				return err
			}
		}

		if tag == "neglect" {
			res, err := find(tField.Name, vField.Kind(), r)
			if err == nil {
				err := write(vField, res)
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}

func find(nField string, tField reflect.Kind, r map[string]interface{}) (interface{}, error) {
	name := strings.ToLower(nField)
	if nObjField, ok := r[name]; ok {

		tObjField := reflect.ValueOf(nObjField).Kind()
		if tObjField == reflect.Float64 && tField == reflect.Int64 {
			tField = reflect.Float64
		}

		if tObjField == tField {
			return nObjField, nil
		}

		return nil, fmt.Errorf("Invalid type field '%v', expected %v", name, tField)
	}
	return nil, fmt.Errorf("Field '%v' not found", name)
}

func write(f reflect.Value, r interface{}) error {
	switch f.Kind() {
	case reflect.String:
		f.SetString(r.(string))
	case reflect.Int64:
		f.SetInt(int64(r.(float64)))
	case reflect.Bool:
		f.SetBool(r.(bool))
	default:
		return fmt.Errorf("Unsupported kind '%s'", f.Kind())
	}
	return nil
}
