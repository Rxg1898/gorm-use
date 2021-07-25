package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 定义表
type NewUser struct {
	ID           uint
	Name         string  `gorm:"column:name"`
	Email        *string // 也可以解决空值问题
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Deleted      gorm.DeletedAt
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
	db.AutoMigrate(&NewUser{})

	users := []NewUser{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
	db.CreateInBatches(users, 2)
	for _, user := range users {
		fmt.Println(user.ID)
	}

	// 删除对象需要指定主键，否则会触发批量Delete
	db.Delete(&NewUser{}, 2) //没有软删除列字段 DELETE FROM `users` WHERE `users`.`id` = 10

	// 软删除
	db.Delete(&NewUser{}, 2) //UPDATE `new_users` SET `deleted`='2021-07-25 00:42:55.479' WHERE `new_users`.`id` = 2 AND `new_users`.`deleted` IS NULL

	// 硬删除
	db.Unscoped().Delete(&NewUser{ID: 2}) // DELETE FROM `new_users` WHERE `new_users`.`id` = 2
}
