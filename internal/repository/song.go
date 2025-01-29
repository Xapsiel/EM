package repository

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Xapsiel/EffectiveMobile/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type songRepository struct {
	db *pgxpool.Pool
}

func NewSongRepository(db *pgxpool.Pool) *songRepository {
	return &songRepository{
		db: db,
	}
}

func (r *songRepository) GetSongs(filter model.Song, page int, limit int) ([]model.Song, error) {
	slog.Info("Начало выполнения GetSongs", "page", page, "limit", limit)

	offset := (page - 1) * limit
	query := `SELECT 
			s.id, g.name, s.song_name, 
			s.release_date, s.link, s.text
			FROM songs as s
			JOIN public.groups g on g.id = s.group_id
			WHERE 1=1`

	var args []interface{}
	argIndex := 1

	if *filter.SongName != "" {
		query += fmt.Sprintf(" AND s.song_name LIKE $%d", argIndex)
		args = append(args, "%"+*filter.SongName+"%")
		argIndex++
	}
	if *filter.Group != "" {
		query += fmt.Sprintf(" AND g.name LIKE $%d", argIndex)
		args = append(args, "%"+*filter.Group+"%")
		argIndex++
	}
	if *filter.Text != "" {
		query += fmt.Sprintf(" AND s.text LIKE $%d", argIndex)
		args = append(args, "%"+*filter.Text+"%")
		argIndex++
	}
	if *filter.Link != "" {
		query += fmt.Sprintf(" AND s.link LIKE $%d", argIndex)
		args = append(args, "%"+*filter.Link+"%")
		argIndex++
	}
	if !(*filter.ReleaseDate == "01.01.0001") {
		date, err := time.Parse("02.01.2006", *filter.ReleaseDate)
		if err == nil || *filter.ReleaseDate == "" {
			query += fmt.Sprintf(" AND s.release_date > $%d", argIndex)
			args = append(args, date)
			argIndex++
		} else {
			slog.Warn(err.Error())
		}

	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	slog.Debug("Сформированный SQL-запрос", "query", query, "args", args)

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		slog.Error("Ошибка при выполнении запроса", "error", err)
		return nil, err
	}
	defer rows.Close()

	var songs []model.Song
	for rows.Next() {
		var song model.Song
		var group, songName, link, text string
		var releaseDate time.Time
		var id int
		if err := rows.Scan(&id, &group, &songName, &releaseDate, &link, &text); err != nil {
			slog.Error("Ошибка при сканировании строки", "error", err)
			return nil, err
		}
		song.ID = &id
		song.Group = &group
		song.SongName = &songName
		var tmp = (fmt.Sprintf("%02d.%02d.%d", releaseDate.Day(), releaseDate.Month(), releaseDate.Year()))
		song.ReleaseDate = &tmp
		song.Link = &link
		song.Text = &text
		songs = append(songs, song)
	}

	slog.Info("Успешно получены песни", "количество песен", len(songs))
	return songs, nil
}

func (r *songRepository) GetSongVerse(song model.Song, verse int) (string, int, error) {
	slog.Info("Начало выполнения GetSongVerse", "song", song, "verse", verse)

	var result struct {
		Text string `json:"text" db:"text"`
		ID   int    `json:"id" db:"id"`
	}

	query := `SELECT s.text, s.id 
			  FROM songs AS s
			  JOIN groups g ON g.id = s.group_id
			  WHERE s.song_name = $1 AND g.name = $2`

	slog.Debug("Сформированный SQL-запрос", "query", query, "song_name", *song.SongName, "group", *song.Group)

	row := r.db.QueryRow(context.Background(), query, song.SongName, song.Group)
	err := row.Scan(&result.Text, &result.ID)
	if err != nil {
		slog.Error("Ошибка при выполнении запроса", "error", err)
		return "", 0, err
	}

	verses := strings.Split(result.Text, "\n\n")
	if verse < 1 || verse > len(verses) {
		slog.Error("Куплет не найден", "verse", verse, "total_verses", len(verses))
		return "", 0, fmt.Errorf("Куплет %d не найден", verse)
	}

	slog.Info("Успешно получен куплет", "id", result.ID, "verse", verse)
	return verses[verse-1], result.ID, nil
}

