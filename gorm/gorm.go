package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	dm "github.com/go-sql-driver/mysql"
	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/hints"
	"gorm.io/plugin/dbresolver"
	"gorm.io/plugin/prometheus"
	"gorm.io/plugin/soft_delete"
	"gorm.io/rawsql"
	"gorm.io/sharding"
)

func main() {}

type Bool bool

func (b *Bool) Scan(value any) error {
	if v, ok := value.([]byte); ok {
		*b = len(v) > 0 && v[0] > 0
		return nil
	} else {
		return fmt.Errorf("cannot convert %v to bool", value)
	}
}

func (b *Bool) Value() (driver.Value, error) {
	if *b {
		return []byte{1}, nil
	} else {
		return nil, nil
	}
}

// GORM 支持的结构字段类型：基础类型、指针、类型别名、实现了 Scanner 和 Valuer 的自定义类型
// 表名和字段名默认采用 snake_cases 风格
type Person struct {
	Id           int // 该字段命名，将识别为主码
	Email        *string
	Hobby        sql.NullString
	ignored      string         // 忽略不导出字段
	CreatedAt    time.Time      // 该字段命名，当添加记录且字段值是零值时，自动设置为当前时间
	UpdatedAt    time.Time      // 该字段命名，当更新记录且字段值是零值时，自动设置为当前时间。标签设置为 gorm:"autoUpdateTime:false 停用该行为
	UpdatedNano  int64          `gorm:"autoUpdateTime:nano"`  // 设置该标签，更新记录时自动设置为纳秒时间戳
	UpdatedMilli int64          `gorm:"autoUpdateTime:milli"` // 设置该标签，更新记录时自动设置为毫秒时间戳
	Created      int64          `gorm:"autoCreateTime"`       // 设置该标签，添加记录时自动设置为秒时间戳
	DeletedAt    gorm.DeletedAt // 该字段类型，将识别为软删除字段
	Name1        string         `gorm:"-"`           // 设置该标签，该字段忽略读写
	Name2        string         `gorm:"-:all"`       // 设置该标签，该字段忽略读写和建表
	Name3        string         `gorm:"-:migration"` // 设置该标签，该字段忽略建表
	Name4        string         `gorm:"->"`          // 设置该标签，该字段只读取
	Name5        string         `gorm:"<-"`          // 设置该标签，该字段可读取和写入
	Name6        string         `gorm:"<-:false"`    // 设置该标签，该字段只读取
	Name7        string         `gorm:"<-:create"`   // 设置该标签，该字段可读取和新增
	Name8        string         `gorm:"->:false"`    // 设置该标签，该字段忽略可读取
	BlogID       int
	Blog                               // 内嵌结构体相当于字段位于父结构体效果
	Blog1        Blog                  `gorm:"embedded"`                      // 设置该标签，相当于内嵌字段效果
	Blog2        Blog                  `gorm:"embedded;embeddedPrefix:blog_"` // 设置该标签，内嵌字段的数据库字段名有前缀 blog_
	gorm.Model                         // 提供的通用字段
	DeletedAt1   soft_delete.DeletedAt `gorm:"softDelete:milli;DeletedAtField:DeletedAt"` // 该字段类型，将识别为软删除字段。默认单位秒，使用标签改变行为，milli nano flag。DeletedAtField 指定同步设置的删除时间字段
	GormValue

	// 标签大小写不敏感，使用分号分隔，使用反斜线转义特殊字符
	// column 指定数据库表字段名
	// type 指定数据库表字段类型
	// serializer 设置序列化和反序列化方式，当写和读数据库时。json gob unixtime
	// size 指定数据库表字段大小长度
	// primaryKey 指定字段是主码。默认设置自增
	// unique 指定数据库表字段为唯一约束
	// default 指定数据库表字段默认值。当添加记录时，字段是零值时，会使用这里的默认值
	// precision 指定数据库表字段精度
	// scale 指定数据库表字段大小
	// not null 指定数据库表字段非空
	// autoIncrement 指定数据库表字段可自增。autoIncrement:false
	// autoIncrementIncrement 指定数据库表字段步长
	// embedded 设置字段为内嵌字段效果
	// embeddedPrefix 设置内嵌字段的数据库表字段名前缀
	// autoCreateTime 设置字段跟踪添加时间
	// autoUpdateTime 设置字段跟踪更新时间
	// index 创建索引，多个字段的索引名相同意味着创建联合索引。里面有字段：class, type, where, comment, expression, sort, collate, option, unique, priority, composite
	// uniqueIndex 创建唯一索引
	// check 指定数据库表字段检查约束
	// <- 设置字段写入权限
	// -> 设置字段读取权限
	// comment 指定数据库表字段注释
	// foreignKey 指定外码在本结构体的字段名，默认是 BlogID
	// references 指定关联结构的外码字段名，默认是 Id
	// constraint 建表约束，例子：OnUpdate:CASCADE,OnDelete:SET NULL;
	// many2many 连接表名
	// joinForeignKey
	// joinReferences
	// polymorphic
	// polymorphicType
	// polymorphicId
	// polymorphicValue
}

