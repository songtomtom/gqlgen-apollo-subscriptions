package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	Post struct {
		gorm.Model
	}

	Comment struct {
		gorm.Model
		PostID  uint
		Content string
	}
)

func open(username, password, host, port, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbName,
	)

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	if err = db.AutoMigrate(&Post{}, &Comment{}); err != nil {
		return nil, fmt.Errorf("failed to auto migration schema: %v", err)
	}

	post := Post{}
	db.Create(&post)

	comment := Comment{
		PostID:  post.ID,
		Content: "test_content",
	}
	db.Create(&comment)

	return db, nil

}
