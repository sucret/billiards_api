package cmd

import (
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	"billiards/pkg/tool"

	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type generateOption struct {
	Type  string
	Value string
	Table string
}

// var generator =
var option generateOption

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		option.run()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVarP(&option.Table, "table", "t", "", "生成model")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (o *generateOption) run() {

	if o.Table != "" {
		generateTable()
	} else {
		db := mysql.New()

		var adminList []model.Admin
		db.Preload("RoleList").Preload("RoleList.NodeList", "node_id in(?)", []int{1, 2, 21}).Where("admin_id = 3").Find(&adminList)
		tool.Dump(adminList)
	}
}

func generateTable() {
	db := mysql.New()

	config := gen.Config{
		OutPath:      "pkg/mysql/query",
		ModelPkgPath: "pkg/mysql/model",

		// WithDefaultQuery 生成默认查询结构体(作为全局变量使用), 即`Q`结构体和其字段(各表模型)
		// WithoutContext 生成没有context调用限制的代码供查询
		// WithQueryInterface 生成interface形式的查询代码(可导出), 如`Where()`方法返回的就是一个可导出的接口类型
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,

		// 表字段可为 null 值时, 对应结体字段使用指针类型
		// FieldNullable: false, // generate pointer when field is nullable

		// 表字段默认值与模型结构体字段零值不一致的字段, 在插入数据时需要赋值该字段值为零值的, 结构体字段须是指针类型才能成功, 即`FieldCoverable:true`配置下生成的结构体字段.
		// 因为在插入时遇到字段为零值的会被GORM赋予默认值. 如字段`age`表默认值为10, 即使你显式设置为0最后也会被GORM设为10提交.
		// 如果该字段没有上面提到的插入时赋零值的特殊需要, 则字段为非指针类型使用起来会比较方便.
		FieldCoverable: false, // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values

		// 模型结构体字段的数字类型的符号表示是否与表字段的一致, `false`指示都用有符号类型
		FieldSignable: false, // detect integer field's unsigned type, adjust generated data type
		// 生成 gorm 标签的字段索引属性
		FieldWithIndexTag: false, // generate with gorm index tag
		// 生成 gorm 标签的字段类型属性
		FieldWithTypeTag: true, // generate with gorm column type tag
	}
	g := gen.NewGenerator(config)

	dataMap := map[string]func(detailType string) (dataType string){
		"timestamp": func(detailType string) (dataType string) { return "Time" },
		"datetime":  func(detailType string) (dataType string) { return "Time" },
	}

	g.WithDataTypeMap(dataMap)

	// 设置目标 db
	g.UseDB(db)

	allModel := g.GenerateAllTable()

	role := g.GenerateModel("role")
	node := g.GenerateModel("node")
	g.GenerateModel("task")
	g.GenerateModel("task_log")

	g.GenerateModel("admin_role")
	g.GenerateModel("role_node")

	// 用户表，关联角色
	admin := g.GenerateModel("admin", []gen.ModelOpt{
		// 用户和角色表关系 many2many
		gen.FieldRelate(field.Many2Many, "RoleList", role, &field.RelateConfig{GORMTag: "many2many:admin_role;foreignKey:AdminID;joinForeignKey:AdminID;joinReferences:RoleID"}),
	}...)

	// 角色表，关联用户 和 权限
	g.GenerateModel("role", []gen.ModelOpt{
		// 角色和用户表关系 many2many
		gen.FieldRelate(field.Many2Many, "AdminList", admin, &field.RelateConfig{GORMTag: "many2many:admin_role;foreignKey:RoleID;joinForeignKey:RoleID;joinReferences:AdminID"}),

		// 角色和权限表关系 many2many
		gen.FieldRelate(field.Many2Many, "NodeList", node, &field.RelateConfig{GORMTag: "many2many:role_node;foreignKey:RoleID;joinForeignKey:RoleID;joinReferences:NodeID"}),
	}...)

	terminal := g.GenerateModel("terminal")

	table := g.GenerateModel("table", []gen.ModelOpt{
		gen.FieldRelate(field.HasMany, "TerminalList", terminal, &field.RelateConfig{GORMTag: "hasmany:terminal;foreignKey:TableID;joinForeignKey:TableID;joinReferences:TableID"}),
	}...)

	g.GenerateModel("shop", []gen.ModelOpt{
		gen.FieldRelate(field.HasMany, "TableList", table, &field.RelateConfig{GORMTag: "hasmany:table;foreignKey:ShopID;joinForeignKey:ShopID;joinReferences:ShopID"}),
	}...)

	g.ApplyBasic(allModel...)
	g.Execute()
}
