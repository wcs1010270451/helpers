package _struct

import "reflect"

func Copy(src, dest interface{}) {
	vSrc := reflect.ValueOf(src).Elem()
	vDest := reflect.ValueOf(dest).Elem()

	for i := 0; i < vSrc.NumField(); i++ {
		fieldName := vSrc.Type().Field(i).Name
		vDestField := vDest.FieldByName(fieldName) //获取dest中相同名字的字段
		if vDestField.IsValid() && vDestField.CanSet() {
			vDestField.Set(vSrc.Field(i)) //如果这个字段有效并且能被设置，则进行字段值复制
		}
	}
}
