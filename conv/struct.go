package conv

import (
	"reflect"
)

// 使用reflect实现struct同名字段转换，支持Struct、Ptr迭代
func StructToStruct(src, dst interface{}) {
	s := reflect.ValueOf(src).Elem()
	d := reflect.ValueOf(dst).Elem()
	copy(s, d)
}

func copy(src, dst reflect.Value) {
	// 更多考虑的场景是struct字段压缩，所以使用dval.NumField()遍历
	for i := 0; i < dst.NumField(); i++ {
		dValue := dst.Field(i)
		name := dst.Type().Field(i).Name

		sValue := src.FieldByName(name)
		if sValue.IsValid() == false {
			continue
		}

		if dValue.Type() != sValue.Type() {
			switch dValue.Type().Kind() {
			case reflect.Struct:
				switch sValue.Kind() {
				case reflect.Struct:
					copy(sValue, dValue)
				case reflect.Ptr:
					copy(sValue.Elem(), dValue)
				}
			case reflect.Ptr:
				if dValue.IsZero() {
					dValue.Set(reflect.New(dValue.Type().Elem()))
				}
				switch sValue.Kind() {
				case reflect.Struct:
					copy(sValue, dValue.Elem())
				case reflect.Ptr:
					copy(sValue.Elem(), dValue.Elem())
				}
			}

			continue
		} else {
			dValue.Set(sValue)
		}
	}
}
