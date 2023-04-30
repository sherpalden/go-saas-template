package httpi

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BuildQueryParams(ctx *gin.Context, dest interface{}) error {
	typOf := reflect.TypeOf(dest)
	if typOf.Kind() != reflect.Pointer {
		return errors.New("destination type must a pointer")
	}
	valOf := reflect.ValueOf(dest)
	for i := 0; i < typOf.Elem().NumField(); i++ {
		queryVal := ctx.Query(typOf.Elem().Field(i).Tag.Get("json"))
		switch typOf.Elem().Field(i).Type.Kind() {
		default:
			return fmt.Errorf("invalid type %T for field %v", typOf.Elem().Field(i).Type.Kind(), typOf.Elem().Field(i).Name)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(queryVal)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetBool(boolVal)
		case reflect.Int:
			intVal, err := strconv.Atoi(queryVal)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetInt(int64(intVal))
		case reflect.Int8:
			intVal, err := strconv.Atoi(queryVal)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetInt(int64(intVal))
		case reflect.Int16:
			intVal, err := strconv.Atoi(queryVal)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetInt(int64(intVal))
		case reflect.Int32:
			intVal, err := strconv.Atoi(queryVal)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetInt(int64(intVal))
		case reflect.Int64:
			intVal, err := strconv.Atoi(queryVal)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetInt(int64(intVal))
		case reflect.Uint:
			uintVal, err := strconv.ParseUint(queryVal, 10, 32)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetUint(uintVal)
		case reflect.Uint8:
			uintVal, err := strconv.ParseUint(queryVal, 10, 8)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetUint(uintVal)
		case reflect.Uint16:
			uintVal, err := strconv.ParseUint(queryVal, 10, 16)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetUint(uintVal)
		case reflect.Uint32:
			uintVal, err := strconv.ParseUint(queryVal, 10, 32)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetUint(uintVal)
		case reflect.Uint64:
			uintVal, err := strconv.ParseUint(queryVal, 10, 64)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetUint(uintVal)
		case reflect.Float32:
			floatVal, err := strconv.ParseFloat(queryVal, 32)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetFloat(floatVal)
		case reflect.Float64:
			floatVal, err := strconv.ParseFloat(queryVal, 64)
			if err != nil {
				return err
			}
			valOf.Elem().Field(i).SetFloat(floatVal)
		case reflect.String:
			valOf.Elem().Field(i).SetString(queryVal)
		}
	}
	return nil
}
