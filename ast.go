package migrate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// ParseGoAst 解析go语法树
func ParseGoAst(filePath string, funcNameList []string) (funcTableList map[string][]tableAst, err error) {
	if len(funcNameList) == 0 {
		return
	}
	funcNameMap := map[string]bool{}
	for _, v := range funcNameList {
		funcNameMap[v] = true
	}
	funcTableList = map[string][]tableAst{}
	fSet := token.NewFileSet()
	filePath, err = filepath.Abs(filePath)
	if err != nil {
		return
	}
	var f *ast.File
	f, err = parser.ParseFile(fSet, filePath, nil, 0)
	if err != nil {
		return
	}
	funcMap, importMap, err := getFuncMap(fSet, f)
	if err != nil {
		return
	}
	var tableList []tableAst
	for k, v := range funcMap {
		if !funcNameMap[k] {
			continue
		}
		tableList, err = getTableList(fSet, v, importMap)
		if err != nil {
			return
		}
		funcTableList[k] = tableList
	}
	return
}

// 获取各个方法的ast内容
func getFuncMap(fSet *token.FileSet, f *ast.File) (funcMap map[string]*ast.FuncDecl, importMap map[string]string, err error) {
	funcMap = map[string]*ast.FuncDecl{}
	importMap = map[string]string{}
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			funcMap[funcDecl.Name.Name] = funcDecl
		}
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok.String() == "import" {
			for _, spec := range genDecl.Specs {
				if importSpec, ok := spec.(*ast.ImportSpec); ok {
					importPath := strings.Trim(importSpec.Path.Value, "\"|`")
					if !strings.Contains(importPath, modName) {
						err = fmt.Errorf("%v: package %v is not required", fSet.Position(importSpec.Path.ValuePos), importPath)
						return
					}
					importPathList := strings.Split(importPath, "/")
					importName := importPathList[len(importPathList)-1]
					if importSpec.Name != nil {
						importName = importSpec.Name.Name
					}
					importMap[importName] = importPathList[len(importPathList)-1]
				}
			}
		}
	}
	return
}

// 根据indexExpr获取方法名称
func getFunNameByIndexExpr(fSet *token.FileSet, indexExpr *ast.IndexExpr, importMap map[string]string) (table tableAst, err error) {
	var ok bool
	var xSelectExpr *ast.SelectorExpr
	if xSelectExpr, ok = indexExpr.X.(*ast.SelectorExpr); !ok {
		err = fmt.Errorf("%v: the code does not meet the definition, parsing failed",
			fSet.Position(indexExpr.X.Pos()))
		return
	}
	var xIdent *ast.Ident
	if xIdent, ok = xSelectExpr.X.(*ast.Ident); !ok {
		err = fmt.Errorf("%v: the code does not meet the definition, parsing failed",
			fSet.Position(xSelectExpr.X.Pos()))
		return
	}
	importVal := xIdent.Name
	if importMap[xIdent.Name] != "" {
		importVal = importMap[xIdent.Name]
	}
	if importVal != "schema" {
		err = fmt.Errorf("%v: the code does not meet the definition, parsing failed",
			fSet.Position(xIdent.Pos()))
		return
	}
	if schemaFuncValid[xSelectExpr.Sel.Name] == 0 {
		err = fmt.Errorf("%v: schema does not have method %v",
			fSet.Position(xIdent.Pos()), xSelectExpr.Sel.Name)
		return
	}
	table.Func = xSelectExpr.Sel.Name
	// 获取values泛型
	var indexSelectExpr *ast.SelectorExpr
	if indexSelectExpr, ok = indexExpr.Index.(*ast.SelectorExpr); !ok {
		err = fmt.Errorf("%v: the code does not meet the definition, parsing failed",
			fSet.Position(indexExpr.Index.Pos()))
		return
	}
	if xIdent, ok = indexSelectExpr.X.(*ast.Ident); !ok {
		err = fmt.Errorf("%v: the code does not meet the definition, parsing failed",
			fSet.Position(indexSelectExpr.X.Pos()))
		return
	}
	tableType := xIdent.Name
	if importMap[xIdent.Name] != "" {
		tableType = importMap[xIdent.Name]
	}
	table.Type = tableType
	table.Active = indexSelectExpr.Sel.Name
	return
}

