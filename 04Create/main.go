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

	// 迁移自动建表
	_ = db.AutoMigrate(&User{})
	user := User{
		Name: "an",
	}

	// Create
	result := db.Create(&user)
	fmt.Println(user.ID)             // 返回数据主键'
	fmt.Println(result.Error)        // 错误
	fmt.Println(result.RowsAffected) // 插入多少行

	// Update 会更新零值，Updates不会
	db.Model(&User{ID: 2}).Update("Name", "")      // UPDATE `users` SET `name`='',`updated_at`='2021-07-24 22:59:07.994' WHERE `id` = 2
	db.Model(&User{ID: 2}).Updates(User{Name: ""}) // UPDATE `users` SET `updated_at`='2021-07-24 22:59:08.167' WHERE `id` = 2
	empty := ""
	db.Model(&User{ID: 2}).Updates(User{Email: &empty}) // UPDATE `users` SET `email`='',`updated_at`='2021-07-24 23:01:22.732' WHERE `id` = 2

	// 批量插入
	var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
	db.Create(&users)
	for _, user = range users {
		fmt.Println(user.ID)
	}

	// 指定批量创建的数量，最多一次插入2条数据
	// 为什么不一次性插入所有？sql语句有长度限制!
	users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
	db.CreateInBatches(users, 2)
	for _, user = range users {
		fmt.Println(user.ID)
	}

	// 创建钩子,角色控制admin才可以插入数据

	// map创建
	db.Model(&User{}).Create(map[string]interface{}{ // INSERT INTO `users` (`age`,`name`) VALUES (18,'jinzhu')
		"Name": "jinzhu",
		"Age":  18,
	})

	// 关联创建
	

}
