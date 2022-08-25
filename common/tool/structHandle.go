package tool

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func FindSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func IsEmptyValue(args ...interface{}) bool {
	v := args[0].(reflect.Value)
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func Struct2Map(obj interface{}, fun func(...interface{}) bool) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if fun(field, v, i) {
			data[field.Tag.Get("db")] = v.Field(i).Interface()
		}
	}
	return data
}

func Struct2MapPre(obj interface{}) map[string]interface{} {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data[typeOfType.Field(i).Name] = field.Interface()
	}
	return data
}

func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

func BuildFields(obj interface{}, fun func(...interface{}) bool) []string {
	t := reflect.TypeOf(obj)
	//判断结构体类型是否为指针类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(obj)
	//判断结构体值是否为指针类型
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	sqls := []string{}
	//循环结构体所有字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		b := field.Tag.Get("db")
		if b == "create_time" || b == "update_time" || b == "createTime" || b == "updateTime" {
			break
		}
		//判断是否可以拼接sql
		if ! fun(v.Field(i), field) {
			sqls = append(sqls, fmt.Sprintf("`%s`=%#v", field.Tag.Get("db"), v.Field(i).Interface()))
		}
	}
	return sqls
}

func Map2Struct(map_ map[string]interface{}, stu interface{}) error {
	pd_tmp,err := json.Marshal(map_)
	if err == nil {
		err = json.Unmarshal(pd_tmp,&stu)
	}
	return err
}