// 获取字段
func getColumn(fSet *token.FileSet, col ast.Expr, internetCallName string) (column columnAst, err error) {
	column.InternetFunc = map[string][]argAst{}
	var ok bool
	switch val := col.(type) {
	case *ast.CallExpr:
		var funcExpr *ast.SelectorExpr
		if funcExpr, ok = val.Fun.(*ast.SelectorExpr); !ok {
			err = fmt.Errorf("%v: format error",
				fSet.Position(val.Pos()))
			return
		}
		args := make([]argAst, 0)
		for _, arg := range val.Args {
			switch v := arg.(type) {
			case *ast.Ident:
				if v.Name == "nil" {
					args = append(args, argAst{
						Type: "nil",
					})
				}
			case *ast.BasicLit:
				typeString := v.Kind.String()
				var newVal string
				if typeString == "STRING" {
					newVal = strings.Trim(v.Value, "\"|`")
				} else {
					newVal = v.Value
				}
				args = append(args, argAst{
					Val:  newVal,
					Type: strings.ToLower(typeString),
				})
			}
		}
		column.InternetFunc[funcExpr.Sel.Name] = args
		var allCol columnAst
		allCol, err = getColumn(fSet, funcExpr.X, internetCallName)
		if err != nil {
			return
		}
		for k, v := range allCol.InternetFunc {
			column.InternetFunc[k] = v
		}
		column.LianFuncSort = append(column.LianFuncSort, allCol.LianFuncSort...)
		column.LianFuncSort = append(column.LianFuncSort, funcExpr.Sel.Name)
	case *ast.SelectorExpr:
		// 内部方法的变量
		column.InternetFunc[val.Sel.Name] = make([]argAst, 0)
		column.LianFuncSort = append(column.LianFuncSort, val.Sel.Name)
	case *ast.Ident:
		if val.Obj == nil {
			err = fmt.Errorf("%v: format error",
				fSet.Position(val.Pos()))
			return
		}
		var field *ast.Field
		if field, ok = val.Obj.Decl.(*ast.Field); !ok {
			err = fmt.Errorf("%v: format error",
				fSet.Position(val.Pos()))
			return
		}
		selfCallName := getInternetCallName(field)
		if selfCallName != internetCallName {
			err = fmt.Errorf("%v: internal call method should be %v instead of %v",
				fSet.Position(field.Pos()), internetCallName, selfCallName)
			return
		}
	}
	return
}

// 获取内部调用方法名称
func getInternetCallName(file *ast.Field) (internetCallName string) {
	internetCallName = file.Names[0].Name
	var ok bool
	// 获取类型
	var typeExpr *ast.SelectorExpr
	if typeExpr, ok = file.Type.(*ast.SelectorExpr); !ok {
		return ""
	}
	var xTypeExpr *ast.Ident
	if xTypeExpr, ok = typeExpr.X.(*ast.Ident); !ok {
		return ""
	}
	internetCallName += " " + xTypeExpr.Name + "." + typeExpr.Sel.Name
	return
}

