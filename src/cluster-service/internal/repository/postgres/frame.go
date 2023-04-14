package postgres

import (
	"cluster/internal/entity"
	"cluster/pkg/psqlcli"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"sort"
)

type coeff struct {
	Index int
	Value float64
}

type FramePostgres struct {
	db *sqlx.DB
}

func NewFramePostgres(cli *psqlcli.Client) *FramePostgres {
	return &FramePostgres{db: cli.DB}
}

func (f *FramePostgres) GetOne(sampleHash string, sampleNum int) (entity.Frame, error) {
	fm := entity.Frame{
		SampleHash: sampleHash,
	}

	err := f.db.Get(&fm,
		`SELECT uuid, index FROM frame WHERE sample_hash = $1 AND index = $2`, sampleHash, sampleNum)
	if err != nil {
		panic(err)
	}

	fm.MFCCs, err = f.getMFCCs(fm.ID)
	if err != nil {
		return entity.Frame{}, err
	}

	return fm, nil
}

func (f *FramePostgres) getMFCCs(frameId string) ([]float64, error) {
	var coefficients []coeff
	err := f.db.Select(&coefficients,
		`SELECT index, value FROM mfcc WHERE frame_uuid = $1`, frameId)
	if err != nil {
		panic(err)
	}

	sort.Slice(coefficients, func(i, j int) bool {
		return coefficients[i].Index < coefficients[j].Index
	})

	values := make([]float64, len(coefficients))
	for i, c := range coefficients {
		values[i] = c.Value
	}

	return values, err
}

func (f *FramePostgres) AssignCluster(clusterID, frameID string) error {
	query := "UPDATE frame SET cluster_uuid = $1 WHERE uuid = $2"
	_, err := f.db.Exec(query, clusterID, frameID)

	return err
}

func (f *FramePostgres) CountPerSample(sampleHash string) (int, error) {
	query := "select count(uuid) from frame where sample_hash =  $1"

	var count int
	if err := f.db.Get(&count, query, sampleHash); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("sample does not exist")
		} else {
			return 0, err
		}
	}

	return count, nil
}
