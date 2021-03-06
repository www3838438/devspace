package configutil

import (
	"reflect"
	"unsafe"
)

// Merge deeply merges two objects
func Merge(object interface{}, overwriteObject interface{}) {
	objectPointerUnsafe := unsafe.Pointer(&object)
	overwriteObjectPointerUnsafe := unsafe.Pointer(&overwriteObject)

	merge(object, overwriteObject, objectPointerUnsafe, overwriteObjectPointerUnsafe)
}

func merge(objectPointer interface{}, overwriteObjectPointer interface{}, objectPointerUnsafe unsafe.Pointer, overwriteObjectPointerUnsafe unsafe.Pointer) {
	overwriteObjectRef := reflect.ValueOf(overwriteObjectPointer)

	if !overwriteObjectRef.IsNil() {
		if overwriteObjectRef.Kind() == reflect.Ptr {
			overwriteObjectRef = overwriteObjectRef.Elem()
		}
		overwriteObject := overwriteObjectRef.Interface()
		overwriteObjectType := reflect.TypeOf(overwriteObject)
		overwriteObjectKind := overwriteObjectType.Kind()
		objectPointerRef := reflect.ValueOf(objectPointer)
		var objectRef reflect.Value

		if !objectPointerRef.IsNil() {
			objectRef = reflect.ValueOf(objectPointer).Elem()
		}

		switch overwriteObjectKind {
		case reflect.Slice:
			if objectPointerRef.IsNil() {
				objectRef.Set(reflect.New(overwriteObjectType))
			}

			for i := 0; i < overwriteObjectRef.Len(); i++ {
				overwriteValue := overwriteObjectRef.Index(i)

				objectRef.Set(reflect.Append(objectRef, overwriteValue))
			}
		case reflect.Map:
			if objectPointerRef.IsNil() {
				objectRef.Set(overwriteObjectRef)
			} else {
				genericPointerType := reflect.TypeOf(overwriteObject)

				for _, keyRef := range overwriteObjectRef.MapKeys() {
					key := keyRef.Interface()
					overwriteValue := getMapValue(overwriteObject, key, genericPointerType)
					valuePointerRef := objectRef.MapIndex(keyRef)

					if isZero(valuePointerRef) == false {
						valuePointer := valuePointerRef.Interface()

						merge(valuePointer, overwriteValue, unsafe.Pointer(&valuePointer), unsafe.Pointer(&overwriteValue))
					} else {
						keyRef := reflect.ValueOf(key)
						overwriteValueRef := reflect.ValueOf(overwriteValue)

						objectRef.SetMapIndex(keyRef, overwriteValueRef)
					}
				}
			}
		case reflect.Struct:
			for i := 0; i < overwriteObjectRef.NumField(); i++ {
				//fieldName := objectRef.Type().Field(i).Name
				overwriteValueRef := overwriteObjectRef.Field(i)
				overwriteValuePointerRef := reflect.ValueOf(overwriteValueRef.Interface())

				if !overwriteValuePointerRef.IsNil() {
					overwriteValue := overwriteValueRef.Interface()
					valuePointerRef := objectRef.Field(i)

					if valuePointerRef.IsNil() {
						objectRef.Field(i).Set(reflect.ValueOf(overwriteValue))
					} else {
						valuePointer := objectRef.Field(i).Interface()

						merge(valuePointer, overwriteValue, unsafe.Pointer(&valuePointer), unsafe.Pointer(&overwriteValue))
					}
				}
			}
		default:
			*(*unsafe.Pointer)(objectPointerUnsafe) = overwriteObjectPointerUnsafe
		}
	}
}
