package model

type Song struct {
	ID          *int    `json:"id,omitempty" example:"1" db:"id"`
	GroupId     *int    `json:"-"  db:"group_id"`
	Group       *string `json:"group_name,omitempty" example:"Muse" db:"group_name"`
	SongName    *string `json:"song_name,omitempty" example:"Supermassive Black Hole" db:"song_name"`
	ReleaseDate *string `json:"releaseDate,omitempty" example:"19.07.2006" db:"release_date"`
	Link        *string `json:"link,omitempty" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw" db:"link"`
	Text        *string `json:"text,omitempty" example:"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?" db:"text"`
}
