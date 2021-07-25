package main

import (
	"database/sql"
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
	Name         string  `gorm:"column:name"`
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

	// Save会保存所有字段，即使是零值
	var user User
	db.First(&user)
	user.Age = 100
	user.Name = ""
	db.Save(&user) // UPDATE `users` SET `name`='',`email`=NULL,`age`=100,`birthday`=NULL,`member_number`=NULL,`activated_at`=NULL,`created_at`='2021-07-24 23:26:02.481',`updated_at`='2021-07-25 00:24:29.073' WHERE `id` = 1

	db.First(&user)
	user.Age = 100
	user.Name = ""
	user.ID = 0    // ID为0则插入
	db.Save(&user) // INSERT INTO `users` (`name`,`email`,`age`,`birthday`,`member_number`,`activated_at`,`created_at`,`updated_at`) VALUES ('',NULL,100,NULL,NULL,NULL,'2021-07-24 23:26:02.481','2021-07-25 00:27:23.808')

	// 通过update方法更新
	db.Model(&User{}).Where("name = ?", "an").Update("name", "hello")

	// 更新选定字段Select 忽略字段Omit
	// 使用 Map 进行 Select
	db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18})
	// UPDATE users SET name='hello' WHERE id=111;

	db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18})
	// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

	
}
