package example

import (
	"github.com/goodluckxu-go/migrate/schema"
	"github.com/goodluckxu-go/migrate/schema/mysql"
	mysqlAlias "github.com/goodluckxu-go/migrate/schema/mysql" // 使用别名
	"github.com/goodluckxu-go/migrate/schema/mysql/tb"
	"github.com/goodluckxu-go/migrate/schema/pgsql"
	tb2 "github.com/goodluckxu-go/migrate/schema/pgsql/tb"
)

type MyTable struct {
}

func (m *MyTable) Up() {
	schema.Drop[mysql.Drop]("user")
	schema.Drop[mysqlAlias.DropIfExists]("user_info")
	schema.Edit[mysql.Create]("user", func(table *tb.CreateTable) {
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
	createIfNotExists("user_info", func(table *tb.CreateTable) {
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
	update("user", func(table *tb.UpdateTable) {
		table.ChangeColumn("tel", "phone").Varchar(50).Default("").Comment("手机号").After("name")
		table.ModifyColumn("password").Varchar(50).Default("").Comment("密码").After("phone")
		table.AddIndexes("phone", "password").BTREE()
	})
	schema.Edit[mysql.Update]("user_info", func(table *tb.UpdateTable) {
		table.DropForeignKey("test_a_b_c")
		table.DropIndexes("phone_password")
		table.ModifyColumn("id").Int().Unsigned().Comment("主键")
		table.DropPrimaryKey()
	})
	schema.Drop[mysql.Drop]("user_info")
	schema.Drop[mysql.Drop]("user")
}

func (m *MyTable) Down() {
	schema.Edit[pgsql.Create]("user_info", func(table *tb2.CreateTable) {
		table.Column("uuid").Uuid()
		table.Column("ids").Int4().Dimension(2)
		table.Column("num").Int4().GeneratedByDefaultAsIdentity(1, 1, 1).Comment("自增数字")
		table.Column("num1").Int4().GeneratedAlwaysAsIdentity(1, 1, 1).Comment("自增数字")
		table.Column("name").Varchar(100).CollateMode("pg_catalog").Collate("default").Comment("名称")
		table.Column("pwd").Varchar(50).Nullable().Default(nil).Comment("密码")
		table.Indexes(func(column *tb2.IndexesColumn) {
			column.Field("name").CollateMode("pg_catalog").Collate("default").
				OperationSymbolMode("pg_catalog").OperationSymbol("text_ops").ASC().NULLSLAST()
			column.Field("id").ASC()
		}).UNIQUE().BTREE().Comment("索引注释")
		table.PrimaryKey("uuid")
		table.ForeignKey("name").QuoteTable("user").QuoteColumn("id").DeleteCascade().Comment("外键注释")
		table.Schema.OwnerTo("postgres").Inherits("abc").ClusterOn("test_idx").Fillfactor(10).Comment("测试表")
	})
	schema.Edit[pgsql.Update]("user_info", func(table *tb2.UpdateTable) {
		table.AddColumn("user_add").Int4().GeneratedByDefaultAsIdentity(1, 1, 1).Comment("自增数字")
	})
}
