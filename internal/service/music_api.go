package service

import (
	"encoding/json"
	"fmt"
	"github.com/effectivemobile/music-library/internal/dto"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type MusicAPIService struct {
	baseURL string
}

func NewMusicAPIService(baseURL string) *MusicAPIService {
	return &MusicAPIService{
		baseURL: baseURL,
	}
}

func (s *MusicAPIService) GetSongInfo(group, song string) (*dto.SongDetail, error) {
	// URL encode parameters
	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)
	
	url := fmt.Sprintf("%s/info?%s", s.baseURL, params.Encode())
	log.Printf("Making request to: %s", url)
	
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("API returned non-200 status code: %d, body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	log.Printf("Response body: %s", string(body))

	var songDetail dto.SongDetail
	if err := json.Unmarshal(body, &songDetail); err != nil {
		log.Printf("Failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &songDetail, nil
}
