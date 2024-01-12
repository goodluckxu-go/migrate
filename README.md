# go版本的migrate创建

使用ast语法解析器解析go语言生成sql语句，仿照php语言laravel框架的artisan命令的sql生成方式

## 实现(achieve)
- 使用go语言ast语法糖实现源码解读
- go语言解析，不符合定义规范的语法报错误的文件以及对应文件的行列未知
- go语言解析，引入可使用别名
- go语言解析，基层方法可赋值 方法有[schema.Drop,schema.Edit]
- go语言解析，schema.Drop,schema.Edit两个方法可以泛型


## 用法(usage)
~~~go
package example

import (
	"github.com/goodluckxu-go/migrate/schema"
	"github.com/goodluckxu-go/migrate/schema/mysql"
	mysqlAlias "github.com/goodluckxu-go/migrate/schema/mysql" // 使用别名
	"github.com/goodluckxu-go/migrate/schema/mysql/tb"
)

type MyTable struct {
}

func (m *MyTable) Up() {
	schema.Drop[mysql.Drop]("user")
	schema.Drop[mysqlAlias.DropIfExists]("user_info")
	schema.Edit[mysql.Create]("user", func(table tb.CreateTable) {
		table.Column("id").Int().Unsigned().AutoIncrement().Comment("主键")
		table.Column("name").Varchar(50).Default("").Comment("用户名")
		table.Column("tel").Varchar(20).Default("").Comment("电话号码")
		table.Column("password").Varchar(32).Default("").Comment("密码")
		table.Column("created_at").Timestamp().Nullable().Default(nil).Comment("创建时间")
		table.Column("update_at").Timestamp().Nullable().Default(nil).Comment("修改时间")
		table.PrimaryKey("id")
		table.Indexes("name").UNIQUE().BTREE()
		table.Schema.Engine("InnoDB").Charset("utf8mb4")          // 设置表属性
		table.Schema.Collate("utf8mb4_unicode_ci").Comment("用户表") // 设置表属性
	})
	// := 赋值
	createIfNotExists := schema.Edit[mysql.CreateIfNotExists]
	createIfNotExists("user_info", func(table tb.CreateTable) {
		table.Column("id").Int().Unsigned().AutoIncrement().Comment("主键")
		table.Column("user_id").Int().Unsigned().Default("0").Comment("用户uuid")
		table.Column("number").Int(3).Zerofill().Default("000").Comment("编号")
		table.Column("content").Json().Nullable().Default(nil).Comment("内容")
		table.Column("created_at").Timestamp().Nullable().Default(nil).Comment("创建时间")
		table.Column("update_at").Timestamp().Nullable().Default(nil).Comment("修改时间")
		table.PrimaryKey("id")
		table.ForeignKey("user_id").Name("test_a_b_c").QuoteTable("user").QuoteColumn("id").UpdateCascade().DeleteCascade()
		table.Schema.Engine("InnoDB").Charset("utf8mb4").Collate("utf8mb4_unicode_ci").Comment("用户表") // 设置表属性
	})
	// var 赋值
	var update = schema.Edit[mysql.Update]
	update("user", func(table tb.UpdateTable) {
		table.ChangeColumn("tel", "phone").Varchar(50).Default("").Comment("手机号").After("name")
		table.ModifyColumn("password").Varchar(50).Default("").Comment("密码").After("phone")
		table.AddIndexes("phone", "password").BTREE()
	})
	schema.Edit[mysql.Update]("user_info", func(table tb.UpdateTable) {
		table.DropForeignKey("test_a_b_c")
		table.DropIndexes("phone_password")
		table.ModifyColumn("id").Int().Unsigned().Comment("主键")
		table.DropPrimaryKey()
	})
	schema.Drop[mysql.Drop]("user_info")
	schema.Drop[mysql.Drop]("user")
}

func (m *MyTable) Down() {
}

~~~

执行该文件代码
~~~go
sqlMapList, err := ParseSQL("./example/table.go", []string{"Up"})
if err != nil {
	log.Fatal(err)
}
for funcName, sqlAll := range sqlMapList {
	fmt.Println("### func name " + funcName + " ###")
	for sqlType, sqlList := range sqlAll {
		fmt.Println("### " + sqlType + " ###")
		fmt.Println(strings.Join(sqlList, ";\n"))
	}
}
~~~

执行结果
~~~sql
### func name Up ###
### mysql ###
DROP TABLE `user`;
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user` (
 `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
 `name` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
 `tel` varchar(20) NOT NULL DEFAULT '' COMMENT '电话号码',
 `password` varchar(32) NOT NULL DEFAULT '' COMMENT '密码',
 `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
 `update_at` timestamp NULL DEFAULT NULL COMMENT '修改时间',
 PRIMARY KEY (`id`) USING BTREE,
 UNIQUE KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
CREATE TABLE IF NOT EXISTS `user_info` (
 `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
 `user_id` int unsigned NOT NULL DEFAULT '0' COMMENT '用户uuid',
 `number` int(3) zerofill NOT NULL DEFAULT '000' COMMENT '编号',
 `content` json NULL DEFAULT NULL COMMENT '内容',
 `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
 `update_at` timestamp NULL DEFAULT NULL COMMENT '修改时间',
 PRIMARY KEY (`id`) USING BTREE,
 CONSTRAINT `test_a_b_c` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
ALTER TABLE `user` CHANGE COLUMN `tel` `phone` varchar(50) NOT NULL DEFAULT '' COMMENT '手机号' AFTER `name`;
ALTER TABLE `user` MODIFY COLUMN `password` varchar(50) NOT NULL DEFAULT '' COMMENT '密码' AFTER `phone`;
ALTER TABLE `user` ADD KEY `phone_password` (`phone`,`password`) USING BTREE;
ALTER TABLE `user_info` DROP CONSTRAINT `test_a_b_c`;
ALTER TABLE `user_info` DROP INDEX `phone_password`;
ALTER TABLE `user_info` MODIFY COLUMN `id` int unsigned NOT NULL COMMENT '主键';
ALTER TABLE `user_info` DROP PRIMARY KEY;
DROP TABLE `user_info`;
DROP TABLE `user`

~~~

## 问题(problem)

目前只支持mysql的解析

后期陆续更新其他sql
