package migrate

import (
	"fmt"
	"strings"
)

type parseMySQL struct {
}

func (p parseMySQL) ParseSQL(table tableAst) (sqlList []string) {
	switch table.Active {
	case "Create":
		sql := "CREATE TABLE"
		sql += p.mergeCreate(table)
		sqlList = append(sqlList, sql)
	case "CreateIfNotExists":
		sql := "CREATE TABLE IF NOT EXISTS"
		sql += p.mergeCreate(table)
		sqlList = append(sqlList, sql)
	case "Update":
		for _, column := range table.ColumnList {
			sqlList = append(sqlList, p.mergeColumn(column, table.Name))
		}
	case "Drop":
		sql := "DROP TABLE"
		sql = fmt.Sprintf("%v `%v`", sql, table.Name)
		sqlList = append(sqlList, sql)
	case "DropIfExists":
		sql := "DROP TABLE"
		sql = fmt.Sprintf("%v IF EXISTS `%v`", sql, table.Name)
		sqlList = append(sqlList, sql)
	}
	return
}

func (p parseMySQL) mergeColumn(column columnAst, tableName string) (sql string) {
	firstDao := column.LianFuncSort[0]
	switch firstDao {
	case "Column":
		sql += p.handleColumnData(column.InternetFunc["Column"], "`", "", "", ",", " ?", "")
		sql += p.mergeField(column)
	case "Indexes":
		sql += p.mergeIndexes(column)
	case "PrimaryKey":
		sql += p.handleColumnData(column.InternetFunc["PrimaryKey"], "`", "", "", ",", " PRIMARY KEY (?) USING BTREE", "")
	case "ForeignKey":
		sql += p.mergeForeignKey(column)
	case "AddColumn":
		sql += "ALTER TABLE"
		sql = fmt.Sprintf("%s `%s`", sql, tableName)
		sql += p.handleColumnData(column.InternetFunc["AddColumn"], "`", "", "", "", " ADD COLUMN ?", "")
		sql += p.mergeField(column)
		sql += p.handleColumnData(column.InternetFunc["First"], "", "", "", "", " FIRST", "")
		sql += p.handleColumnData(column.InternetFunc["After"], "`", "", "", "", " AFTER ?", "")
	case "ModifyColumn":
		sql += "ALTER TABLE"
		sql = fmt.Sprintf("%s `%s`", sql, tableName)
		sql += p.handleColumnData(column.InternetFunc["ModifyColumn"], "`", "", "", "", " MODIFY COLUMN ?", "")
		sql += p.mergeField(column)
		sql += p.handleColumnData(column.InternetFunc["First"], "", "", "", "", " FIRST", "")
		sql += p.handleColumnData(column.InternetFunc["After"], "`", "", "", "", " AFTER ?", "")
	case "ChangeColumn":
		sql += "ALTER TABLE"
		sql = fmt.Sprintf("%s `%s`", sql, tableName)
		sql += p.handleColumnData(column.InternetFunc["ChangeColumn"], "`", "", "", " ", " CHANGE COLUMN ?", "")
		sql += p.mergeField(column)
		sql += p.handleColumnData(column.InternetFunc["First"], "", "", "", "", " FIRST", "")
		sql += p.handleColumnData(column.InternetFunc["After"], "`", "", "", "", " AFTER ?", "")
	case "DropColumn":
		sql += "ALTER TABLE"
		sql = fmt.Sprintf("%s `%s`", sql, tableName)
		sql += p.handleColumnData(column.InternetFunc["DropColumn"], "`", "", "", "", " DROP COLUMN ?", "")
	case "AddIndexes":
		sql += "ALTER TABLE"
		sql = fmt.Sprintf("%s `%s` ADD", sql, tableName)
		sql += p.mergeIndexes(column)
	case "DropIndexes":
		sql += "ALTER TABLE"
		sql = fmt.Sprintf("%s `%s`", sql, tableName)
		sql += p.handleColumnData(column.InternetFunc["DropIndexes"], "`", "", "", "", " DROP INDEX ?", "")
	case "AddPrimaryKey":
		sql += "ALTER TABLE"
		sql = fmt.Sprintf("%s `%s`", sql, tableName)
		sql += p.handleColumnData(column.InternetFunc["AddPrimaryKey"], "`", "", "", ",", " ADD PRIMARY KEY (?) USING BTREE", "")
	case "DropPrimaryKey":
		sql += "ALTER TABLE"
		sql = fmt.Sprintf("%s `%s` DROP PRIMARY KEY", sql, tableName)
	case "AddForeignKey":
		sql += "ALTER TABLE"
		sql = fmt.Sprintf("%s `%s` ADD", sql, tableName)
		sql += p.mergeForeignKey(column)
	case "DropForeignKey":
		sql += "ALTER TABLE"
		sql = fmt.Sprintf("%s `%s`", sql, tableName)
		sql += p.handleColumnData(column.InternetFunc["DropForeignKey"], "`", "", "", "", " DROP CONSTRAINT ?", "")
	}
	return
}

func (p parseMySQL) mergeCreate(table tableAst) (sql string) {
	sql = fmt.Sprintf("%s `%s` (\n", sql, table.Name)
	schemaColumn := columnAst{
		InternetFunc: map[string][]argAst{},
	}
	for index, column := range table.ColumnList {
		columnSql := p.mergeColumn(column, table.Name)
		if columnSql == "" {
			lianLen := len(column.LianFuncSort)
			schemaColumn.LianFuncSort = column.LianFuncSort
			if column.LianFuncSort[0] == "Schema" {
				for i := 1; i < lianLen; i++ {
					schemaColumn.InternetFunc[column.LianFuncSort[i]] = column.InternetFunc[column.LianFuncSort[i]]
				}
			}
			continue
		}
		if index != 0 {
			sql += ",\n"
		}
		sql += columnSql
	}
	sql += "\n)"
	sql += p.mergeSchema(schemaColumn)
	return
}