type Blog struct {
	Id       int
	PersonId int
}

// 钩子函数
func (*Person) BeforeCreate(tx *gorm.DB) error {
	// 更改 SQL 语句
	tx.Statement.AddClause(nil)
	tx.Statement.Select("")

	// 获取字段 GormDataType 值
	_ = tx.Statement.Schema.LookUpField("GormValue").DataType == "string"
	return nil
}
func (*Person) AfterCreate(tx *gorm.DB) error { return nil }
func (*Person) BeforeSave(tx *gorm.DB) error  { return nil }
func (*Person) AfterSave(tx *gorm.DB) error   { return nil }
func (*Person) AfterFind(tx *gorm.DB) error   { return nil }
func (*Person) AfterUpdate(tx *gorm.DB) error { return nil }
func (*Person) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Name1") {
		// 如果字段 name1 将被更新
	}

	// 如果字段 name1 name2 有一个要更新
	if tx.Statement.Changed("Name1", "Name2") {
		tx.Statement.SetColumn("Name3", "ww") // 设置更新的值
	}

	// 如果有任何字段要更新
	if tx.Statement.Changed() {
		tx.Statement.SetColumn("UpdatedAt", time.Now())
	}
	return nil
}
func (*Person) AfterDelete(tx *gorm.DB) error  { return nil }
func (*Person) BeforeDelete(tx *gorm.DB) error { return nil }
func (*Person) TableName() string {
	return "person" // 该表名会被 GORM 缓存，可使用 Scopes 动态表名
}

func Connect() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"),
		&gorm.Config{
			CreateBatchSize:                          1,    // 影响批量添加记录行为
			QueryFields:                              true, // 影响智能字段选择
			AllowGlobalUpdate:                        true, // 允许没有 where 条件更新
			TranslateError:                           true, // Error 转换为 GORM 的
			PrepareStmt:                              true,
			SkipDefaultTransaction:                   true,
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   newLogger,
			DisableNestedTransaction:                 false,
			DisableAutomaticPing:                     true,
		})

	// 复用连接
	sqlDB, _ := db.DB()
	gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

var db *gorm.DB

func Create() {
	// 添加记录，会回填主码字段
	db.Create(&Person{})

	// 主码有值就更新所有字段，没值就添加记录。零值字段也保存更新了
	db.Save(&Person{})

	// 添加多条记录，里面执行一条 SQL 语句完成
	result := db.Create([]*Person{})
	_ = result.RowsAffected // 返回添加的记录数

	// 只添加选择的字段
	db.Select("Email").Create(&Person{})
	// 选择更新关联表字段
	db.Select("Blog.Name").Create(&Person{})

	// 忽略字段添加
	db.Omit("Email").Create(&Person{})

	// 忽略关联表和连接表添加
	db.Omit("Blog").Create(&Person{})
	db.Omit(clause.Associations).Create(&Person{})
	// 忽略关联表添加
	db.Omit("Blog.*").Create(&Person{})
	// 忽略关联表字段更新
	db.Omit("Blog.Name").Create(&Person{})

	// 每次添加指定数量的记录
	db.CreateInBatches([]*Person{}, 1)

	// 使用 map 和 []map 添加记录，不会触发钩子函数和主码回填
	db.Model(&Person{}).Create(map[string]any{})
	db.Model(&Person{}).Create([]map[string]any{})

	// 添加关联表记录
	db.Model(&Person{Id: 1}).Association("Blog").Append(&Blog{})
	db.Model(&Person{Id: 1}).Association("Blog").Append(&[]*Blog{})

	// Blog 不是零值，会创建 SQL 语句添加记录，它的钩子函数也会被调用
	// 开启事务添加
	db.Create(&Person{})

	// 检测是否有联系
	_ = db.Model(&Person{}).Association("Blog").Error
}

type GormValue struct{}

func (*GormValue) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "CONCAT(?, '@126.com')",
		Vars: []any{"ivfzhou"},
	}
}

