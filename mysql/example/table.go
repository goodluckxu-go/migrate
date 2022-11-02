package example

import (
	"github.com/goodluckxu-go/migrate/mysql/schema"
	"github.com/goodluckxu-go/migrate/mysql/schema/tb"
)

type MyTable struct {
}

func (m *MyTable) Up() {
	schema.DropIfExists("user")
	schema.DropIfExists("user_info")
	schema.Create("user", func(table tb.CreateTable) {
		table.Column("id").Int().Unsigned().AutoIncrement().Comment("主键")
		table.Column("name").Varchar(50).Default("").Comment("用户名")
		table.Column("tel").Varchar(20).Default("").Comment("电话号码")
		table.Column("password").Varchar(32).Default("").Comment("密码")
		table.Column("created_at").Timestamp().Nullable().Default(nil).Comment("创建时间")
		table.Column("update_at").Timestamp().Nullable().Default(nil).Comment("修改时间")
		table.PrimaryKey("id")
		table.Indexes("name").UNIQUE().BTREE()
	}).Engine("InnoDB").Charset("utf8mb4").Collate("utf8mb4_unicode_ci").Comment("用户表")
	schema.CreateIfNotExists("user_info", func(table tb.CreateTable) {
		table.Column("id").Int().Unsigned().AutoIncrement().Comment("主键")
		table.Column("user_id").Int().Unsigned().Default("0").Comment("用户uuid")
		table.Column("number").Int(3).Zerofill().Default("000").Comment("编号")
		table.Column("content").Json().Nullable().Default(nil).Comment("内容")
		table.Column("created_at").Timestamp().Nullable().Default(nil).Comment("创建时间")
		table.Column("update_at").Timestamp().Nullable().Default(nil).Comment("修改时间")
		table.PrimaryKey("id")
		table.ForeignKey("user_id").Name("test_a_b_c").QuoteTable("user").QuoteColumn("id").UpdateCascade().DeleteCascade()
	}).Engine("InnoDB").Charset("utf8mb4").Collate("utf8mb4_unicode_ci").Comment("用户详情表")
	schema.Table("user", func(table tb.Table) {
		table.ChangeColumn("tel", "phone").Varchar(50).Default("").Comment("手机号").After("name")
		table.ModifyColumn("password").Varchar(50).Default("").Comment("密码").After("phone")
		table.AddIndexes("phone", "password").BTREE()
	})
	schema.Table("user_info", func(table tb.Table) {
		table.DropForeignKey("test_a_b_c")
		table.DropIndexes("phone_password")
		table.ModifyColumn("id").Int().Unsigned().Comment("主键")
		table.DropPrimaryKey()
	})
	schema.Drop("user_info")
	schema.Drop("user")
}

func (m *MyTable) Down() {
}
