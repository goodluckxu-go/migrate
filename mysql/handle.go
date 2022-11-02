package mysql

import (
	"fmt"
	"github.com/goodluckxu/migrate/mysql/schema/tb"
	"reflect"
	"strings"
)

type Handle struct {
}

func (h *Handle) mergeCreate(columnList []map[string][]string, table string) (sql string, err error) {
	typeList := h.getStructMethods(&tb.Types{})
	var oneSql string
	for k, columnMap := range columnList {
		if oneSql, err = h.mergeColumn(columnMap, typeList, table); err != nil {
			return
		}
		sql += oneSql
		if k+1 != len(columnList) {
			sql += ","
		}
		sql += "\n"
	}
	return
}

func (h *Handle) mergeUpdate(columnList []map[string][]string, table string) (sqlList []string, err error) {
	typeList := h.getStructMethods(&tb.Types{})
	var oneSql string
	for _, columnMap := range columnList {
		if oneSql, err = h.mergeColumn(columnMap, typeList, table); err != nil {
			return
		}
		sqlList = append(sqlList, oneSql)
	}
	return
}

func (h *Handle) mergeColumn(columnMap map[string][]string, typeList []string, table string) (sql string, err error) {
	validSort := columnMap["validSort"]
	if len(validSort) == 0 {
		return
	}
	columnType := h.getList(validSort, 1)
	switch validSort[0] {
	case "Column":
		sql += h.handleColumnData(columnMap["Column"], "`", ",", " ?", "")
		if !h.inArray(columnType, typeList) {
			err = validErr("wrong column type '%s'", columnType)
			return
		}
		sql += h.handleColumnData(columnMap[columnType], "'", ",",
			fmt.Sprintf(" %s?(??)", strings.ToLower(columnType)), "")
		sql += h.handleColumnData(columnMap["Unsigned"], "", "", " unsigned", "")
		sql += h.handleColumnData(columnMap["Zerofill"], "", "", " zerofill", "")
		sql += h.handleColumnData(columnMap["Nullable"], "", "", " NULL", " NOT NULL")
		sql += h.handleColumnData(columnMap["Default"], "'", "", " DEFAULT ?", "")
		sql += h.handleColumnData(columnMap["AutoIncrement"], "", "", " AUTO_INCREMENT", "")
		sql += h.handleColumnData(columnMap["Comment"], "'", "", " COMMENT ?", "")
	case "Indexes":
		sql += h.handleColumnData(columnMap["FULLTEXT"], "", "", " FULLTEXT", "")
		sql += h.handleColumnData(columnMap["UNIQUE"], "", "", " UNIQUE", "")
		sql += h.handleColumnData(columnMap["SPATIAL"], "", "", " SPATIAL", "")
		sql += h.handleColumnData(columnMap["Name"], "`", "", " KEY ?",
			h.handleColumnData(columnMap["Indexes"], "", "_", " KEY `?`", ""))
		sql += h.handleColumnData(columnMap["Indexes"], "`", ",", " ?(??)", "")
		sql += h.handleColumnData(columnMap["BTREE"], "", "", " USING BTREE", "")
		sql += h.handleColumnData(columnMap["HASH"], "", "", " USING HASH", "")
	case "PrimaryKey":
		sql += h.handleColumnData(columnMap["PrimaryKey"], "`", ",", " PRIMARY KEY ?(??) USING BTREE", "")
	case "ForeignKey":
		sql += h.handleColumnData(columnMap["Name"], "`", "", " CONSTRAINT ?",
			h.handleColumnData(columnMap["ForeignKey"], "", "_", " CONSTRAINT `?`", ""))
		sql += h.handleColumnData(columnMap["ForeignKey"], "`", ",", " FOREIGN KEY ?(??)", "")
		sql += h.handleColumnData(columnMap["QuoteTable"], "`", ",", " REFERENCES ?", "")
		sql += h.handleColumnData(columnMap["QuoteColumn"], "`", ",", " ?(??)", "")
		sql += h.handleColumnData(columnMap["DeleteCascade"], "`", ",", " ON DELETE CASCADE", "")
		sql += h.handleColumnData(columnMap["UpdateCascade"], "`", ",", " ON UPDATE CASCADE", "")
	case "AddColumn":
		sql += fmt.Sprintf("ALTER TABLE `%s`", table)
		sql += h.handleColumnData(columnMap["AddColumn"], "`", "", " ADD COLUMN ?", "")
		if !h.inArray(columnType, typeList) {
			err = validErr("wrong column type '%s'", columnType)
			return
		}
		sql += h.handleColumnData(columnMap[columnType], "'", ",",
			fmt.Sprintf(" %s?(??)", strings.ToLower(columnType)), "")
		sql += h.handleColumnData(columnMap["Unsigned"], "", "", " unsigned", "")
		sql += h.handleColumnData(columnMap["Zerofill"], "", "", " zerofill", "")
		sql += h.handleColumnData(columnMap["Nullable"], "", "", " NULL", " NOT NULL")
		sql += h.handleColumnData(columnMap["Default"], "'", "", " DEFAULT ?", "")
		sql += h.handleColumnData(columnMap["AutoIncrement"], "", "", " AUTO_INCREMENT", "")
		sql += h.handleColumnData(columnMap["Comment"], "'", "", " COMMENT ?", "")
		sql += h.handleColumnData(columnMap["First"], "", "", " FIRST", "")
		sql += h.handleColumnData(columnMap["After"], "`", "", " AFTER ?", "")
	case "ModifyColumn":
		sql += fmt.Sprintf("ALTER TABLE `%s`", table)
		sql += h.handleColumnData(columnMap["ModifyColumn"], "`", "", " MODIFY COLUMN ?", "")
		if !h.inArray(columnType, typeList) {
			err = validErr("wrong column type '%s'", columnType)
			return
		}
		sql += h.handleColumnData(columnMap[columnType], "'", ",",
			fmt.Sprintf(" %s?(??)", strings.ToLower(columnType)), "")
		sql += h.handleColumnData(columnMap["Unsigned"], "", "", " unsigned", "")
		sql += h.handleColumnData(columnMap["Zerofill"], "", "", " zerofill", "")
		sql += h.handleColumnData(columnMap["Nullable"], "", "", " NULL", " NOT NULL")
		sql += h.handleColumnData(columnMap["Default"], "'", "", " DEFAULT ?", "")
		sql += h.handleColumnData(columnMap["AutoIncrement"], "", "", " AUTO_INCREMENT", "")
		sql += h.handleColumnData(columnMap["Comment"], "'", "", " COMMENT ?", "")
		sql += h.handleColumnData(columnMap["First"], "", "", " FIRST", "")
		sql += h.handleColumnData(columnMap["After"], "`", "", " AFTER ?", "")
	case "ChangeColumn":
		sql += fmt.Sprintf("ALTER TABLE `%s`", table)
		sql += h.handleColumnData(columnMap["ChangeColumn"], "`", " ", " CHANGE COLUMN ?", "")
		if !h.inArray(columnType, typeList) {
			err = validErr("wrong column type '%s'", columnType)
			return
		}
		sql += h.handleColumnData(columnMap[columnType], "'", ",",
			fmt.Sprintf(" %s?(??)", strings.ToLower(columnType)), "")
		sql += h.handleColumnData(columnMap["Unsigned"], "", "", " unsigned", "")
		sql += h.handleColumnData(columnMap["Zerofill"], "", "", " zerofill", "")
		sql += h.handleColumnData(columnMap["Nullable"], "", "", " NULL", " NOT NULL")
		sql += h.handleColumnData(columnMap["Default"], "'", "", " DEFAULT ?", "")
		sql += h.handleColumnData(columnMap["AutoIncrement"], "", "", " AUTO_INCREMENT", "")
		sql += h.handleColumnData(columnMap["Comment"], "'", "", " COMMENT ?", "")
		sql += h.handleColumnData(columnMap["First"], "", "", " FIRST", "")
		sql += h.handleColumnData(columnMap["After"], "`", "", " AFTER ?", "")
	case "DropColumn":
		sql += fmt.Sprintf("ALTER TABLE `%s`", table)
		sql += h.handleColumnData(columnMap["DropColumn"], "`", "", " DROP COLUMN ?", "")
	case "AddIndexes":
		sql += fmt.Sprintf("ALTER TABLE `%s` ADD", table)
		sql += h.handleColumnData(columnMap["FULLTEXT"], "", "", " FULLTEXT", "")
		sql += h.handleColumnData(columnMap["UNIQUE"], "", "", " UNIQUE", "")
		sql += h.handleColumnData(columnMap["SPATIAL"], "", "", " SPATIAL", "")
		sql += h.handleColumnData(columnMap["Name"], "`", "", " KEY ?",
			h.handleColumnData(columnMap["AddIndexes"], "", "_", " KEY `?`", ""))
		sql += h.handleColumnData(columnMap["AddIndexes"], "`", ",", " ?(??)", "")
		sql += h.handleColumnData(columnMap["BTREE"], "", "", " USING BTREE", "")
		sql += h.handleColumnData(columnMap["HASH"], "", "", " USING HASH", "")
	case "DropIndexes":
		sql += fmt.Sprintf("ALTER TABLE `%s`", table)
		sql += h.handleColumnData(columnMap["DropIndexes"], "`", "", " DROP INDEX ?", "")
	case "AddPrimaryKey":
		sql += fmt.Sprintf("ALTER TABLE `%s`", table)
		sql += h.handleColumnData(columnMap["AddPrimaryKey"], "`", ",", " ADD PRIMARY KEY ?(??)", "")
	case "DropPrimaryKey":
		sql += fmt.Sprintf("ALTER TABLE `%s` DROP PRIMARY KEY", table)
	case "AddForeignKey":
		sql += fmt.Sprintf("ALTER TABLE `%s` ADD", table)
		sql += h.handleColumnData(columnMap["Name"], "`", "", " CONSTRAINT ?",
			h.handleColumnData(columnMap["AddForeignKey"], "", "_", " CONSTRAINT `?`", ""))
		sql += h.handleColumnData(columnMap["AddForeignKey"], "`", ",", " FOREIGN KEY ?(??)", "")
		sql += h.handleColumnData(columnMap["QuoteTable"], "`", ",", " REFERENCES ?", "")
		sql += h.handleColumnData(columnMap["QuoteColumn"], "`", ",", " ?(??)", "")
		sql += h.handleColumnData(columnMap["DeleteCascade"], "`", ",", " ON DELETE CASCADE", "")
		sql += h.handleColumnData(columnMap["UpdateCascade"], "`", ",", " ON UPDATE CASCADE", "")
	case "DropForeignKey":
		sql += fmt.Sprintf("ALTER TABLE `%s`", table)
		sql += h.handleColumnData(columnMap["DropForeignKey"], "`", "", " DROP CONSTRAINT ?", "")
	}
	return
}

