package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"taiko-web/config"
	"taiko-web/models"
)

var db *gorm.DB

func Init(conf *config.TaikoWebConfig) {
	var err error
	db, err = gorm.Open(conf.DB.Dialect, conf.DB.ConnString)
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(
		&models.Categories{},
		&models.SongSkins{},
		&models.Makers{},
		&models.Songs{},
	)

	// Init Categories
	db.Table("categories").
		FirstOrCreate(
			&models.Categories{}, models.Categories{Title: "J-POP"}).
		FirstOrCreate(
			&models.Categories{}, models.Categories{Title: "アニメ"}).
		FirstOrCreate(
			&models.Categories{}, models.Categories{Title: "ボーカロイド™曲"}).
		FirstOrCreate(
			&models.Categories{}, models.Categories{Title: "バラエティ"}).
		FirstOrCreate(
			&models.Categories{}, models.Categories{Title: "クラシック"}).
		FirstOrCreate(
			&models.Categories{}, models.Categories{Title: "ゲームミュージック"}).
		FirstOrCreate(
			&models.Categories{}, models.Categories{Title: "ナムコオリジナル"})
}

func GetDB() *gorm.DB {
	return db
}