func (*GormValue) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "string"
	case "postgres":
		return "string"
	}
	return ""
}

func (*GormValue) GormDataType() string {
	return "string"
}

func UseExpr() {
	db.Model(&Person{}).Create(map[string]any{"Email": clause.Expr{
		SQL:  "CONCAT(?, '@126.com')",
		Vars: []any{"ivfzhou"},
	}})

	db.Model(&Person{}).Create(map[string]any{"Email": &GormValue{}})
}

func Conflict() {
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&Person{})

	// 冲突时更新 name1 字段为 zs
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]any{"name1": gorm.Expr("?", "zs")}),
	}).Create(&Person{})

	// 冲突时更新 name1 字段为新设置的值
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name1"}),
	}).Create(&Person{})

	// 发生冲突时，更新所有字段
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&Person{})
}

func Retrieve() {
	// 查询一个，order by id limit 1
	db.First(&Person{})

	// 查询最后一个，order by id desc limit 1
	db.Last(&Person{})

	// 以上的排序字段获取自结构体 Person，或者 db.Model()。如果结构体中没有主码则使用第一个字段排序

	// 获取一个，limit 1
	result := db.Take(&Person{})
	_ = result.RowsAffected                         // 查询到的记录数量
	errors.Is(result.Error, gorm.ErrRecordNotFound) // 是否未找到记录

	// 将结果扫描到结构体中
	db.Model(&Person{}).Scan(&Person{})

	// 使用主键查询
	db.First(&Person{}, 1)
	db.First(&Person{}, "1")
	db.Find(&Person{}, []int{1, 2, 3})
	db.First(&Person{Id: 1})

	// 查询包含软删除的记录
	db.Unscoped().Find(&Person{})

	// 使用查询条件
	db.Where("id = ?", 1).First(&Person{})
	db.Where("id in ?", []int{1}).First(&Person{})
	db.Where(&Person{Name1: "zs", Name2: ""}).First(&Person{})                                // 忽略结构体零值字段，where name1 = 'zs'
	db.Where(&Person{Name1: "zs"}, "name1", "Name2").First(&Person{})                         // 指定使用结构哪些字段作为查询条件，where name1 = 'zs' and name2 = ''
	db.Where(map[string]any{"name1": "zs", "name2": ""}).First(&Person{})                     // where name1 = 'zs' and name2 = ''
	db.Where([]int{1}).First(&Person{})                                                       // where id in (1)
	db.Where(db.Where("id = 1")).First(&Person{})                                             // where id = 1
	db.Where("(name1, name2) in ?", [][]any{{"zs", "admin"}, {"ls", "user"}}).Find(&Person{}) // where (name, name2) in (("zs", "admin"), ("ls", "user"))
	db.Table("person").Row()
	rows, _ := db.Table("person").Rows()
	defer rows.Close()
	rows.Next()
	db.ScanRows(rows, &Person{})

	db.Find(&Person{}, "id = ?", 1)
	db.Find(&Person{}, &Person{Name1: "zs"})
	db.Find(&Person{}, db.Where("id = 1"))
	db.Find(&Person{}, map[string]any{"name1": "zs", "name2": ""})

	db.Not("id = ?", 1).Find(&Person{})                                  // where not id = 1
	db.Not([]int{1}).Find(&Person{})                                     // where id not in (1)
	db.Not(&Person{Name1: "zs"}).Find(&Person{})                         // where name1 != 'zs'
	db.Not(map[string]any{"name1": "zs"}).Find(&Person{})                // where name1 != 'zs'
	db.Where("id = 1").Or("id = ?", 1).Find(&Person{})                   // where id = 1 or id = 1
	db.Where("id = 1").Or(&Person{Name1: "zs"}).Find(&Person{})          // where id = 1 or name1 = 'zs'
	db.Where("id = 1").Or(map[string]any{"name1": "zs"}).Find(&Person{}) // where id = 1 or name1 = 'zs'

	// 选择查询字段
	db.Select("name1").Find(&Person{})
	db.Select([]string{"name1"}).Find(&Person{})
	db.Select("coalesce(name1, ?)", "zs").Find(&Person{})

	// 根据目标结构体选择查询的字段，find 的默认行为
	type ps struct {
		Name string
	}
	db.Session(&gorm.Session{QueryFields: true}).Model(&Person{}).Find(&ps{})

	// 排序
	db.Order("id desc, name1").Find(&Person{})         // order by id desc name1
	db.Order("id desc").Order("name1").Find(&Person{}) // order by id desc name1
	db.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "field(id,?)", Vars: []any{[]int{1, 2, 3}}, WithoutParentheses: true},
	}).Find(&Person{}) // order by field(id,1,2,3)

	// 分页
	db.Limit(1).Find(&Person{})
	db.Offset(1).Find(&Person{})

	// 分组过滤
	db.Group("name1").Having("name1 = ?", "zs").Find(&Person{})

	// 联合查询
	db.Table("person").Joins("left join person on id = id and name1 = ?", "zs").Scan(&Person{})
	db.Joins("Blog").Find(&Person{})
	db.Joins("Blog.User").Find(&Person{})
	db.Joins("Blog", db.Where("")).Find(&Person{}) // 连接条件
	db.InnerJoins("Blog").Find(&Person{})
	db.Joins("left join (?) t on t.id = id", db.Find(&Person{})).Find(&Person{}) // 连接衍生表

	// 锁表查询
	db.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&Person{}) // from person for update
	db.Clauses(clause.Locking{
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).Find(&Person{}) // from person for share of person
	db.Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "NOWAIT",
	}).Find(&Person{}) // from person for update nowait
	db.Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "SKIP LOCKED",
	}).Find(&Person{}) // from person for update skip locked

	// 子查询
	db.Where("id in ?", db.Select("id").Find(&Person{})).Find(&Person{})

	// 命令参数
	db.Where("name1 = @name or name2 = @name", sql.Named("name", "zs")).Find(&Person{})
	db.Where("name1 = @name or name2 = @name", map[string]any{"name": "zs"}).First(&Person{})
	type Named struct {
		Name string
	}
	db.Where("name1 = @Name or name2 = @Name", &Named{"zs"}).First(&Person{})

	// 查询到 map
	var res map[string]any
	db.Model(&Person{}).First(&res, "id = ?", 1)
	var results []map[string]any
	db.Table("person").Find(&results)

	// 查询不到就添加记录
	db.FirstOrInit(&Person{})
	// Attrs：当查不到记录时，在 Person 中设置 Name1=zs
	db.Attrs(&Person{Name1: "zs"}).FirstOrInit(&Person{})
	// Assign：不管记录是否查询，都在 Person 中设置 Name1=zs
	db.Assign(&Person{Name1: "zs"}).FirstOrInit(&Person{})

	// 查不到记录就创建记录
	db.FirstOrCreate(&Person{})
	// 当记录没查到，person 表字段 name1 设置为 zs
	db.Attrs(&Person{Name1: "zs"}).FirstOrCreate(&Person{})
	// 不管记录查到与否，都将 person 表字段 name1 设置为 zs
	db.Assign(&Person{Name1: "zs"}).FirstOrCreate(&Person{})

	// 指令
	db.Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find(&Person{}) // select * /*+ MAX_EXECUTION_TIME(10000) */ from person
	db.Clauses(hints.New("hint")).Find(&Person{})                      // select * /*+ hint */ from person

	// 选择索引
	db.Clauses(hints.UseIndex("idx_name1")).Find(&Person{})                          // select * from person use index (idx_name1)
	db.Clauses(hints.ForceIndex("idx_name1", "idx_name2").ForJoin()).Find(&Person{}) // select * from person force index for join (idx_name1,idx_name2)

	// 批量查询
	var persons []*Person
	result = db.FindInBatches(&persons, 100, func(tx *gorm.DB, batch int) error {
		for _, _ = range persons {

		}
		// 返回 err 终止处理
		return nil
	})
	_ = result.RowsAffected // 记录处理数量

	// 查询单个字段
	db.Model(&Person{}).Distinct().Pluck("id", &[]int{})

	// 合并查询条件
	db.Scopes(
		func(tx *gorm.DB) *gorm.DB {
			return tx.Where("id = 1")
		},
		func(tx *gorm.DB) *gorm.DB {
			return tx.Where("name = ?", "zs")
		},
	).Find(&Person{})

	// 计数
	var count int64
	db.Model(&Person{}).Count(&count)

	// 关联表总数
	db.Model(&Person{Id: 1}).Association("Blog").Count()

	// 查询关联表
	db.Model(&Person{Id: 1}).Where("name = ?", "zs").Association("Blog").Find(&Blog{}) // select from blog join person on xxx and id = 1 where name = 'zs'
	db.Preload("Blog").Find(&Person{})
	db.Preload("Blog.User").Find(&Person{})                       // 加载关联表的关联表
	db.Preload(clause.Associations).Find(&Person{})               // 加载所有关联表，但不会加载关联表的关联表
	db.Preload("Blog", "id in ?", []int{1, 2, 3}).Find(&Person{}) // 根据条件查关联表
	db.Preload("Blog", func(db *gorm.DB) *gorm.DB {
		return db.Order("id desc")
	}).Find(&Person{}) // 根据条件查关联表
	db.Preload("NestedBlog.User").Find(&Person{}) // 加载内嵌字段表
	db.Preload("User").Find(&Person{})            // 有 Blog.User 时，加载关联表

	// 运行 explain
	db.Dialector.Explain("")

	db.Find(&Person{}, datatypes.JSONQuery("attributes").HasKey("role"))         // where json_extract(`attributes`, '$.role') is not null
	db.Find(&Person{}, datatypes.JSONQuery("attributes").HasKey("orgs", "orga")) // where json_extract(`attributes`, '$.orgs.orga') is not null

	db.Find(&Person{}, datatypes.JSONQuery("attributes").Equals("zs", "name")) // where json_extract(`attributes`, '$.name') = "zs"
}

