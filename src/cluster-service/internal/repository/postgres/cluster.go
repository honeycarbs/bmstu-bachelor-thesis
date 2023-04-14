package postgres

import (
	"cluster/internal/entity"
	"cluster/pkg/psqlcli"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ClusterPostgres struct {
	db *sqlx.DB
}

func NewClusterPostgres(cli *psqlcli.Client) *ClusterPostgres {
	return &ClusterPostgres{db: cli.DB}
}

func (c *ClusterPostgres) Create(cluster entity.Cluster) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	err = createCentroid(tx, cluster.Centroid)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = createCentroidCoords(tx, cluster.Centroid)
	if err != nil {
		tx.Rollback()
		return err
	}

	query := "INSERT INTO cluster(uuid, index, centroid_uuid) VALUES ($1, $2, $3)"
	_, err = tx.Exec(query, cluster.ID, cluster.Index, cluster.Centroid.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func createCentroid(tx *sql.Tx, centroid entity.Centroid) error {
	query := "INSERT INTO centroid(uuid) VALUES ($1)"
	_, err := tx.Exec(query, centroid.ID)

	return err
}

func createCentroidCoords(tx *sql.Tx, centroid entity.Centroid) error {
	for i, coord := range centroid.Value {
		query := "INSERT INTO centroid_coords (uuid, centroid_uuid, index, value) VALUES ($1, $2, $3, $4)"
		_, err := tx.Exec(query, uuid.New().String(), centroid.ID, i+1, coord)
		if err != nil {
			return err
		}
	}
	return nil
}
