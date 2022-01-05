package tool

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"strings"
)

func StructName(oi interface{}) string {
	typ := typeElem(reflect.TypeOf(oi))
	name := typ.Name()
	if name == "" {
		return ""
	}

	if packageDotInd := strings.LastIndex(name, "."); packageDotInd != -1 {
		return name[packageDotInd:]
	}
	return name
}
func StructNewAddr(oi interface{}) interface{} {
	return valueInterface(structNewAddr(typeElem(reflect.TypeOf(oi))))
}
func IsStruct(oi interface{}) (bool) {
	typ := typeElem(reflect.TypeOf(oi))
	return isTypeStruct(typ)
}
func StructTypeNewAddr(typ reflect.Type) interface{} {
	return valueInterface(reflect.New(typeElem(typ)))
}
func IsStructSame(o1 interface{}, o2 interface{}) bool {
	o1Type := reflect.TypeOf(o1)
	o2Type := reflect.TypeOf(o2)
	return isTypeSame(typeElem(o1Type), typeElem(o2Type))
}

func StructDeepCopy(oi interface{}, dstoi interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(oi); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dstoi)
}

var DeepCopy = StructDeepCopy

func FieldNameSlice(oi interface{}) []string {
	fields := []string{}
	o := valueElem(reflect.ValueOf(&oi))
	for i := 0; i < o.NumField(); i++ {
		fieldInfo := o.Type().Field(i)
		fields = append(fields, fieldInfo.Name)
	}
	return fields
}
func Field(oi interface{}, fieldName string) interface{} {
	return valueInterface(getField(reflect.ValueOf(oi), fieldName))
}
func Field2(oi interface{}, fieldName1 string, fieldName2 string) interface{} {
	return Field(Field(oi, fieldName1), fieldName2)
}
func AsField(oi interface{}, typ interface{}, fieldName string) interface{} {
	return valueInterface(asField(
		reflect.ValueOf(oi),
		reflect.TypeOf(typ),
		fieldName,
	))
}
func AsField2(oi interface{}, typ1 interface{}, fieldName1 string, typ2 interface{}, fieldName2 string) interface{} {
	oi = AsField(oi, typ1, fieldName1)
	if IsStructSame(oi, typ2) {
		return oi
	} else {
		return AsField(oi, typ2, fieldName2)
	}
}

func FieldTypeMap(oi interface{}, flattenSkip func(interface{}) bool) map[string]reflect.Type {
	return fieldTypeMap(valueElem(reflect.ValueOf(oi)), func(value reflect.Value) bool {
		return flattenSkip(valueInterface(value))
	})
}
func FieldMap(oi interface{}, flattenSkip func(interface{}) bool, tagNameSlice ...string) (map[string]interface{}, map[string]reflect.Type, map[string]*[]string) {
	return fieldMap(valueElem(reflect.ValueOf(oi)), func(value reflect.Value) bool {
		return flattenSkip(valueInterface(value))
	}, tagNameSlice...)
}

func IsFieldZero(oi interface{}, fieldName string) bool {
	return isFieldZero(reflect.ValueOf(oi), fieldName)
}

func SetField(oi interface{}, fieldName string, field interface{}) {
	setField(reflect.ValueOf(oi), fieldName, reflect.ValueOf(field))
}

func SetFieldZero(oi interface{}, fieldName string) {
	oValue := valueElem(reflect.ValueOf(oi))
	fieldValue := oValue.FieldByName(fieldName)
	valueSetZero(fieldValue)
}

func fieldFillZero(value reflect.Value) {
	value = valueElem(value)
	typ := value.Type()
	if !isTypeStruct(typ) {
		return
	}

	fieldNum := value.NumField()
	for i := 0; i < fieldNum; i++ {
		fieldType := typ.Field(i).Type
		field := value.Field(i)
		//fmt.Println(typ.Field(i).Name)

		if isTypeStruct(fieldType) {
			//fmt.Println(fieldType.Name())
			field.Set(typeZeroValue(fieldType))
			if !isTypeTime(fieldType) {
				fieldFillZero(field)
			}
		}

		if isPtr(fieldType) {
			fieldTypeElem := fieldType.Elem()
			//fmt.Println(fieldTypeElem.Name())
			if isTypeStruct(fieldTypeElem) {
				field.Set(typeZeroValueAddr(fieldTypeElem))
				if !isTypeTime(fieldTypeElem) {
					fieldFillZero(field)
				}
			}
		}
	}
}
func FieldFillZero(oi interface{}) {
	fieldFillZero(reflect.ValueOf(oi))
}

func StructToMap(oi interface{}) map[string]interface{} {
	structMap := map[string]interface{}{}

	o := valueElem(reflect.ValueOf(&oi))
	for i := 0; i < o.NumField(); i++ {
		fieldInfo := o.Type().Field(i)
		fieldName := fieldInfo.Name
		structMap[fieldName] = valueInterface(getField(o, fieldName))
	}
	return structMap
}