func Update() {
	// 更新字段
	db.Model(&Person{}).Update("name", "zs")
	db.Model(&Person{Id: 1}).Update("name", "zs") // 如果主键有值，会作为条件
	db.Model(&Person{}).Updates(map[string]any{}) // 更新多个字段值
	db.Updates(&Person{})                         // 只更新非零值字段。会触发关联表更新

	// 选择更新字段
	db.Select("Name1").Updates(&Person{})
	db.Select("name1").Updates(map[string]any{})
	db.Select("*").Updates(&Person{})             // 选择更新所有字段，包括零值字段
	result := db.Omit("Name1").Updates(&Person{}) // 忽略某字段更新
	_ = result.RowsAffected                       // 受影响的行数
	db.Model(&Person{}).Update("name1", db.Find(&Person{}))

	// 忽略 where 更新所有字段
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&Person{}).Update("name1", "zs")

	// 不触发钩子函数的更新
	db.Model(&Person{}).UpdateColumn("names", "zs")
	db.Model(&Person{}).UpdateColumns(map[string]any{"name1": "zs"})

	// 返回更新后的记录
	var persons []*Person
	db.Model(&persons).Clauses(clause.Returning{}).Where("id = 1").Update("name1", "zs")                                          // update person set name1 = 'zs' where id = 1 returning *
	db.Model(&persons).Clauses(clause.Returning{Columns: []clause.Column{{Name: "name1"}}}).Where("id = 1").Update("name1", "zs") // update person set name1 = 'zs' where id = 1 returning name1
	_ = persons

	// 更新联系的实体
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&Person{})

	// 更新关联表
	db.Model(&Person{}).Association("Blog").Replace(&Blog{})
	db.Model(&Person{}).Association("Blog").Replace(&[]*Blog{})
}

