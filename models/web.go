package models

type Song struct {
	ID           uint       `json:"id"`
	Title        string     `json:"title"`
	TitleLang    *string    `json:"title_lang"`
	Subtitle     *string    `json:"subtitle"`
	SubtitleLang *string    `json:"subtitle_lang"`
	Stars        []*string  `json:"stars"`
	Preview      *float32   `json:"preview"`
	Category     string     `json:"category"`
	Type         *string    `json:"type"`
	Offset       *float32   `json:"offset"`
	SongSkin     *SongSkins `json:"song_skin"`
	Volume       *float32   `json:"volume"`
	Maker        *Makers    `json:"maker"`
}

type Message struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
