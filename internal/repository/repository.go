package repository

import (
	"github.com/Xapsiel/EffectiveMobile/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Song interface {
	GetSongs(filter model.Song, page int, limit int) ([]model.Song, error)
	GetSongVerse(song model.Song, verse int) (string, int, error)
	DeleteSong(song model.Song) (bool, error)
	UpdateSong(song_name, group_name string, song model.Song) (bool, model.Song, error)
	Add(song model.Song) (int, error)
}
type Repository struct {
	Song
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{
		Song: NewSongRepository(db),
	}
}
