package serializer

import (
	"context"
	"fmt"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
)

type StringArraySerializer struct{}

func init() {
	schema.RegisterSerializer("string_array", StringArraySerializer{})
}

// Scan implements serializer interface
func (StringArraySerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	fieldValue := reflect.New(field.FieldType)

	if dbValue != nil {
		var bytes []byte
		switch v := dbValue.(type) {
		case []byte:
			bytes = v
		case string:
			bytes = []byte(v)
		default:
			return fmt.Errorf("failed to unmarshal StringArray value: %#v", dbValue)
		}

		array := strings.Split(string(bytes), ",")
		fieldValue.Elem().Set(reflect.ValueOf(array))
	}

	field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
	return
}

// Value implements serializer interface
func (StringArraySerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if _, ok := fieldValue.([]string); !ok {
		return nil, fmt.Errorf("failed to marshal StringArray value: %#v", fieldValue)
	}
	return strings.Join(fieldValue.([]string), ","), nil
}