func (h *Handle) handleColumnData(columnData []string, replaceNew interface{}, sep, sql, notSql string) string {
	if columnData == nil {
		return notSql
	}
	columnList := []string{}
	for _, v := range columnData {
		if v == "nil" {
			v = "NULL"
		}
		if replaceNew != nil {
			v = strings.ReplaceAll(v, "\"", fmt.Sprintf("%s", replaceNew))
		}
		columnList = append(columnList, v)
	}
	replaceList := []string{"?(", "?)", "?'"}
	for _, old := range replaceList {
		if len(columnList) > 0 {
			sql = strings.ReplaceAll(sql, old, strings.TrimPrefix(old, "?"))
		} else {
			sql = strings.ReplaceAll(sql, old, "")
		}
	}
	sql = strings.ReplaceAll(sql, "?", strings.Join(columnList, sep))
	return sql
}

func (h *Handle) getList(list []string, index int) string {
	if len(list) > index {
		return list[index]
	}
	return ""
}

func (h *Handle) getStructMethods(val interface{}) []string {
	value := reflect.ValueOf(val)
	rs := []string{}
	for i := 0; i < value.NumMethod(); i++ {
		rs = append(rs, value.Type().Method(i).Name)
	}
	return rs
}

func (h *Handle) inArray(val interface{}, array interface{}) (exists bool) {
	exists = false
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
				return
			}
		}
	}
	return
}
