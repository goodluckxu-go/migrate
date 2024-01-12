package migrate

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	sqlMapList, err := ParseSQL("./example/table.go", []string{"Up"})
	fmt.Println(err)
	return
	if err != nil {
		t.Errorf("错误: %s", err.Error())
	}
	f, _ := os.OpenFile("./example/table.sql", os.O_WRONLY|os.O_CREATE, 0666)
	for funcName, sqlAll := range sqlMapList {
		_, _ = f.Write([]byte("### func name " + funcName + " ###\n"))
		for sqlType, sqlList := range sqlAll {
			_, _ = f.Write([]byte("### " + sqlType + " ###\n"))
			_, _ = f.Write([]byte(strings.Join(sqlList, ";\n")))
		}
	}
	_ = f.Close()
}
