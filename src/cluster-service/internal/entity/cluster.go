package entity

type Centroid struct {
	ID    string    `db:"uuid"`
	Value []float64 `db:"value"`
}

type Cluster struct {
	ID       string   `db:"uuid"`
	Index    int      `db:"index"`
	Centroid Centroid `db:"centroid_id"`
}
