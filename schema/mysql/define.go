package mysql

import (
	"fmt"
	"github.com/goodluckxu-go/migrate/schema/mysql/tb"
	"reflect"
)

var (
	funcArgValid = map[string]string{
		"Create":            "tb.CreateTable",
		"CreateIfNotExists": "tb.CreateTable",
		"Update":            "tb.UpdateTable",
	}
	internetFuncValid = map[string]map[string]map[string]Arg{}
)

type Arg struct {
	Type     string
	ArgTypes []string
}

func init() {
	internetFuncValid = map[string]map[string]map[string]Arg{}
	internetFuncSet("tb.CreateTable", reflect.TypeOf(&tb.CreateTable{}), true)
	internetFuncSet("tb.UpdateTable", reflect.TypeOf(&tb.UpdateTable{}), true)
}

func internetFuncSet(funcName string, valType reflect.Type, init bool) {
	if internetFuncValid[funcName] == nil {
		internetFuncValid[funcName] = map[string]map[string]Arg{}
	}
	fType := "init"
	if !init {
		fType = valType.String()
	}
	if internetFuncValid[funcName][fType] == nil {
		internetFuncValid[funcName][fType] = map[string]Arg{}
	}
	// 方法
	numMethod := valType.NumMethod()
	for i := 0; i < numMethod; i++ {
		var argTypes []string
		inNum := valType.Method(i).Type.NumIn()
		for j := 1; j < inNum; j++ {
			argTypes = append(argTypes, valType.Method(i).Type.In(j).String())
		}
		if valType.Method(i).Type.NumOut() != 0 {
			out := valType.Method(i).Type.Out(0)
			internetFuncValid[funcName][fType][valType.Method(i).Name] = Arg{
				Type:     out.String(),
				ArgTypes: argTypes,
			}
			if internetFuncValid[funcName][out.String()] == nil {
				internetFuncSet(funcName, out, false)
			}
		} else {
			internetFuncValid[funcName][fType][valType.Method(i).Name] = Arg{
				Type:     "nil",
				ArgTypes: argTypes,
			}
		}
	}
	// 参数
	valElem := valType.Elem()
	numField := valElem.NumField()
	for i := 0; i < numField; i++ {
		out := valElem.Field(i).Type
		internetFuncValid[funcName][fType][valElem.Field(i).Name] = Arg{
			Type: out.String(),
		}
		if internetFuncValid[funcName][out.String()] == nil {
			internetFuncSet(funcName, out, false)
		}
	}
}

// 验证mysql内部方法字段
func validInternetFunc(funcName string, funcNameList []Arg) (err error) {
	funcMap := internetFuncValid[funcName]
	n := len(funcNameList)
	if funcMap == nil || n == 0 {
		return
	}
	initMap := funcMap["init"]
	isStop := false
	for i := 0; i < n; i++ {
		input := funcNameList[i]
		if isStop {
			err = fmt.Errorf("method %v not exist", input.Type)
			return
		}
		funcType := initMap[input.Type]
		if funcType.Type == "" {
			err = fmt.Errorf("method %v not exist", input.Type)
			return
		} else {
			// 验证参数以及类型
			if len(funcType.ArgTypes) == 0 {
				if len(funcNameList[i].ArgTypes) != 0 {
					err = fmt.Errorf("the number of parameters for method %v should be %v, but there are actually %v",
						input.Type, len(funcType.ArgTypes), len(funcNameList[i].ArgTypes))
					return
				}
			} else {
				total := len(funcType.ArgTypes)
				for j := 0; j < total-1; j++ {
					if len(input.ArgTypes) <= j {
						err = fmt.Errorf("the %vst parameter of method %v does not exist",
							j+1, input.Type)
						return
					}
					if funcType.ArgTypes[j] != "interface {}" && funcType.ArgTypes[j] != input.ArgTypes[j] {
						err = fmt.Errorf("the %vst parameter type of method %v should be %v, not %v",
							j+1, input.Type, funcType.ArgTypes[j], input.ArgTypes[j])
						return
					}
				}
				lastArg := funcType.ArgTypes[total-1]
				if lastArg[0:2] == "[]" {
					for j := total - 1; j < len(input.ArgTypes); j++ {
						if input.ArgTypes[j] != lastArg[2:] {
							err = fmt.Errorf("the %vst parameter type of method %v should be %v, not %v",
								j+1, input.Type, lastArg[2:], input.ArgTypes[j])
							return
						}
					}
				} else if len(input.ArgTypes) != total {
					err = fmt.Errorf("the number of parameters for method %v should be %v, but there are actually %v",
						input.Type, total, len(input.ArgTypes))
					return
				} else if lastArg != "interface {}" && lastArg != input.ArgTypes[total-1] {
					err = fmt.Errorf("the %vst parameter type of method %v should be %v, not %v",
						total, input.Type, lastArg, input.ArgTypes[total-1])
					return
				}
			}
			initMap = funcMap[funcType.Type]
			if funcType.Type == "nil" {
				isStop = true
			}
		}
	}
	return
}

type Create func(table tb.CreateTable)

type CreateIfNotExists func(table tb.CreateTable)

type Update func(table tb.UpdateTable)

type Drop uint8

type DropIfExists uint8
