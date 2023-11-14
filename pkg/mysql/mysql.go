package mysql

import (
	"gin-api/pkg/config"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	once     sync.Once
	instance *gorm.DB
)

func New() *gorm.DB {
	dbConfig := config.GetConfig().Database

	dsn := dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + ")/" + dbConfig.Database + "?charset=" + dbConfig.Charset + "&parseTime=True&loc=Local"

	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,            // 禁用自动创建外键约束
		Logger:                                   getGormLogger(), // 使用自定义 Logger
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})

	if err != nil {
		panic("connection failed!")
	}

	return db
}

// 单例模式调用
func GetDB() *gorm.DB {
	once.Do(func() {
		instance = New()
	})

	return instance
}

func getGormLogger() logger.Interface {
	logMode := logger.Info

	return logger.New(getGormLogWriter(), logger.Config{
		SlowThreshold:             1 * time.Second, // 慢 SQL 阈值
		LogLevel:                  logMode,         // 日志级别
		IgnoreRecordNotFoundError: false,           // 忽略ErrRecordNotFound（记录未找到）错误
	})
}

// 自定义 gorm Writer
func getGormLogWriter() logger.Writer {
	var writer io.Writer
	conf := config.GetConfig()

	// 是否启用日志文件
	if conf.Database.EnableFileLogWriter {
		// 自定义 Writer
		writer = &lumberjack.Logger{
			Filename:   conf.Log.RootDir + "/" + conf.Database.LogFilename,
			MaxSize:    conf.Log.MaxSize,
			MaxBackups: conf.Log.MaxBackups,
			MaxAge:     conf.Log.MaxAge,
			Compress:   conf.Log.Compress,
		}
	} else {
		// 默认 Writer
		writer = os.Stdout
	}
	return log.New(writer, "\r", log.LstdFlags)
}
