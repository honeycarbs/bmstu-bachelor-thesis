package postgres

import (
	"github.com/jmoiron/sqlx"
	"ml/internal/entity"
	"ml/pkg/psqlcli"
)

type FramePostgres struct {
	db *sqlx.DB
}

func NewFramePostgres(cli *psqlcli.Client) *FramePostgres {
	return &FramePostgres{db: cli.DB}
}

func (f *FramePostgres) GetBySample(sampleHash string) ([]entity.Frame, error) {
	var frames []entity.Frame

	query := `SELECT f.uuid, f.index as frame_index, c.index as cluster_index FROM
            	cluster c INNER JOIN frame f ON f.cluster_uuid = c.uuid 
                WHERE sample_uuid = $1`
	err := f.db.Select(&frames, query, sampleHash)
	if err != nil {
		panic(err)
	}

	return frames, nil
}
