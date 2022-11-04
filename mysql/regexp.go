package mysql

import (
	"reflect"
	"regexp"
	"strings"
)

type Regexp struct {
}

func (r *Regexp) getFunc(content string, funcType string) (string, error) {
	funcTypeValue := reflect.ValueOf(FuncType)
	regString := `type( |\t|\n)+?(\w*?)( |\t)+?struct( |\t|\n)*?\{(?s).*?\}`
	reg := regexp.MustCompile(regString)
	list := reg.FindAllStringSubmatch(content, -1)
	if len(list) == 0 {
		return "", validErr("struct does not exist")
	}
	if len(list) > 1 {
		return "", validErr("struct has %d quantities, but only 1 is required", len(list))
	}
	tableStruct := list[0][2]
	regString = `func( |\t|\n)*?\(\w+( |\t)*\*(\w+)( |\t)*?\)( |\t)*?(\w+)( |\t)*?\(( |\t|\n)*?\)( |\t)*?\{(?s).*?\}( |\t)*\n`
	reg = regexp.MustCompile(regString)
	upString := ""
	for _, funcV := range reg.FindAllStringSubmatch(content, -1) {
		if tableStruct != funcV[3] {
			return "", validErr("struct '%s' is different from '%s'", funcV[3], tableStruct)
		}
		if !funcTypeValue.FieldByName(funcV[6]).IsValid() {
			return "", validErr("func '%s' must be in 'FuncType'", funcV[6])
		}
		if funcV[6] != funcType {
			continue
		}
		upString = funcV[0]
	}
	return upString, nil
}

func (r *Regexp) getSchema(content string) []map[string]interface{} {
	regString := `schema\.(\w+)( |\t)*?\(( |\t|\n)*?\"(\w*?)\"`
	regString += `(( |\t)*?,( |\t\n)*?func( |\t)*?\(.*?\)( |\t)*?\{((?s).*?)\}( |\t|\n)*?\)((\.(?s).+?)*?)\n)*`
	reg := regexp.MustCompile(regString)
	rs := []map[string]interface{}{}
	for _, vMap := range reg.FindAllStringSubmatch(content, -1) {
		rs = append(rs, map[string]interface{}{
			"type":    vMap[1],
			"table":   vMap[4],
			"content": vMap[10],
			"other":   vMap[12],
		})
	}
	return rs
}

func (r *Regexp) getColumn(content string) []map[string][]string {
	regString := `table( |\t)*(\.((?s).+?)\))( |\t)*\n`
	reg := regexp.MustCompile(regString)
	rs := []map[string][]string{}
	for _, vMap := range reg.FindAllStringSubmatch(content, -1) {
		columnString := vMap[2]
		columnInfo := r.getOneColumn(columnString)
		rs = append(rs, columnInfo)
	}
	return rs
}

func (r *Regexp) getOneColumn(content string) map[string][]string {
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\t", "")
	content = regexp.MustCompile(`(\))\.`).ReplaceAllString(content, "$1..")
	content += "."
	regString := `\.( \t|\n)*(\w+)\((.*?)*\)\.`
	reg := regexp.MustCompile(regString)
	rs := map[string][]string{}
	for _, columnMap := range reg.FindAllStringSubmatch(content, -1) {
		columnType := columnMap[2]
		rs["validSort"] = append(rs["validSort"], columnType)
		columnTmp := strings.Split(columnMap[3], ",")
		columnValList := []string{}
		for _, columnVal := range columnTmp {
			columnVal = strings.Trim(columnVal, " |\t")
			if columnVal == "" {
				continue
			}
			columnValList = append(columnValList, columnVal)
		}
		rs[columnType] = columnValList
	}
	return rs
}
