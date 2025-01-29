package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Xapsiel/EffectiveMobile/internal/model"
	"github.com/gin-gonic/gin"
)

type Song struct {
	SongName string `json:"song"`
	Group    string `json:"group"`
}

//	@Summary		Получение списка песен
//
// @Tags			songs
// @Description	Получение списка песен из базы данных с фильтрацией по параметрам
// @Accept			json
// @Produce		json
// @Param			song	query		string	false	"Название песни"	default(Supermassive Black Hole)
// @Param			group	query		string	false	"Группа"			default(Muse)
// @Param			link	query		string	false	"Ссылка на клип"			default(https://www.youtube.com/watch?v=Xsp3_a-PMTw)
// @Param			text	query		string	false	"Текст песни"			default(Ooh baby, don't you know I suffer?\nOoh baby, can my soul alight)
// @Param			date		query		string	false	"Дата публикации" example(19.07.2006)
// @Param			page	query		int		false	"Номер страницы"		default(1)
// @Param			limit	query		int		false	"Количество на странице"		default(10)
// @Success		200		{object}	[]model.Song
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Router			/info [get]
func (h *Handler) GetSongs(c *gin.Context) {
	slog.Info("Начало обработки запроса GetSongs")

	group := c.DefaultQuery("group", "")
	group_id, err := strconv.Atoi(c.DefaultQuery("group_id", "-1"))
	if err != nil {
		slog.Error("Ошибка при парсинге group_id", "error", err)
		newErrorResponce(c, http.StatusBadRequest, fmt.Sprintf("Invalid group id %v", group))
		return
	}

	song := c.DefaultQuery("song", "")
	id, err := strconv.Atoi(c.DefaultQuery("id", "-1"))
	if err != nil {
		slog.Error("Ошибка при парсинге id", "error", err)
		newErrorResponce(c, http.StatusBadRequest, fmt.Sprintf("Invalid id %v", id))
		return
	}

	date := c.DefaultQuery("date", "0001-01-01")
	link := c.DefaultQuery("link", "")
	text := c.DefaultQuery("text", "")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		slog.Error("Ошибка при парсинге page", "error", err)
		newErrorResponce(c, http.StatusBadRequest, fmt.Sprintf("Invalid page %v", page))
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		slog.Error("Ошибка при парсинге limit", "error", err)
		newErrorResponce(c, http.StatusBadRequest, fmt.Sprintf("Invalid limit %v", limit))
		return
	}

	filter := model.Song{
		ID:          &id,
		SongName:    &song,
		GroupId:     &group_id,
		Group:       &group,
		ReleaseDate: &date,
		Link:        &link,
		Text:        &text,
	}

	slog.Debug("Параметры фильтра", "filter", filter)

	res, err := h.service.GetSongs(filter, page, limit)
	if err != nil {
		slog.Error("Ошибка при получении песен", "error", err)
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(res) == 0 {
		slog.Info("Список песен пуст")
		res = make([]model.Song, 0)
		c.AbortWithStatusJSON(200, res)
		return
	}

	slog.Info("Успешно получен список песен", "количество песен", len(res))
	c.AbortWithStatusJSON(200, res)
}

