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

// User 拥有并属于多种 language，`user_languages` 是连接表
type User struct {
	gorm.Model
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	gorm.Model
	Name string
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
	db.AutoMigrate(&User{})

	// 插入数据
	// languages := []Language{}
	// languages = append(languages, Language{Name: "go"})
	// languages = append(languages, Language{Name: "python"})
	// user := User{
	// 	Languages: languages,
	// }

	// db.Create(&user)

	// 查询
	// var user User
	// db.Preload("Languages").First(&user)
	// for _, language := range user.Languages {
	// 	fmt.Println(language.Name)
	// }

	// 取出用户，根据用户查询
	var user User
	db.First(&user)
	var languages []Language
	_ = db.Model(&user).Association("Languages").Find(&languages)
	for _, language := range languages {
		fmt.Println(language.Name)
	}
}
