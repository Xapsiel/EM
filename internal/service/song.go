package service

import (
	"fmt"

	"github.com/Xapsiel/EffectiveMobile/internal/api"
	"github.com/Xapsiel/EffectiveMobile/internal/model"
	"github.com/Xapsiel/EffectiveMobile/internal/repository"
)

type songService struct {
	api  api.Client
	repo repository.Song
}

func NewSongService(repo repository.Song) *songService {
	return &songService{repo: repo}
}

func (s *songService) GetSongs(filter model.Song, page int, limit int) ([]model.Song, error) {
	return s.repo.GetSongs(filter, page, limit)

}
func (s *songService) Add(song string, group string) (int, error) {
	if song == "" || group == "" {
		return -1, fmt.Errorf("invalid params")
	}
	res, err := s.api.GetInfo(group, song)
	if err != nil {
		return -1, err
	}
	res.SongName = res.SongName
	res.Group = res.Group
	res.Link = res.Link
	res.ReleaseDate = res.ReleaseDate
	res.Text = res.Text

	//if err != nil {
	//	return 0, err
	//}

	return s.repo.Add(res)
}

func (s *songService) GetSongVerse(song model.Song, verse int) (string, int, error) {
	if song.SongName == nil || song.Group == nil {
		return "", -1, fmt.Errorf("song name or group is empty")
	}
	return s.repo.GetSongVerse(song, verse)
}
func (s *songService) DeleteSong(song model.Song) (bool, error) {
	if song.SongName == nil || song.Group == nil {
		return false, fmt.Errorf("song name or group is empty")
	}
	return s.repo.DeleteSong(song)
}
func (s *songService) UpdateSong(song_name, group_name string, song model.Song) (bool, model.Song, error) {
	if song_name == "" && group_name == "" {
		return false, model.Song{}, nil
	}
	return s.repo.UpdateSong(song_name, group_name, song)
}