func Delete() {
	// 删除记录
	db.Delete(&Person{})

	// 主键做条件
	db.Delete(&Person{}, 1)
	db.Delete(&Person{}, "1")
	db.Delete(&Person{}, []int{1})
	db.Delete(&Person{Id: 1})
	persons := []*Person{{Id: 1}}
	db.Delete(&persons)

	db.Delete(&Person{}, "name1 = ?", "zs")

	// 返回被删除的记录
	db.Clauses(clause.Returning{}).Where("id = 1").Delete(&persons)
	db.Clauses(clause.Returning{Columns: []clause.Column{{Name: "name1"}}}).Where("id = 1").Delete(&persons)

	// 避免软删除
	db.Unscoped().Delete(&Person{})

	// 删除连接表记录，条件必须有主码
	db.Select(clause.Associations).Delete(&Person{Id: 1})
	db.Select("Blog").Delete(&Person{Id: 1})

	// 删除关联表
	db.Model(&Person{}).Association("Blog").Delete(&Blog{})
	db.Model(&Person{}).Association("Blog").Delete(&[]*Blog{})
	db.Model(&Person{}).Association("Blog").Clear()
}

func RawSQL() {
	// 查询
	db.Raw("select * from person").Scan(&Person{})

	// 更新
	db.Exec("update person set name1 = ? where id = ?", "zs", 1)
}

func SQLString() {
	// 没有实际运行 SQL
	stmt := db.Session(&gorm.Session{DryRun: true}).First(&Person{}, 1).Statement
	stmt.SQL.String() // select * from person where id = $1 order by id
	_ = stmt.Vars     // []any{1}

	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&Person{}).Where("id = ?", 100).Limit(10).Order("age desc").Find(&[]Person{})
	})
	_ = sql
}

