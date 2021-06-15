package session

import (
	"github.com/ZhaoxingZhang/go-gees/geecommon/log"
	"reflect"
)

// Hooks constants
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// CallMethod calls the registered hooks
// 每一个钩子的入参类型均是 *Session
// s.RefTable().Model 或 value 即当前会话正在操作的对象，
// 使用 MethodByName 方法反射得到 **对象** 的方法
func (s *Session) CallMethod(method string, value interface{}) {
	hookMethod := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		hookMethod = reflect.ValueOf(value).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if hookMethod.IsValid() { // 判断hook方法存在
		if v := hookMethod.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}
