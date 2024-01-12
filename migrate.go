package migrate

func ParseSQL(filePath string, funcNameList []string) (sqlMapList map[string][]string, err error) {
	sqlMapList = map[string][]string{}
	var funcList map[string][]tableAst
	funcList, err = ParseGoAst(filePath, funcNameList)
	if err != nil {
		return
	}
	for _, tableList := range funcList {
		for _, table := range tableList {
			switch table.Type {
			case "mysql":
				sqlMapList[table.Type] = append(sqlMapList[table.Type], new(parseMySQL).ParseSQL(table)...)
			}
		}
	}
	return
}
