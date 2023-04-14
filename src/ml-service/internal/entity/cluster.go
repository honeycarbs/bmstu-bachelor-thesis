package entity

type Centroid struct {
	ID    string    `db:"id"`
	Value []float64 `db:"value"`
}

type Cluster struct {
	ID       string   `db:"id"`
	Centroid Centroid `db:"centroid_id"`
}
