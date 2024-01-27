package internal

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mic-trainning-lesson-part2/model"
	"os"
	"time"
)

// 初始化数据库连接
// 注：大写是为了让其他包也能调用
var DB *gorm.DB
var err error

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbName"`
	UserName string `mapstructure:"userName"`
	Password string `mapstructure:"password"`
}

func InitDB() {
	//	配置gorm日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	//	声明数据库连接信息
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConf.DBConfig.UserName, AppConf.DBConfig.Password, "localhost",
		AppConf.DBConfig.Port, AppConf.DBConfig.DBName)
	zap.S().Infof(conn)
	//	连接数据库
	DB, err = gorm.Open(mysql.Open(conn), &gorm.Config{
		//	配置gorm选项
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}
	fmt.Println("Mysql初始化完成...")
	//	根据结构体自动生成表
	err = DB.AutoMigrate(&model.Category{}, &model.Brand{}, &model.Advertise{}, &model.Product{})
	if err != nil {
		fmt.Println("自动生成失败：" + err.Error())
	}
}

// 根据分页查询所定义的一个分页函数模板
func MyPaging(pageNo, pageSize int) func(db *gorm.DB) *gorm.DB {
	//	这个模板函数是用于返回分页的条件的
	return func(db *gorm.DB) *gorm.DB {
		//	判断分页条件
		if pageNo < 1 {
			pageNo = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize < 0:
			pageSize = 5
		}
		offset := (pageNo - 1) * pageSize
		//	声明好分页条件
		return db.Offset(offset).Limit(pageSize)
	}
}