// @Summary		Добавление новой песни
// @Description	Добавление новой песни в базу данных (Обязательные параметры - song, group)
// @Tags			songs
// @Accept			json
// @Produce		json
// @Param song body Song true "Данные песни" default({ "group": "Muse", "song": "Supermassive Black Hole" })
// @Success		200		{object}	resultResponse
// @Failure		400		{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Router			/songs [post]
func (h *Handler) AddSong(c *gin.Context) {
	slog.Info("Начало обработки запроса AddSong")

	var song Song
	if err := c.ShouldBindJSON(&song); err != nil {
		slog.Error("Ошибка при парсинге JSON", "error", err)
		newErrorResponce(c, http.StatusBadRequest, fmt.Sprintf("Parameter error: %v", err))
		return
	}

	slog.Debug("Данные песни", "song", song)

	id, err := h.service.Add(song.SongName, song.Group)
	if err != nil {
		slog.Error("Ошибка при добавлении песни", "error", err)
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("Песня успешно добавлена", "id", id)
	c.AbortWithStatusJSON(200, resultResponse{
		Status: "success",
		Id:     id,
		Text:   "Песня добавлена",
	})
}

// @Summary Получение текста куплета песни
// @Description Получение текста конкретного куплета
// @Tags songs
// @Accept json
// @Produce json
// @Param song query string false "Название песни" default(Supermassive Black Hole)
// @Param group query string false "Группа" default(Muse)
// @Param verse query int false "Номер куплета" default(1)
// @Success 200 {object} resultResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /info/verse [get]
func (h *Handler) GetSongVerse(c *gin.Context) {
	slog.Info("Начало обработки запроса GetSongVerse")

	group := c.DefaultQuery("group", "")
	song_name := c.DefaultQuery("song", "")
	verseNumber, err := strconv.Atoi(c.DefaultQuery("verse", "1"))
	if err != nil {
		slog.Warn("Ошибка при парсинге verse, установлено значение по умолчанию", "verseNumber", 1)
		verseNumber = 1
	}

	slog.Debug("Параметры запроса", "group", group, "song_name", song_name, "verseNumber", verseNumber)

	song := model.Song{SongName: &song_name, Group: &group}
	verse, id, err := h.service.GetSongVerse(song, verseNumber)
	if err != nil {
		slog.Error("Ошибка при получении куплета", "error", err)
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("Куплет успешно получен", "id", id)
	c.AbortWithStatusJSON(200, resultResponse{
		Status: "success",
		Id:     id,
		Text:   verse,
	})
}

// @Summary Удаление песни
// @Description Удаление песни по ID, названию или группе
// @Tags songs
// @Accept json
// @Produce json
// @Param song body model.Song true "Данные песни для удаления"
// @Success 200 {object} resultResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /songs [delete]
func (h *Handler) DeleteSong(c *gin.Context) {
	slog.Info("Начало обработки запроса DeleteSong")

	var song model.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		slog.Error("Ошибка при парсинге JSON", "error", err)
		newErrorResponce(c, http.StatusUnprocessableEntity, "Ошибка парсинга структуры")
		return
	}

	slog.Debug("Данные песни для удаления", "song", song)

	ok, err := h.service.DeleteSong(song)
	if err != nil || !ok {
		slog.Error("Ошибка при удалении песни", "error", err)
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("Песня успешно удалена", "song_name", *song.SongName)
	c.AbortWithStatusJSON(200, resultResponse{
		Status: "success",
		Text:   "Удаление прошло успешно",
	})
}

// @Summary Обновление информации о песне
// @Description Обновление данных о песне
// @Tags songs
// @Accept json
// @Produce json
// @Param			song	query		string	false	"Название песни"			default(Supermassive Black Hole)
// @Param			group	query		string	false	"Группа"			default(Muse)
// @Param song body model.Song true "Данные песни для обновления"
// @Success 200 {object} resultResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /songs [put]
func (h *Handler) UpdateSong(c *gin.Context) {
	slog.Info("Начало обработки запроса UpdateSong")

	song_name := c.DefaultQuery("song", "")
	group_name := c.DefaultQuery("group", "")

	slog.Debug("Параметры запроса", "song_name", song_name, "group_name", group_name)

	var song model.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		slog.Error("Ошибка при парсинге JSON", "error", err)
		newErrorResponce(c, http.StatusUnprocessableEntity, "Ошибка парсинга структуры")
		return
	}

	slog.Debug("Данные песни для обновления", "song", song)

	ok, song, err := h.service.UpdateSong(song_name, group_name, song)
	if err != nil || !ok {
		slog.Error("Ошибка при обновлении песни", "error", err)
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("Информация о песне успешно обновлена", "song_name", *song.SongName, "id", *song.ID)
	c.AbortWithStatusJSON(200, resultResponse{
		Status: "success",
		Text:   "Обновление прошло успешно",
	})
}
