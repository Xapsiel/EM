package service

import (
	"github.com/Xapsiel/EffectiveMobile/internal/model"
	"github.com/Xapsiel/EffectiveMobile/internal/repository"
)

type Service struct {
	Song
}

type Song interface {
	GetSongs(filter model.Song, page int, limit int) ([]model.Song, error)
	GetSongVerse(song model.Song, verse int) (string, int, error)
	DeleteSong(song model.Song) (bool, error)
	UpdateSong(song_name, group_name string, song model.Song) (bool, model.Song, error)
	Add(song string, group string) (int, error)
}

func NewService(repo repository.Repository) Service {
	return Service{
		Song: NewSongService(repo),
	}

}