func RunOneTCP() {
	db.Connection(func(tx *gorm.DB) error {
		tx.Exec("")

		tx.First(&Person{})
		return nil
	})
}

func Clause() {
	db.ToSQL(func(tx *gorm.DB) *gorm.DB { return db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&Person{}) })

	var limit = 1
	statement := db.Clauses(
		clause.Select{Columns: []clause.Column{{Name: "*"}}},
		clause.From{Tables: []clause.Table{{Name: clause.CurrentTable}}},
		clause.Limit{Limit: &limit},
		clause.OrderBy{Columns: []clause.OrderByColumn{
			{
				Column: clause.Column{
					Table: clause.CurrentTable,
					Name:  clause.PrimaryKey,
				},
			},
		}}).Statement
	statement.Build("SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT", "FOR")
	println(statement.SQL.String()) // SELECT `*` FROM ` ORDER BY `. LIMIT ?
}

func UserContext() {
	tx := db.WithContext(context.Background())
	_ = tx.Statement.Context
}

func HandleError() {
	err := db.Create(&Person{}).Error
	me, ok := err.(*dm.MySQLError)
	if ok {
		_ = me.Number
	}

	// TranslateError=true
	errors.Is(err, gorm.ErrDuplicatedKey)
	errors.Is(err, gorm.ErrForeignKeyViolated)
}

func Session() {
	db.Session(&gorm.Session{
		CreateBatchSize:          1,    // 影响批量添加记录行为
		SkipHooks:                true, // 将跳过执行钩子函数
		NewDB:                    true, // 操作 db.Statement 不影响 SQL 语句构建
		Initialized:              true,
		DisableNestedTransaction: true,                                // 关闭 savepoint, rollback to
		Context:                  context.Background(),                // 等于 db.WithContext()
		Logger:                   logger.Default.LogMode(logger.Info), // 等于 db.Debug()
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
}

func Transaction() {
	db.Transaction(func(tx *gorm.DB) error {

		// 内部事务
		tx.Transaction(func(tx2 *gorm.DB) error {
			return nil
		})

		return nil
	})

	tx := db.Begin()
	tx.SavePoint("sp1")
	tx.RollbackTo("sp1")
	tx.Rollback()
	tx.Commit()
}

func Migration() {
	db.Migrator().CurrentDatabase()
	db.Migrator().CreateTable(&Person{})
	db.Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&Person{})
	db.Migrator().HasTable(&Person{})
	db.Migrator().HasTable("person")
	db.Migrator().DropTable(&Person{})
	db.Migrator().DropTable("person")
	db.Migrator().RenameTable(&Person{}, &Person{})
	db.Migrator().RenameTable("person", "person_infos")
	db.Migrator().AddColumn(&Person{}, "Name")
	db.Migrator().DropColumn(&Person{}, "Name")
	db.Migrator().AlterColumn(&Person{}, "Name")
	db.Migrator().HasColumn(&Person{}, "Name")
	db.Migrator().RenameColumn(&Person{}, "Name", "NewName")
	db.Migrator().ColumnTypes(&Person{})
}

func View() {
	query := db.Model(&Person{}).Where("age > ?", 20)
	db.Migrator().CreateView("users_pets", gorm.ViewOption{Query: query})                                   // CREATE VIEW `users_view` AS SELECT * FROM `users` WHERE age > 20
	db.Migrator().CreateView("users_pets", gorm.ViewOption{Query: query, Replace: true})                    // CREATE OR REPLACE VIEW `users_pets` AS SELECT * FROM `users` WHERE age > 20
	db.Migrator().CreateView("users_pets", gorm.ViewOption{Query: query, CheckOption: "WITH CHECK OPTION"}) // CREATE VIEW `users_pets` AS SELECT * FROM `users` WHERE age > 20 WITH CHECK OPTION
	db.Migrator().DropView("users_pets")                                                                    // DROP VIEW IF EXISTS "users_pets"
}

func Constraint() {
	type UserIndex struct {
		Name string `gorm:"check:name_checker,name <> 'zs'"`
	}
	db.Migrator().CreateConstraint(&UserIndex{}, "name_checker")
	db.Migrator().DropConstraint(&UserIndex{}, "name_checker")
	db.Migrator().HasConstraint(&UserIndex{}, "name_checker")

	type CreditCard struct {
		gorm.Model
		Number string
		UserID uint
	}
	type User struct {
		gorm.Model
		CreditCards []CreditCard
	}
	db.Migrator().CreateConstraint(&User{}, "CreditCards")
	db.Migrator().CreateConstraint(&User{}, "fk_users_credit_cards")

	db.Migrator().HasConstraint(&User{}, "CreditCards")
	db.Migrator().HasConstraint(&User{}, "fk_users_credit_cards")

	db.Migrator().DropConstraint(&User{}, "CreditCards")
	db.Migrator().DropConstraint(&User{}, "fk_users_credit_cards")
}

func Index() {
	type User struct {
		gorm.Model
		Name string `gorm:"size:255;index:idx_name,unique"`
	}

	db.Migrator().CreateIndex(&User{}, "Name")
	db.Migrator().CreateIndex(&User{}, "idx_name")

	db.Migrator().DropIndex(&User{}, "Name")
	db.Migrator().DropIndex(&User{}, "idx_name")

	db.Migrator().HasIndex(&User{}, "Name")
	db.Migrator().HasIndex(&User{}, "idx_name")

	type User1 struct {
		gorm.Model
		Name  string `gorm:"size:255;index:idx_name,unique"`
		Name2 string `gorm:"size:255;index:idx_name_2,unique"`
	}
	db.Migrator().RenameIndex(&User1{}, "Name", "Name2")
	db.Migrator().RenameIndex(&User1{}, "idx_name", "idx_name_2")
}

func Set() {
	tx := db.Set("key", "value")
	tx.Get("key")

	tx = db.InstanceSet("key", "value")
	tx.InstanceGet("key")
}

func DBResolver() {
	db.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:           []gorm.Dialector{mysql.Open("db2_dsn")},                        // 主节点
			Replicas:          []gorm.Dialector{mysql.Open("db3_dsn"), mysql.Open("db4_dsn")}, // 复制节点
			Policy:            dbresolver.RandomPolicy{},                                      // 随机选择
			TraceResolverMode: true,                                                           // 日志打印选择的节点
		}). // 指定表使用哪个复制节点
			Register(dbresolver.Config{
				// 主节点是 db
				Replicas: []gorm.Dialector{mysql.Open("db5_dsn")},
			}, &Person{}).
			Register(dbresolver.Config{
				Sources:  []gorm.Dialector{mysql.Open("db6_dsn"), mysql.Open("db7_dsn")},
				Replicas: []gorm.Dialector{mysql.Open("db8_dsn")},
			}, "person").
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)

	// 指定节点和读写
	db.Clauses(dbresolver.Write).First(&Person{})
	db.Clauses(dbresolver.Use("person")).First(&Person{})
	db.Clauses(dbresolver.Use("person"), dbresolver.Write).First(&Person{})

	// 指定事务运行节点
	tx := db.Clauses(dbresolver.Read).Begin()
	tx = db.Clauses(dbresolver.Write).Begin()
	tx = db.Clauses(dbresolver.Use("person"), dbresolver.Write).Begin()
	_ = tx
}