func (p parseMySQL) mergeField(column columnAst) (sql string) {
	firstDao := column.LianFuncSort[0]
	switch firstDao {
	case "Column", "AddColumn", "ModifyColumn", "ChangeColumn":
		columnType := column.LianFuncSort[1]
		sql += p.handleColumnData(column.InternetFunc[columnType], "", "", "(?)", ",",
			fmt.Sprintf(" %v?", strings.ToLower(columnType)), "")
		sql += p.handleColumnData(column.InternetFunc["Charset"], "", "", "", "", " CHARACTER SET ?", "")
		sql += p.handleColumnData(column.InternetFunc["Collate"], "", "", "", "", " COLLATE ?", "")
		sql += p.handleColumnData(column.InternetFunc["Unsigned"], "", "", "", "", " unsigned", "")
		sql += p.handleColumnData(column.InternetFunc["Zerofill"], "", "", "", "", " zerofill", "")
		sql += p.handleColumnData(column.InternetFunc["Nullable"], "", "", "", "", " NULL", " NOT NULL")
		sql += p.handleColumnData(column.InternetFunc["Default"], "'", "", "", "", " DEFAULT ?", "")
		sql += p.handleColumnData(column.InternetFunc["AutoIncrement"], "", "", "", "", " AUTO_INCREMENT", "")
		sql += p.handleColumnData(column.InternetFunc["Comment"], "'", "", "", "", " COMMENT ?", "")
	}
	return
}

func (p parseMySQL) mergeIndexes(column columnAst) (sql string) {
	firstDao := column.LianFuncSort[0]
	switch firstDao {
	case "Indexes", "AddIndexes":
		sql += p.handleColumnData(column.InternetFunc["FULLTEXT"], "", "", "", "", " FULLTEXT", "")
		sql += p.handleColumnData(column.InternetFunc["UNIQUE"], "", "", "", "", " UNIQUE", "")
		sql += p.handleColumnData(column.InternetFunc["SPATIAL"], "", "", "", "", " SPATIAL", "")
		sql += p.handleColumnData(column.InternetFunc["Name"], "`", "", "", "", " KEY ?",
			p.handleColumnData(column.InternetFunc[firstDao], "", "", "", "_", " KEY `?`", ""))
		sql += p.handleColumnData(column.InternetFunc[firstDao], "`", "", "", ",", " (?)", "")
		sql += p.handleColumnData(column.InternetFunc["BTREE"], "", "", "", "", " USING BTREE", "")
		sql += p.handleColumnData(column.InternetFunc["HASH"], "", "", "", "", " USING HASH", "")
	}
	return
}

func (p parseMySQL) mergeForeignKey(column columnAst) (sql string) {
	firstDao := column.LianFuncSort[0]
	switch firstDao {
	case "ForeignKey", "AddForeignKey":
		sql += p.handleColumnData(column.InternetFunc["Name"], "`", "", "", "", " CONSTRAINT ?",
			p.handleColumnData(column.InternetFunc["ForeignKey"], "", "", "", "_", " CONSTRAINT `?`", ""))
		sql += p.handleColumnData(column.InternetFunc["ForeignKey"], "`", "", "", ",", " FOREIGN KEY (?)", "")
		sql += p.handleColumnData(column.InternetFunc["QuoteTable"], "`", "", "", ",", " REFERENCES ?", "")
		sql += p.handleColumnData(column.InternetFunc["QuoteColumn"], "`", "", "", ",", " (?)", "")
		sql += p.handleColumnData(column.InternetFunc["DeleteCascade"], "`", "", "", ",", " ON DELETE CASCADE", "")
		sql += p.handleColumnData(column.InternetFunc["UpdateCascade"], "`", "", "", ",", " ON UPDATE CASCADE", "")
	}
	return
}

func (p parseMySQL) mergeSchema(column columnAst) (sql string) {
	if len(column.LianFuncSort) == 0 {
		return
	}
	firstDao := column.LianFuncSort[0]
	switch firstDao {
	case "Schema":
		sql += p.handleColumnData(column.InternetFunc["Engine"], "", "", "", "", " ENGINE=?", "")
		sql += p.handleColumnData(column.InternetFunc["AutoIncrement"], "", "", "", "", " AUTO_INCREMENT=?", "")
		sql += p.handleColumnData(column.InternetFunc["Charset"], "", "", "", "", " DEFAULT CHARSET=?", "")
		sql += p.handleColumnData(column.InternetFunc["Collate"], "", "", "", "", " COLLATE=?", "")
		sql += p.handleColumnData(column.InternetFunc["Comment"], "'", "", "", "", " COMMENT=?", "")
	}
	return
}

func (p parseMySQL) handleColumnData(funcArgs []argAst, argSingleAround, argMulAround, argsExistsAround, sep, sql, defaultSql string) string {
	if funcArgs == nil {
		return defaultSql
	}
	var argList []string
	for _, arg := range funcArgs {
		ar := arg.Val
		if arg.Type == "nil" {
			ar = "NULL"
		} else if argSingleAround != "" {
			ar = argSingleAround + ar + argSingleAround
		}
		argList = append(argList, ar)
	}
	argStr := strings.Join(argList, sep)
	if len(funcArgs) > 1 && argMulAround != "" {
		argStr = strings.ReplaceAll(argMulAround, "?", argStr)
	}
	if len(funcArgs) > 0 && argsExistsAround != "" {
		argStr = strings.ReplaceAll(argsExistsAround, "?", argStr)
	}
	sql = strings.ReplaceAll(sql, "?", argStr)
	return sql
}
