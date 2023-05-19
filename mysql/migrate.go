package mysql

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	FuncType = struct {
		Up   string
		Down string
	}{Up: "Up", Down: "Down"}
	// 注释字符串
	siteList = [][]string{
		{"//"},
		{"/*", "*/"},
	}
)

func ParseSql(filePath string, funcType string) ([]string, error) {
	reg := new(Regexp)
	had := new(Handle)
	content, err := readAll(filePath)
	if err != nil {
		return nil, err
	}
	funcString, err := reg.getFunc(content, funcType)
	if err != nil {
		return nil, err
	}
	schemaList := reg.getSchema(funcString)
	sqlList := []string{}
	for _, schemaMap := range schemaList {
		tableType, _ := schemaMap["type"].(string)
		schemaContent, _ := schemaMap["content"].(string)
		schemaTable, _ := schemaMap["table"].(string)
		schemaOther, _ := schemaMap["other"].(string)
		columnList := reg.getColumn(schemaContent)
		switch tableType {
		case "Create":
			tableMap := reg.getOneColumn(schemaOther)
			sql, errs := had.mergeCreate(columnList, schemaTable)
			if errs != nil {
				return nil, errs
			}
			sql = fmt.Sprintf("CREATE TABLE `%s` (\n", schemaTable) + sql
			sql += ")"
			sql += had.handleColumnData(tableMap["Engine"], "", "", " ENGINE=?", "")
			sql += had.handleColumnData(tableMap["AutoIncrement"], "", "", " AUTO_INCREMENT=?", "")
			sql += had.handleColumnData(tableMap["Charset"], "", "", " DEFAULT CHARSET=?", "")
			sql += had.handleColumnData(tableMap["Collate"], "", "", " COLLATE=?", "")
			sql += had.handleColumnData(tableMap["Comment"], "'", "", " COMMENT=?", "")
			sqlList = append(sqlList, sql)
		case "CreateIfNotExists":
			tableMap := reg.getOneColumn(schemaOther)
			sql, errs := had.mergeCreate(columnList, schemaTable)
			if errs != nil {
				return nil, errs
			}
			sql = fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n", schemaTable) + sql
			sql += ")"
			sql += had.handleColumnData(tableMap["Engine"], "", "", " ENGINE=?", "")
			sql += had.handleColumnData(tableMap["AutoIncrement"], "", "", " AUTO_INCREMENT=?", "")
			sql += had.handleColumnData(tableMap["Charset"], "", "", " DEFAULT CHARSET=?", "")
			sql += had.handleColumnData(tableMap["Collate"], "", "", " COLLATE=?", "")
			sql += had.handleColumnData(tableMap["Comment"], "'", "", " COMMENT=?", "")
			sqlList = append(sqlList, sql)
		case "Table":
			sqls, errs := had.mergeUpdate(columnList, schemaTable)
			if errs != nil {
				return nil, errs
			}
			sqlList = append(sqlList, sqls...)
		case "Drop":
			sqlList = append(sqlList, fmt.Sprintf("DROP TABLE `%s`", schemaTable))
		case "DropIfExists":
			sqlList = append(sqlList, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", schemaTable))
		default:
			return nil, validErr("schema not method exists '%s'", tableType)
		}
	}
	return sqlList, nil
}

func readAll(filePath string) (string, error) {
	fi, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	isSite := false
	content := ""
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(a)
		// 过滤掉注释
		tLine := ""
		for _, site := range siteList {
			if len(site) == 1 {
				siteIndex := strings.Index(line, site[0])
				if siteIndex != -1 {
					line = line[0:siteIndex]
				}
			}
			if len(site) == 2 {
				siteStartIndex := strings.Index(line, site[0])
				siteStopIndex := strings.Index(line, site[1])
				if !isSite && siteStartIndex == -1 && siteStopIndex == -1 {
					tLine = line
				}
				if siteStartIndex != -1 && siteStopIndex != -1 {
					rep, _ := regexp.Compile(`/\*.*?\*/`)
					tLine = rep.ReplaceAllString(line, "")
				} else {
					if siteStartIndex != -1 {
						tLine = line[0:siteStartIndex]
						isSite = true
					}
					if siteStopIndex != -1 {
						tLine = line[siteStopIndex+len(site[1]):]
						isSite = false
					}
				}
			}
		}
		line = tLine
		content += line + "\n"
	}
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	return content, nil
}

func validErr(err string, args ...interface{}) error {
	err = fmt.Sprintf(err, args...)
	return errors.New(err)
}
