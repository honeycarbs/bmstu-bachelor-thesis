package postgres

import (
	"github.com/jmoiron/sqlx"
	"ml/internal/entity"
	"ml/pkg/psqlcli"
)

type SamplePostgres struct {
	db *sqlx.DB
}

func NewSamplePostgres(cli *psqlcli.Client) *SamplePostgres {
	return &SamplePostgres{db: cli.DB}
}

func (s *SamplePostgres) GetByLabel(label entity.Label) ([]entity.Sample, error) {
	var samples []entity.Sample

	query := "SELECT uuid, audio_path, emotion FROM sample WHERE emotion = $1"

	err := s.db.Select(&samples,
		query, label)
	if err != nil {
		return nil, err
	}

	return samples, nil
}
