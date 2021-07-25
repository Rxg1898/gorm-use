package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 有多张 CreditCard，UserID 是外键
type User struct {
	gorm.Model
	CreditCards []CreditCard `gorm:"foreignKey:UserRefer"` // 并没有生成外键约束，大型系统不建议使用外键约束，业务层保证数据的一致性。外键约束可以保证数据的完整性。
}

type CreditCard struct {
	gorm.Model
	Number    string
	UserRefer uint
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:root@tcp(127.0.0.1:3306)/gormtest?charset=utf8mb4&parseTime=True&loc=Local"

	// 设置全局的logger。打印每次的sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	// 建表
	// db.AutoMigrate(&User{})
	// db.AutoMigrate(&CreditCard{})

	// 插入数据
	// var user User
	// db.Create(&user)
	// db.Create(&CreditCard{
	// 	Number: "12",
	// 	UserRefer: user.ID,
	// })
	// db.Create(&CreditCard{
	// 	Number: "34",
	// 	UserRefer: user.ID,
	// })

	// 查询
	var user User
	db.Preload("CreditCards").First(&user) // Preload表名
	for _, card := range user.CreditCards {
		fmt.Println(card.Number)
	}
}
