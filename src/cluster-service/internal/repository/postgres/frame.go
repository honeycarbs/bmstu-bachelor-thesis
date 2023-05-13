package postgres

import (
	"cluster/internal/entity"
	"cluster/pkg/psqlcli"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
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

func (f *FramePostgres) AssignCluster(clusterID, frameID string) error {
	query := "UPDATE frame SET cluster_uuid = $1 WHERE uuid = $2"
	_, err := f.db.Exec(query, clusterID, frameID)

	return err
}

func (f *FramePostgres) GetAll() ([]entity.Frame, error) {
	var frames []entity.Frame
	query := `SELECT f.uuid, f.sample_uuid, f.index, array_agg(m.value ORDER BY m.index) AS mfccs
			FROM frame f INNER JOIN mfcc m ON f.uuid = m.frame_uuid
			GROUP BY f.uuid, f.sample_uuid, f.index, f.cluster_uuid`
	rows, err := f.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			uuid       string
			sampleUUID string
			index      int
			mfccsStr   string
		)

		err = rows.Scan(&uuid, &sampleUUID, &index, &mfccsStr)

		if err != nil {
			return nil, err
		}
		mfccs, err := parseMFCC(mfccsStr)
		if err != nil {
			return nil, err
		}

		frames = append(frames, entity.Frame{
			ID:         uuid,
			SampleUUID: sampleUUID,
			Index:      index,
			MFCCs:      mfccs,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return frames, nil
}

func parseMFCC(mfccStr string) ([]float64, error) {
	mfccStr = strings.Trim(mfccStr, "{}")    // Remove curly braces
	strValues := strings.Split(mfccStr, ",") // Split by comma
	floatValues := make([]float64, 0, len(strValues))

	for _, strValue := range strValues {
		floatValue, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error while parsing float value %v: %v", strValue, err))
		}
		floatValues = append(floatValues, floatValue)
	}
	return floatValues, nil
}
