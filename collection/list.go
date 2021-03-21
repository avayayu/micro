package collection

import (
	"fmt"
	"reflect"
)

//CheckSlice 确认参数
func CheckSlice(a interface{}) error {
	typ := reflect.TypeOf(a)
	if typ.Kind() != reflect.Slice {
		return fmt.Errorf("a must be slice,but a is %v", typ.Name())
	}
	return nil
}

func CheckSlicePTR(a interface{}) error {
	typ := reflect.TypeOf(a)

	if typ.Kind() != reflect.Ptr {
		return fmt.Errorf("a must be ptr,but a is %v", typ.Name())
	}

	if typ.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("ptr of a must be slice,but ptr of a is %s", typ.Elem().Name())
	}
	return nil
}

//GetColumFromSlice 从slice结构体中抽取一列作为新的slice
func GetStructColumFromSlice(in interface{}, column string, out interface{}) error {
	inType := reflect.TypeOf(in)
	outType := reflect.TypeOf(out)
	if inType.Kind() != reflect.Slice {
		return fmt.Errorf("in参数必须为slice类型或者slice的引用类型")
	}

	if outType.Kind() != reflect.Ptr {
		return fmt.Errorf("out参数必须为PTR类型")
	}

	if outType.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("out参数必须为slice的PTR类型")
	}

	inValue := reflect.ValueOf(in)
	if inType.Kind() == reflect.Ptr {
		if inType.Elem().Kind() != reflect.Slice {
			return fmt.Errorf("in参数必须为slice类型或者slice的引用类型")
		}

		if inType.Elem().Elem().Kind() != reflect.Struct && inType.Elem().Elem().Elem().Kind() != reflect.Struct {
			return fmt.Errorf("slice的成员必须为结构体或者结构体")
		}
		inValue = inValue.Elem()
	}

	outValue := reflect.ValueOf(out)
	outSlice := reflect.New(reflect.SliceOf(outType.Elem().Elem())).Elem()

	for i := 0; i < inValue.Len(); i++ {
		valueStruct := inValue.Index(i)
		value := valueStruct.FieldByName(column)
		if !value.IsValid() || value.IsZero() {
			return fmt.Errorf("column %s not exists", column)
		}
		outSlice = reflect.Append(outSlice, value)
	}
	outValue.Elem().Set(outSlice)
	return nil
}

//GetDiff 在a中找到属于B但不在a中的元素，并将值返回到out中
// func GetDiff(a interface{}, b interface{}, out interface{}) error {

// 	typA := reflect.TypeOf(a)
// 	typB := reflect.TypeOf(b)

// }
