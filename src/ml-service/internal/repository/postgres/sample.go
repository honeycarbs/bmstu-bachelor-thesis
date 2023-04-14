package postgres

import (
	"github.com/jmoiron/sqlx"
	"ml/internal/entity"
)

type SamplePostgres struct {
	db *sqlx.DB
}

func (s *SamplePostgres) getByLabel(label string) ([]entity.Sample, error) {
	var samples []entity.Sample

	err := s.db.Select(&samples,
		`SELECT hash, audio_path, emotion FROM sample`)
	if err != nil {
		return nil, err
	}

	return samples, nil
}
