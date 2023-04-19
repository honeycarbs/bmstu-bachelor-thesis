package postgres

import (
	"cluster/internal/entity"
	"cluster/pkg/psqlcli"
	"github.com/jmoiron/sqlx"
)

type SamplePostgres struct {
	db *sqlx.DB
}

func NewSamplePostgres(cli *psqlcli.Client) *SamplePostgres {
	return &SamplePostgres{db: cli.DB}
}

func (s *SamplePostgres) Get() ([]entity.Sample, error) {
	var samples []entity.Sample

	err := s.db.Select(&samples,
		`SELECT uuid, audio_path, emotion FROM sample`)
	if err != nil {
		return nil, err
	}

	return samples, nil
}
