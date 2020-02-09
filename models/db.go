package models

type Categories struct {
	ID    uint   `sql:"id" gorm:"primary_key;auto_increment:true"`
	Title string `sql:"title" gorm:"not null"`
}

type SongSkins struct {
	ID    uint    `sql:"id" gorm:"primary_key;auto_increment:true"`
	Name  string  `sql:"name" gorm:"not null"`
	Song  *string `sql:"song"`
	Stage *string `sql:"stage"`
	Don   *string `sql:"don"`
}

type Makers struct {
	ID   uint    `sql:"maker_id" gorm:"primary_key;auto_increment:true"`
	Name string  `sql:"name" gorm:"not null"`
	Url  *string `sql:"url"`
}

type Songs struct {
	ID           uint     `sql:"id" gorm:"primary_key;auto_increment:true"`
	Title        string   `sql:"text" gorm:"not null"`
	TitleLang    *string  `sql:"title_lang"`
	Subtitle     *string  `sql:"subtitle"`
	SubtitleLang *string  `sql:"subtitle_lang"`
	Easy         *string  `sql:"easy"`
	Normal       *string  `sql:"normal"`
	Hard         *string  `sql:"hard"`
	Oni          *string  `sql:"oni"`
	Ura          *string  `sql:"ura"`
	Enabled      int      `sql:"enabled" gorm:"not null"`
	Category     *int     `sql:"category"`
	Type         *string  `sql:"type"`
	Offset       *float32 `sql:"offset"`
	SkinId       *int     `sql:"skin_id"`
	Preview      *float32 `sql:"preview"`
	Volume       *float32 `sql:"volume"`
	MakerId      *int     `sql:"maker_id"`
}

type SongRes struct {
	Songs
	CategoryName string
}
