package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taiko-web/db"
	"taiko-web/models"
)

func SongsHandler(c *gin.Context) {
	db := db.GetDB()

	var songRes []models.SongRes
	db.Table("songs").
		Select("songs.*, categories.title as category_name").
		Joins("left join categories on songs.category = categories.id").
		Where("enabled = 1").
		Find(&songRes)

	songs := []models.Song{}
	for _, song := range songRes {
		var maker *models.Makers
		if song.MakerId != nil {
			db.Table("makers").Where("id = ?", song.MakerId).Find(&maker)
		}
		var songSkin *models.SongSkins
		if song.SkinId != nil {
			db.Table("song_skins").Where("id = ?", song.SkinId).Find(&songSkin)
		}
		songs = append(songs, models.Song{
			song.ID,
			song.Title,
			song.TitleLang,
			song.Subtitle,
			song.SubtitleLang,
			[]*string{song.Easy, song.Normal, song.Hard, song.Oni, song.Ura},
			song.Preview,
			song.CategoryName,
			song.Type,
			song.Offset,
			songSkin,
			song.Volume,
			maker,
		})
	}

	c.JSON(http.StatusOK, songs)
}
