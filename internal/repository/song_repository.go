package repository

import (
	"github.com/effectivemobile/music-library/internal/model"
	"gorm.io/gorm"
)

type SongRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) *SongRepository {
	return &SongRepository{
		db: db,
	}
}

func (r *SongRepository) Create(song *model.Song) error {
	return r.db.Create(song).Error
}

func (r *SongRepository) GetByID(id uint) (*model.Song, error) {
	var song model.Song
	err := r.db.First(&song, id).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (r *SongRepository) Update(song *model.Song) error {
	return r.db.Save(song).Error
}

func (r *SongRepository) Delete(id uint) error {
	return r.db.Delete(&model.Song{}, id).Error
}

func (r *SongRepository) List(page, size int, filters map[string]interface{}) ([]model.Song, int64, error) {
	var songs []model.Song
	var total int64

	query := r.db.Model(&model.Song{})

	// Apply filters
	for field, value := range filters {
		query = query.Where(field+" ILIKE ?", "%"+value.(string)+"%")
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.Offset((page - 1) * size).Limit(size).Find(&songs).Error
	if err != nil {
		return nil, 0, err
	}

	return songs, total, nil
}