func Shard() {
	db.Use(sharding.Register(sharding.Config{
		ShardingKey:         "user_id",
		NumberOfShards:      64,
		PrimaryKeyGenerator: sharding.PKSnowflake,
	}, "person", Person{}))
}

func Serializer() {
	schema.RegisterSerializer("json", nil)
}

type SerializerType string

// 出库
func (es *SerializerType) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue any) (err error) {
	return nil
}

// 进库
func (es *SerializerType) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue any) (any, error) {
	return nil, nil
}

func Plugin() {
	db.Use(prometheus.New(prometheus.Config{
		DBName:          "db1",
		RefreshInterval: 15,
		PushAddr:        "prometheus pusher address",
		StartServer:     true,
		HTTPServerPort:  8080,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		},
	}))

	_ = db.Config.Plugins["pluginName"]
}

func Callback() {
	db.Callback().Create().Register("create_callback_name", func(db *gorm.DB) {})
	db.Callback().Create().Remove("create_callback_name")
	db.Callback().Create().Replace("create_callback_name", func(db *gorm.DB) {})
	db.Callback().Create().Before("create_callback_name").Register("create_callback_name_1", func(db *gorm.DB) {})
	db.Callback().Create().After("create_callback_name").Register("create_callback_name_1", func(db *gorm.DB) {})
	db.Callback().Create().After("*").Register("create_callback_name_1", func(db *gorm.DB) {})
}

