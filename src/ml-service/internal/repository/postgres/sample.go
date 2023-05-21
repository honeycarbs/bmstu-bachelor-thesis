package postgres

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"ml/internal/entity"
	"ml/pkg/logging"
	"ml/pkg/psqlcli"
)

type SamplePostgres struct {
	logger logging.Logger
	db     *sqlx.DB
}

func NewSamplePostgres(cli *psqlcli.Client, logger logging.Logger) *SamplePostgres {
	return &SamplePostgres{logger: logger, db: cli.DB}
}

func (s *SamplePostgres) GetByLabelTrain(label entity.Label) ([]entity.Sample, error) {
	var samples []entity.Sample

	query := "SELECT uuid, audio_path, emotion FROM sample WHERE emotion = $1 AND batch = 'train'"

	err := s.db.Select(&samples,
		query, label)
	if err != nil {
		return nil, err
	}

	return samples, nil
}

func (s *SamplePostgres) GetByLabelTest(label entity.Label) ([]entity.Sample, error) {
	var samples []entity.Sample

	query := "SELECT uuid, audio_path, emotion FROM sample WHERE emotion = $1 AND batch = 'test'"

	err := s.db.Select(&samples,
		query, label)
	if err != nil {
		return nil, err
	}

	return samples, nil
}

func (s *SamplePostgres) GetByPath(path string) (entity.Sample, error) {
	var sample entity.Sample

	query := "SELECT uuid, audio_path FROM sample WHERE audio_path = $1"
	err := s.db.Get(&sample, query, path)
	if err != nil {
		if err == sql.ErrNoRows {
			return sample, fmt.Errorf("sample not found for audio path: %s", path)
		}
		s.logger.Info(err)
		return sample, err
	}

	return sample, nil
}
