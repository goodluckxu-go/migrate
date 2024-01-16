package migrate

// ParseSQL 解析sql语句
// param: filePath string 文件路径
// param: funcNameList []string 需要解析的方法
// result: sqlMapList  第一层为解析的方法，第二层为解析sql类型[mysql等]，第三层为sql语句列表
func ParseSQL(filePath string, funcNameList []string) (sqlMapList map[string]map[string][]string, err error) {
	sqlMapList = map[string]map[string][]string{}
	var funcList map[string][]tableAst
	funcList, err = ParseGoAst(filePath, funcNameList)
	if err != nil {
		return
	}
	for funcName, tableList := range funcList {
		if sqlMapList[funcName] == nil {
			sqlMapList[funcName] = map[string][]string{}
		}
		for _, table := range tableList {
			switch table.Type {
			case "mysql":
				sqlMapList[funcName][table.Type] = append(sqlMapList[funcName][table.Type], new(parseMySQL).ParseSQL(table)...)
			case "pgsql":
				sqlMapList[funcName][table.Type] = append(sqlMapList[funcName][table.Type], new(parsePgSQL).ParseSQL(table)...)
			}
		}
	}
	return
}