func (r *songRepository) UpdateSong(song_name, group_name string, song model.Song) (bool, model.Song, error) {
	slog.Info("Начало выполнения UpdateSong", "song_name", song_name, "group_name", group_name)

	var args []interface{}
	argIndex := 1
	query := `UPDATE songs SET `
	setClauses := []string{}
	if song.SongName != nil {
		setClauses = append(setClauses, fmt.Sprintf("song_name = $%d", argIndex))
		args = append(args, *song.SongName)
		argIndex++
	}
	if song.Text != nil {
		setClauses = append(setClauses, fmt.Sprintf("text = $%d", argIndex))
		args = append(args, *song.Text)
		argIndex++
	}
	if song.Link != nil {
		setClauses = append(setClauses, fmt.Sprintf("link = $%d", argIndex))
		args = append(args, *song.Link)
		argIndex++
	}
	if song.ReleaseDate != nil {
		date, err := time.Parse("02.01.2006", *song.ReleaseDate)
		if err == nil || *song.ReleaseDate == "" {
			setClauses = append(setClauses, fmt.Sprintf("release_date = $%d", argIndex))
			args = append(args, date)
			argIndex++
		} else {
			slog.Warn(err.Error())
		}

	}

	if song.Group != nil {
		g, err := r.selectGroup(*song.Group)
		if err != nil {
			slog.Error("Ошибка при выборе группы", "error", err)
			return false, song, err
		}

		setClauses = append(setClauses, fmt.Sprintf("group_id = $%d ", argIndex))
		args = append(args, g.ID)
		argIndex++
	}

	if len(setClauses) == 0 {
		slog.Error("Нет данных для обновления")
		return false, song, fmt.Errorf("нет данных для обновления")
	}

	query += strings.Join(setClauses, ", ") + fmt.Sprintf(" WHERE song_name = $%d AND EXISTS (SELECT 1 FROM groups WHERE name = $%d)", argIndex, argIndex+1)
	args = append(args, song_name, group_name)

	slog.Debug("Сформированный SQL-запрос", "query", query, "args", args)

	_, err := r.db.Exec(context.Background(), query, args...)
	if err != nil {
		slog.Error("Ошибка при обновлении песни", "error", err)
		return false, song, fmt.Errorf("ошибка обновления песни: %w", err)
	}

	slog.Info("Песня успешно обновлена", "song_name", song_name, "group_name", group_name)
	return true, song, nil
}

func (r *songRepository) Add(song model.Song) (int, error) {
	slog.Info("Начало выполнения Add", "song name", *song.SongName, "group name", *song.Group)

	group, err := r.selectGroup(*song.Group)
	if err != nil {
		slog.Error("Ошибка при выборе группы", "error", err)
		return 0, err
	}

	date, err := time.Parse("02.01.2006", *song.ReleaseDate)
	if err != nil {
		slog.Error("Ошибка при парсинге даты", "error", err)
		return 0, err
	}

	query := `INSERT INTO songs (group_id, song_name,release_date,text,link)
    		VALUES ($1, $2, $3, $4, $5) RETURNING id`

	slog.Debug("Сформированный SQL-запрос", "query", query, "args", []interface{}{*group.ID, *song.SongName, date, *song.Text, *song.Link})

	row := r.db.QueryRow(context.Background(), query, group.ID, song.SongName, date, song.Text, song.Link)

	var id int
	err = row.Scan(&id)
	if err != nil {
		slog.Error("Ошибка при добавлении песни", "error", err)
		return 0, err
	}

	slog.Info("Песня успешно добавлена", "id", id)
	return id, nil
}
func (r *songRepository) DeleteSong(song model.Song) (bool, error) {
	slog.Info("Начало выполнения DeleteSong", "song name", *song.SongName, "group name", *song.Group)

	group, err := r.selectGroup(*song.Group)
	if err != nil {
		slog.Error("Ошибка при выборе группы", "error", err)
		return false, err
	}

	query := `DELETE FROM songs WHERE group_id = $1 AND song_name = $2`
	slog.Debug("Сформированный SQL-запрос", "query", query, "args", []interface{}{group.ID, song.SongName})

	_, err = r.db.Exec(context.Background(), query, group.ID, song.SongName)
	if err != nil {
		slog.Error("Ошибка при удалении песни", "error", err)
		return false, err
	}

	slog.Info("Песня успешно удалена", "song_name", *song.SongName, "group_name", *song.Group)
	return true, nil
}

func (r *songRepository) selectGroup(groupName string) (model.Group, error) {
	slog.Info("Начало выполнения selectGroup", "groupName", groupName)

	query := `SELECT id,name from groups WHERE name = $1`
	slog.Debug("Сформированный SQL-запрос", "query", query, "args", []interface{}{groupName})

	row := r.db.QueryRow(context.Background(), query, groupName)
	var group model.Group
	err := row.Scan(&group.ID, &group.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Info("Группа не найдена, создание новой группы", "groupName", groupName)

			query = `INSERT INTO groups(name) VALUES ($1) RETURNING id`
			slog.Debug("Сформированный SQL-запрос", "query", query, "args", []interface{}{groupName})

			res := r.db.QueryRow(context.Background(), query, groupName)

			var id int
			err = res.Scan(&id)
			if err != nil {
				slog.Error("Ошибка при создании группы", "error", err)
				return group, err
			}
			group.ID = &id
			group.Name = &groupName
			return group, nil
		}
		slog.Error("Ошибка при выборе группы", "error", err)
		return model.Group{}, err
	}

	slog.Info("Группа успешно найдена", "group", group)
	return group, nil
}