// 获取方法中所有表
func getTableList(fSet *token.FileSet, f *ast.FuncDecl, importMap map[string]string) (tableList []tableAst, err error) {
	if f.Body == nil {
		err = fmt.Errorf("%v: there is no available method for the file",
			fSet.Position(f.Pos()))
		return
	}
	var ok bool
	var tableValue tableAst
	tableFuncMap := map[string]tableAst{}
	useFuncMap := map[string]token.Pos{}
	for _, stmt := range f.Body.List {
		switch val := stmt.(type) {
		case *ast.DeclStmt:
			// 获取var定义的方法
			var genDecl *ast.GenDecl
			if genDecl, ok = val.Decl.(*ast.GenDecl); !ok {
				continue
			}
			if genDecl.Tok.String() != "var" || len(genDecl.Specs) == 0 {
				continue
			}
			for _, spec := range genDecl.Specs {
				var valueSpec *ast.ValueSpec
				if valueSpec, ok = spec.(*ast.ValueSpec); !ok {
					continue
				}
				if valueSpec.Names == nil || len(valueSpec.Names) != len(valueSpec.Values) {
					continue
				}
				for index, name := range valueSpec.Names {
					// 获取values内容
					var indexExpr *ast.IndexExpr
					if indexExpr, ok = valueSpec.Values[index].(*ast.IndexExpr); !ok {
						continue
					}
					tableValue, err = getFunNameByIndexExpr(fSet, indexExpr, importMap)
					if err != nil {
						return
					}
					tableFuncMap[name.Name] = tableValue
					useFuncMap[name.Name] = name.Pos()
				}
			}
		case *ast.AssignStmt:
			// 获取:=定义的方法
			if val.Tok.String() != ":=" {
				continue
			}
			if len(val.Lhs) != len(val.Rhs) {
				continue
			}
			for index, name := range val.Lhs {
				var indexExpr *ast.IndexExpr
				if indexExpr, ok = val.Rhs[index].(*ast.IndexExpr); !ok {
					continue
				}
				var ident *ast.Ident
				if ident, ok = name.(*ast.Ident); !ok {
					continue
				}
				tableValue, err = getFunNameByIndexExpr(fSet, indexExpr, importMap)
				if err != nil {
					return
				}
				tableFuncMap[ident.Name] = tableValue
				useFuncMap[ident.Name] = ident.Pos()
			}
		case *ast.ExprStmt:
			// 方法
			var xCallExpr *ast.CallExpr
			if xCallExpr, ok = val.X.(*ast.CallExpr); !ok {
				continue
			}
			switch xFunVal := xCallExpr.Fun.(type) {
			case *ast.Ident:
				// 获取引用方法
				tableValue = tableFuncMap[xFunVal.Name]
				if tableValue.Type == "" {
					err = fmt.Errorf("%v: method %v not found",
						fSet.Position(xFunVal.Pos()), xFunVal.Name)
					return
				}
				delete(useFuncMap, xFunVal.Name)
			case *ast.IndexExpr:
				tableValue, err = getFunNameByIndexExpr(fSet, xFunVal, importMap)
				if err != nil {
					return
				}
			default:
				continue
			}
			// 获取方法参数
			if schemaFuncValid[tableValue.Func] != len(xCallExpr.Args) {
				err = fmt.Errorf("%v: the number of parameters for method %v should be %v, but there are actually %v",
					fSet.Position(xCallExpr.Fun.Pos()), tableValue.Func, schemaFuncValid[tableValue.Func],
					len(xCallExpr.Args))
				return
			}
			// 验证参数1
			var tableArg *ast.BasicLit
			if tableArg, ok = xCallExpr.Args[0].(*ast.BasicLit); !ok {
				err = fmt.Errorf("%v: the 1st parameter type of method %v should be string, not %T",
					fSet.Position(xCallExpr.Fun.Pos()), tableValue.Func, xCallExpr.Args[0])
				return
			}
			if tableArg.Kind.String() != "STRING" {
				err = fmt.Errorf("%v: the 1st parameter type of method %v should be string, not %v",
					fSet.Position(xCallExpr.Fun.Pos()), tableValue.Func, strings.ToLower(tableArg.Kind.String()))
				return
			}
			tableValue.Name = strings.Trim(tableArg.Value, "\"|`")
			// 如果参数2存在则验证
			if len(xCallExpr.Args) > 1 {
				arg2Type := funcArgValid[tableValue.Type][tableValue.Active]
				arg2Note := "func(" + arg2Type + ")"
				var funcArg *ast.FuncLit
				if funcArg, ok = xCallExpr.Args[1].(*ast.FuncLit); !ok {
					err = fmt.Errorf("%v: the 2st parameter type of method %v should be %v, not %v",
						fSet.Position(xCallExpr.Fun.Pos()), tableValue.Func, arg2Note, getArgType(xCallExpr.Args[1]))
					return
				}
				if funcArg.Type == nil || len(funcArg.Type.Params.List) != 1 {
					err = fmt.Errorf("%v: the 2st parameter type of method %v should be %v, not %v",
						fSet.Position(xCallExpr.Fun.Pos()), tableValue.Func, arg2Note, getArgType(xCallExpr.Args[1]))
					return
				}
				param := funcArg.Type.Params.List[0]
				if len(param.Names) != 1 {
					err = fmt.Errorf("%v: the 2st parameter type of method %v should be %v, not %v",
						fSet.Position(xCallExpr.Fun.Pos()), tableValue.Func, arg2Note, getArgType(xCallExpr.Args[1]))
					return
				}
				// 内部调用方法名称
				internetCallName := getInternetCallName(param)
				// 判断类型是否和传入类型一致
				var xSelectExpr *ast.SelectorExpr
				if xSelectExpr, ok = param.Type.(*ast.SelectorExpr); !ok {
					err = fmt.Errorf("%v: the 2st parameter type of method %v should be %v, not %v",
						fSet.Position(xCallExpr.Fun.Pos()), tableValue.Func, arg2Note, getArgType(xCallExpr.Args[1]))
					return
				}
				var xIdent *ast.Ident
				if xIdent, ok = xSelectExpr.X.(*ast.Ident); !ok {
					err = fmt.Errorf("%v: the 2st parameter type of method %v should be %v, not %v",
						fSet.Position(xCallExpr.Fun.Pos()), tableValue.Func, arg2Note, getArgType(xCallExpr.Args[1]))
					return
				}
				argName := xIdent.Name
				if importMap[argName] != "" {
					argName = importMap[argName]
				}
				if arg2Type != argName+"."+xSelectExpr.Sel.Name {
					err = fmt.Errorf("%v: the 2st parameter type of method %v should be %v, not %v",
						fSet.Position(xCallExpr.Fun.Pos()), tableValue.Func, arg2Note, getArgType(xCallExpr.Args[1]))
					return
				}
				var column columnAst
				var columnList []columnAst
				for _, col := range funcArg.Body.List {
					var exprStmt *ast.ExprStmt
					if exprStmt, ok = col.(*ast.ExprStmt); ok {
						column, err = getColumn(fSet, exprStmt.X, internetCallName)
						if err != nil {
							return
						}
						// 验证字段
						switch tableValue.Type {
						case "mysql":
							var args []arg
							for _, v := range column.LianFuncSort {
								var argTypes []string
								for _, internetArg := range column.InternetFunc[v] {
									argTypes = append(argTypes, internetArg.Type)
								}
								args = append(args, arg{
									Type:     v,
									ArgTypes: argTypes,
								})
							}
							if err = validInternetFunc(tableValue.Type, funcArgValid[tableValue.Type][tableValue.Active], args); err != nil {
								err = fmt.Errorf("%v: %v",
									fSet.Position(exprStmt.X.Pos()), err.Error())
								return
							}
						}
						columnList = append(columnList, column)
					}
				}
				tableValue.ColumnList = columnList
			}

			tableList = append(tableList, tableValue)
		}
	}
	for k, v := range useFuncMap {
		err = fmt.Errorf("%v: %v declared and not used",
			fSet.Position(v), k)
		return
	}
	return
}

