package migrate

import (
	"fmt"
	mysqlDb "github.com/goodluckxu-go/migrate/schema/mysql/tb"
	pgsqlDb "github.com/goodluckxu-go/migrate/schema/pgsql/tb"
	"go/token"
	"reflect"
)

type arg struct {
	Type     string
	ArgTypes []string
	Pos      token.Pos
}

var (
	internetFuncValid = map[string]map[string]map[string]map[string]arg{}
)

func init() {
	// 注册验证类
	internetFuncSet("mysql", "*tb.CreateTable", reflect.TypeOf(&mysqlDb.CreateTable{}), true)
	internetFuncSet("mysql", "*tb.UpdateTable", reflect.TypeOf(&mysqlDb.UpdateTable{}), true)
	internetFuncSet("pgsql", "*tb.CreateTable", reflect.TypeOf(&pgsqlDb.CreateTable{}), true)
	internetFuncSet("pgsql", "*tb.UpdateTable", reflect.TypeOf(&pgsqlDb.UpdateTable{}), true)
	internetFuncSet("pgsql", "*tb.IndexesColumn", reflect.TypeOf(&pgsqlDb.IndexesColumn{}), true)
	//for k, v := range internetFuncValid["pgsql"]["*tb.CreateTable"] {
	//	fmt.Println(k, v)
	//}
}

func internetFuncSet(sqlType, funcName string, valType reflect.Type, init bool) {
	if internetFuncValid[sqlType] == nil {
		internetFuncValid[sqlType] = map[string]map[string]map[string]arg{}
	}
	if internetFuncValid[sqlType][funcName] == nil {
		internetFuncValid[sqlType][funcName] = map[string]map[string]arg{}
	}
	fType := "init"
	if !init {
		fType = valType.String()
	}
	if internetFuncValid[sqlType][funcName][fType] == nil {
		internetFuncValid[sqlType][funcName][fType] = map[string]arg{}
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
			internetFuncValid[sqlType][funcName][fType][valType.Method(i).Name] = arg{
				Type:     out.String(),
				ArgTypes: argTypes,
			}
			if internetFuncValid[sqlType][funcName][out.String()] == nil {
				internetFuncSet(sqlType, funcName, out, false)
			}
		} else {
			internetFuncValid[sqlType][funcName][fType][valType.Method(i).Name] = arg{
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
		internetFuncValid[sqlType][funcName][fType][valElem.Field(i).Name] = arg{
			Type: out.String(),
		}
		if internetFuncValid[sqlType][funcName][out.String()] == nil {
			internetFuncSet(sqlType, funcName, out, false)
		}
	}
}

// 验证内部方法字段
func validInternetFunc(sqlType, funcName string, funcNameList []arg) (pos token.Pos, err error) {
	if internetFuncValid[sqlType] == nil {
		pos = -1
		err = fmt.Errorf("database type %v does not exist", sqlType)
		return
	}
	funcMap := internetFuncValid[sqlType][funcName]
	n := len(funcNameList)
	if funcMap == nil || n == 0 {
		return
	}
	initMap := funcMap["init"]
	isStop := false
	for i := 0; i < n; i++ {
		input := funcNameList[i]
		if isStop {
			pos = input.Pos
			err = fmt.Errorf("method %v not exist", input.Type)
			return
		}
		funcType := initMap[input.Type]
		if funcType.Type == "" {
			pos = input.Pos
			err = fmt.Errorf("method %v not exist", input.Type)
			return
		} else {
			// 验证参数以及类型
			if len(funcType.ArgTypes) == 0 {
				if len(funcNameList[i].ArgTypes) != 0 {
					pos = input.Pos
					err = fmt.Errorf("the number of parameters for method %v should be %v, but there are actually %v",
						input.Type, len(funcType.ArgTypes), len(funcNameList[i].ArgTypes))
					return
				}
			} else {
				total := len(funcType.ArgTypes)
				for j := 0; j < total-1; j++ {
					if len(input.ArgTypes) <= j {
						pos = input.Pos
						err = fmt.Errorf("the %vst parameter of method %v does not exist",
							j+1, input.Type)
						return
					}
					if funcType.ArgTypes[j] != "interface {}" && !isEqType(funcType.ArgTypes[j], input.ArgTypes[j]) {
						pos = input.Pos
						err = fmt.Errorf("the %vst parameter type of method %v should be %v, not %v",
							j+1, input.Type, funcType.ArgTypes[j], input.ArgTypes[j])
						return
					}
				}
				lastArg := funcType.ArgTypes[total-1]
				if lastArg[0:2] == "[]" {
					for j := total - 1; j < len(input.ArgTypes); j++ {
						if !isEqType(lastArg[2:], input.ArgTypes[j]) {
							pos = input.Pos
							err = fmt.Errorf("the %vst parameter type of method %v should be %v, not %v",
								j+1, input.Type, lastArg[2:], input.ArgTypes[j])
							return
						}
					}
				} else if len(input.ArgTypes) != total {
					pos = input.Pos
					err = fmt.Errorf("the number of parameters for method %v should be %v, but there are actually %v",
						input.Type, total, len(input.ArgTypes))
					return
				} else if lastArg != "interface {}" && !isEqType(input.ArgTypes[total-1], lastArg) {
					pos = input.Pos
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

func isEqType(t, inputT string) bool {
	if inputT == "int" && len(t) >= len(inputT) && t[0:len(inputT)] == inputT {
		return true
	}
	return t == inputT
}