func Generate() {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "query",
		ModelPkgPath:      "model",
		WithUnitTest:      true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(db)

	g.ApplyBasic(
		// 根据 person 表生成结构体
		g.GenerateModel("person"),
		// 根据 person 表生成结构体 Person
		g.GenerateModelAs("person", "Person"),
		// 修改生成结构体信息
		g.GenerateModel("person", gen.FieldIgnore("address"), gen.FieldType("id", "int64"), gen.FieldType("tags", "datatypes.JSON")),
		// 给生成结构体附加方法
		g.GenerateModel("person", gen.WithMethod((&Person{}).TableName)),
		g.GenerateModel("person", gen.WithMethod(Person{})),
	)

	// 给所有生成的结构体附加方法
	g.WithOpts(gen.WithMethod((&Person{}).TableName))

	// 给所有生成的结构体，附加默认表名方法
	g.WithOpts(gen.WithMethod(gen.DefaultMethodTableWithNamer))

	// 生成所有表的结构体
	g.ApplyBasic(g.GenerateAllTable()...)

	g.WithDataTypeMap(map[string]func(gorm.ColumnType) (dataType string){
		"int": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "*int32"
			}
			return "int32"
		},

		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			ct, _ := columnType.ColumnType()
			if strings.HasPrefix(ct, "tinyint(1)") {
				return "bool"
			}
			return "byte"
		},
	})

	g.Execute()
	// query.SetDefault(db)

	gormdb, _ := gorm.Open(rawsql.New(rawsql.Config{
		// SQL:      rawsql, // 建表 SQL
		FilePath: []string{
			"./sql",        // 建表 SQL 文件夹
			"./deploy.sql", // 建表 SQL 文件
		},
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
	})
	_ = gormdb

	// go install gorm.io/gen/tools/gentool@latest
	// gentool.exe -db mysql -dsn "root:123456@tcp(ivfzhou-debian:3306)/db_certs_local?charset=utf8mb4&parseTime=True&loc=Local" -fieldWithIndexTag -fieldWithTypeTag -modelPkgName model -outPath query -withUnitTest
}

func GenerateDynamic() {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "query",
		ModelPkgPath:      "model",
		WithUnitTest:      true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(db)

	// 动态语句
	type Querier interface {
		// where("name = @name AND age = @age")
		Query1(name string, age int) (gen.M, error)

		// SELECT * FROM @@table WHERE id = @id
		Query2(id int) (*gen.T, error)

		// SELECT * FROM @@table WHERE id = @id
		Query3(id int) (gen.M, error)

		// INSERT INTO @@table (name, age) VALUES (@name, @age)
		Insert(name string, age int) (gen.RowsAffected, error)

		// SELECT * FROM @@table WHERE @@column = @value
		Query4(column string, value string) (*gen.T, error)

		// SELECT * FROM users WHERE
		//  {{if name != ""}}
		//      username = @name AND
		//  {{end}}
		//  role = "admin"
		Query5(name string) (*gen.T, error)

		// SELECT * FROM users
		//  {{if user != nil}}
		//      {{if user.ID > 0}}
		//          WHERE id = @user.ID
		//      {{else if user.Name != ""}}
		//          WHERE username = @user.Name
		//      {{end}}
		//  {{end}}
		Query6(user *gen.T) (gen.T, error)

		// SELECT * FROM @@table
		//  {{where}}
		//      id = @id
		//  {{end}}
		Query7(id int) gen.T

		// SELECT * FROM @@table
		//  {{where}}
		//    {{if !start.IsZero()}}
		//      created_time > @start
		//    {{end}}
		//    {{if !end.IsZero()}}
		//      AND created_time < @end
		//    {{end}}
		//  {{end}}
		Query8(start, end time.Time) ([]gen.T, error)

		// UPDATE @@table
		//  {{set}}
		//    {{if user.Name != ""}} username = @user.Name, {{end}}
		//    {{if user.Age > 0}} age = @user.Age, {{end}}
		//    {{if user.Age >= 18}} is_adult = 1 {{else}} is_adult = 0 {{end}}
		//  {{end}}
		// WHERE id = @id
		Update(user gen.T, id int) (gen.RowsAffected, error)

		// SELECT * FROM @@table
		// {{where}}
		//   {{for _, user := range users}}
		//     {{if user.Name !="" && user.Age > 0}}
		//       (username = @user.Name AND age = @user.Age AND role LIKE concat("%", @user.Role, "%")) OR
		//     {{end}}
		//   {{end}}
		// {{end}}
		Query9(users []*gen.T) ([]*gen.T, error)
	}
	g.ApplyInterface(func(Querier) {}, Person{}, g.GenerateModel("person"))

	g.Execute()
}
