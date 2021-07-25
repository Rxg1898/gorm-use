package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 定义表
type User struct {
	ID           uint
	Name         string
	Email        *string // 也可以解决空值问题
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

	// first查询 主键升序查询第一条
	var user User
	db.First(&user)
	fmt.Println(user.ID) // SELECT * FROM `users` ORDER BY `users`.`id` LIMIT 1

	// take 没有排序

	// last 主键降序排

	// 通过主键查询
	var u User
	result := db.First(&u, 10)                           // db.First(&u, "10") 主键会转换成int类型
	if errors.Is(result.Error, gorm.ErrRecordNotFound) { // 准确判断没有找到数据错误
		fmt.Println("未找到")
	}
	fmt.Println(user.ID)

	// 检索全部对象
	var users []User
	result = db.Find(&users)
	fmt.Println("总共记录:", result.RowsAffected)
	for _, user = range users {
		fmt.Println(user.ID)
	}
}
