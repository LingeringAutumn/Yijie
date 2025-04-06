package client

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/LingeringAutumn/Yijie/pkg/constants"
	"github.com/LingeringAutumn/Yijie/pkg/errno"
	"github.com/LingeringAutumn/Yijie/pkg/logger"
	"github.com/LingeringAutumn/Yijie/pkg/utils"
)

// InitMySQL 通用初始化mysql函数，该函数用于初始化与 MySQL 数据库的连接，并进行相关配置。
// 最终返回一个 *gorm.DB 类型的数据库连接对象和可能出现的错误。
func InitMySQL() (db *gorm.DB, err error) {
	// 调用 utils 包中的 GetMysqlDSN 函数，获取 MySQL 数据库的 DSN（Data Source Name）。
	// DSN 是一个包含数据库连接所需信息的字符串，如用户名、密码、数据库地址等。
	dsn, err := utils.GetMysqlDSN()
	// 检查获取 DSN 时是否出错。
	if err != nil {
		// 如果出错，使用 errno 包的 NewErrNo 函数创建一个自定义错误对象。
		// 这里使用了预定义的错误码 errno.InternalDatabaseErrorCode，并添加了错误信息。
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.InitMySQL get mysql DSN error: %v", err))
	}

	// 使用 gorm 库的 Open 函数打开与 MySQL 数据库的连接。
	// mysql.Open(dsn) 是 gorm 提供的用于 MySQL 数据库的驱动，传入之前获取的 DSN。
	// 第二个参数是一个 &gorm.Config 类型的指针，用于配置 gorm 的行为。
	db, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			// PrepareStmt: true 表示在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，
			// 这样可以提高后续执行相同 SQL 语句的效率。
			PrepareStmt: true,
			// SkipDefaultTransaction: false 表示不禁用默认事务，
			// 即单个创建、更新、删除操作时会使用事务。
			SkipDefaultTransaction: false,
			// TranslateError: true 允许 gorm 翻译错误信息，方便开发者理解。
			TranslateError: true,
			// NamingStrategy 用于配置表名和字段名的命名策略。
			// SingularTable: true 表示使用单数表名。
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			// Logger 用于配置日志记录。
			// glogger.New 是自定义的日志记录器创建函数，传入一个基础的日志记录器和配置。
			Logger: glogger.New(
				logger.GetMysqlLogger(),
				glogger.Config{
					// SlowThreshold: time.Second 表示超过一秒的查询被认为是慢查询。
					SlowThreshold: time.Second,
					// LogLevel: glogger.Warn 表示日志等级为警告，只记录警告及以上级别的日志。
					LogLevel: glogger.Warn,
					// IgnoreRecordNotFoundError: true 表示当未找到记录（RecordNotFoundError）时不记录日志。
					IgnoreRecordNotFoundError: true,
					// ParameterizedQueries: true 表示在 SQL 中不包含参数，提高安全性。
					ParameterizedQueries: true,
					// Colorful: false 表示禁用颜色渲染，使日志输出更简洁。
					Colorful: false,
				}),
		})
	// 检查打开数据库连接时是否出错。
	if err != nil {
		// 如果出错，同样使用 errno 包的 NewErrNo 函数创建自定义错误对象。
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.InitMySQL mysql connect error: %v", err))
	}

	// 尝试从 gorm.DB 对象中获取底层的 *sql.DB 实例对象。
	// 这个对象可以用于进行更底层的数据库连接配置。
	sqlDB, err := db.DB()
	// 检查获取底层数据库对象时是否出错。
	if err != nil {
		// 如果出错，使用 errno 包的 NewErrNo 函数创建自定义错误对象。
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("get generic database object error: %v", err))
	}

	// 设置数据库连接池的最大闲置连接数。
	// constants.MaxIdleConns 是一个预定义的常量，表示最大闲置连接数。
	sqlDB.SetMaxIdleConns(constants.MaxIdleConns)
	// 设置数据库连接池的最大连接数。
	// constants.MaxConnections 是一个预定义的常量，表示最大连接数。
	sqlDB.SetMaxOpenConns(constants.MaxConnections)
	// 设置数据库连接的最大可复用时间。
	// constants.ConnMaxLifetime 是一个预定义的常量，表示最大可复用时间。
	sqlDB.SetConnMaxLifetime(constants.ConnMaxLifetime)
	// 设置数据库连接最长保持空闲状态的时间。
	// constants.ConnMaxIdleTime 是一个预定义的常量，表示最长保持空闲状态时间。
	sqlDB.SetConnMaxIdleTime(constants.ConnMaxIdleTime)
	// 为 gorm.DB 对象设置上下文，这里使用 context.Background() 作为基础上下文。
	db = db.WithContext(context.Background())

	// 进行数据库连通性测试，使用 sqlDB.Ping() 方法尝试与数据库建立连接。
	if err = sqlDB.Ping(); err != nil {
		// 如果测试失败，使用 errno 包的 NewErrNo 函数创建自定义错误对象。
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("ping database error: %v", err))
	}

	// 如果一切正常，返回初始化好的 gorm.DB 对象和 nil 错误。
	return db, nil
}
