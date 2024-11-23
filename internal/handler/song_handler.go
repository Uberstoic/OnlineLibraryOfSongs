package handler

import (
	"github.com/effectivemobile/music-library/internal/dto"
	"github.com/effectivemobile/music-library/internal/model"
	"github.com/effectivemobile/music-library/internal/repository"
	"github.com/effectivemobile/music-library/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SongHandler struct {
	repo     *repository.SongRepository
	musicAPI *service.MusicAPIService
}

func NewSongHandler(repo *repository.SongRepository, musicAPI *service.MusicAPIService) *SongHandler {
	return &SongHandler{
		repo:     repo,
		musicAPI: musicAPI,
	}
}

// Create godoc
// @Summary Create a new song
// @Description Create a new song with external API enrichment
// @Tags songs
// @Accept json
// @Produce json
// @Param song body dto.SongRequest true "Song information"
// @Success 201 {object} dto.SongResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /songs [post]
func (h *SongHandler) Create(c *gin.Context) {
	var req dto.SongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Getting song info for %s - %s", req.Group, req.Song)
	songDetail, err := h.musicAPI.GetSongInfo(req.Group, req.Song)
	if err != nil {
		log.Printf("Error getting song info: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to get song details"})
		return
	}

	releaseDate, err := time.Parse("02.01.2006", songDetail.ReleaseDate)
	if err != nil {
		log.Printf("Error parsing release date: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Invalid release date format"})
		return
	}

	song := &model.Song{
		Group:       req.Group,
		Song:        req.Song,
		ReleaseDate: releaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	if err := h.repo.Create(song); err != nil {
		log.Printf("Error creating song: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create song"})
		return
	}

	c.JSON(http.StatusCreated, dto.SongResponse{
		ID:          song.ID,
		Group:       song.Group,
		Song:        song.Song,
		ReleaseDate: song.ReleaseDate,
		Text:        song.Text,
		Link:        song.Link,
	})
}

// List godoc
// @Summary Get songs with pagination and filtering
// @Description Get a list of songs with pagination and filtering options
// @Tags songs
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Param group query string false "Filter by group name"
// @Param song query string false "Filter by song name"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /songs [get]
func (h *SongHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	
	filters := make(map[string]interface{})
	if group := c.Query("group"); group != "" {
		filters["group"] = group
	}
	if song := c.Query("song"); song != "" {
		filters["song"] = song
	}

	songs, total, err := h.repo.List(page, size, filters)
	if err != nil {
		log.Printf("Error listing songs: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to list songs"})
		return
	}

	var response []dto.SongResponse
	for _, song := range songs {
		response = append(response, dto.SongResponse{
			ID:          song.ID,
			Group:       song.Group,
			Song:        song.Song,
			ReleaseDate: song.ReleaseDate,
			Text:        song.Text,
			Link:        song.Link,
		})
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Total: total,
		Page:  page,
		Size:  size,
		Data:  response,
	})
}

// GetLyrics godoc
// @Summary Get song lyrics with pagination by verses
// @Description Get song lyrics divided into verses with pagination
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Verses per page" default(4)
// @Success 200 {object} dto.LyricsResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /songs/{id}/lyrics [get]
func (h *SongHandler) GetLyrics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid ID format"})
		return
	}

	song, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Song not found"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "4"))

	// Split text into verses (separated by double newline)
	verses := strings.Split(song.Text, "\n\n")
	total := int64(len(verses))

	// Calculate pagination
	start := (page - 1) * size
	end := start + size
	if start >= len(verses) {
		c.JSON(http.StatusOK, dto.LyricsResponse{
			Total:  total,
			Page:   page,
			Size:   size,
			Verses: []string{},
		})
		return
	}
	if end > len(verses) {
		end = len(verses)
	}

	c.JSON(http.StatusOK, dto.LyricsResponse{
		Total:  total,
		Page:   page,
		Size:   size,
		Verses: verses[start:end],
	})
}

// Delete godoc
// @Summary Delete a song
// @Description Delete a song by ID
// @Tags songs
// @Param id path int true "Song ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /songs/{id} [delete]
func (h *SongHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid ID format"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		log.Printf("Error deleting song: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to delete song"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Update godoc
// @Summary Update a song
// @Description Update a song's information
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body dto.SongRequest true "Updated song information"
// @Success 200 {object} dto.SongResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /songs/{id} [put]
func (h *SongHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid ID format"})
		return
	}

	song, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Song not found"})
		return
	}

	var req dto.SongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	songDetail, err := h.musicAPI.GetSongInfo(req.Group, req.Song)
	if err != nil {
		log.Printf("Error getting song info: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to get song details"})
		return
	}

	releaseDate, err := time.Parse("02.01.2006", songDetail.ReleaseDate)
	if err != nil {
		log.Printf("Error parsing release date: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Invalid release date format"})
		return
	}

	song.Group = req.Group
	song.Song = req.Song
	song.ReleaseDate = releaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	if err := h.repo.Update(song); err != nil {
		log.Printf("Error updating song: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, dto.SongResponse{
		ID:          song.ID,
		Group:       song.Group,
		Song:        song.Song,
		ReleaseDate: song.ReleaseDate,
		Text:        song.Text,
		Link:        song.Link,
	})
}