// 获取参数类型
func getArgType(arg ast.Expr) string {
	var ok bool
	switch val := arg.(type) {
	case *ast.BasicLit:
		return strings.ToLower(val.Kind.String())
	case *ast.FuncLit:
		var params []string
		for _, v := range val.Type.Params.List {
			list := strings.Split(getInternetCallName(v), " ")
			params = append(params, list[len(list)-1])
		}
		var results []string
		for _, value := range val.Type.Results.List {
			switch v := value.Type.(type) {
			case *ast.Ident:
				results = append(results, v.Name)
			}
		}
		rs := "func(" + strings.Join(params, ", ") + ")"
		if len(results) == 0 {
			return rs
		} else if len(results) == 1 {
			return rs + " " + results[0]
		}
		return rs + " (" + strings.Join(results, ", ") + ")"
	case *ast.CompositeLit:
		switch tp := val.Type.(type) {
		case *ast.SelectorExpr:
			var xIdent *ast.Ident
			if xIdent, ok = tp.X.(*ast.Ident); ok && tp.Sel != nil {
				return xIdent.Name + "." + tp.Sel.Name + "{}"
			}
		case *ast.Ident:
			return tp.Obj.Name + "{}"
		}
	case *ast.UnaryExpr:
		return val.Op.String() + getArgType(val.X)
	case *ast.Ident:
		if val.Obj.Decl != nil {
			var declStmt *ast.AssignStmt
			if declStmt, ok = val.Obj.Decl.(*ast.AssignStmt); ok {
				var types []string
				for _, v := range declStmt.Rhs {
					types = append(types, getArgType(v))
				}
				return strings.Join(types, ", ")
			}
		}
	}
	return fmt.Sprintf("%T", arg)
}
