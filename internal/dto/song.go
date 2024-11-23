package dto

import "time"

type SongRequest struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}

type SongResponse struct {
	ID          uint      `json:"id"`
	Group       string    `json:"group"`
	Song        string    `json:"song"`
	ReleaseDate time.Time `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type PaginatedResponse struct {
	Total int64         `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
	Data  interface{}  `json:"data"`
}

type LyricsResponse struct {
	Total  int64    `json:"total"`
	Page   int      `json:"page"`
	Size   int      `json:"size"`
	Verses []string `json:"verses"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
