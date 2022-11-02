package mysql

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	sqlList, err := ParseSql("./example/table.go", FuncType.Up)
	if err != nil {
		t.Errorf("错误: %s", err.Error())
	}
	f, _ := os.OpenFile("./example/table.sql", os.O_WRONLY|os.O_CREATE, 0666)
	for _, sql := range sqlList {
		_, _ = f.Write([]byte(sql + ";\n"))
	}
	_ = f.Close()
}
